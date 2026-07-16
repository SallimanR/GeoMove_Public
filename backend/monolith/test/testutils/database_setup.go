package testutils

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	containerName = "postgres"
	containerPort = "5432"

	defaultTemplateDBName = "template_db"
)

func CreateAdminDB(ctx context.Context) (*pgxpool.Pool, string, func(), error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, "", nil, fmt.Errorf("failed to get path of setup file")
	}
	dockerComposeFile := filepath.Join(filepath.Dir(filename), "database/docker-compose.test.yaml")

	dockerCompose, err := compose.NewDockerCompose(dockerComposeFile)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to setup postgres docker compose file: %w", err)
	}

	dockerCompose.WaitForService(
		containerName,
		wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForExec([]string{"pg_isready", "-U", "postgres"}).
				WithPollInterval(20*time.Millisecond).
				WithStartupTimeout(10*time.Second),
		),
	)
	err = dockerCompose.Up(ctx)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to up postgres docker compose: %w", err)
	}

	dockerContainer, err := dockerCompose.ServiceContainer(ctx, containerName)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to get postgres container: %w", err)
	}

	host, err := dockerContainer.Host(ctx)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to get postgres container host: %w", err)
	}
	port, err := dockerContainer.MappedPort(ctx, containerPort)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to get postgres container port: %w", err)
	}

	adminConnString := fmt.Sprintf("postgres://postgres:postgres@%s:%s/postgres?sslmode=disable", host, port.Port())

	dbPool, err := pgxpool.New(ctx, adminConnString)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to create DB pool: %w", err)
	}

	cleanup := func() {
		dbPool.Close()

		err := dockerCompose.Down(ctx)
		if err != nil {
			log.Printf("Faield to close dbPool: %s", err)
		}
	}
	return dbPool, adminConnString, cleanup, nil
}

// FIXME: fix reusable db is dropped by test containers "RYUK" conatiner (GC).
// error in tests: ["Failed to create DB from template: unexpected EOF"]
func CreateAdminDBReusable(ctx context.Context) (*pgxpool.Pool, string, func(), error) {
	req := testcontainers.ContainerRequest{
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "postgres",
		},
		Tmpfs: map[string]string{"/var/lib/postgresql/18/docker": "rw,noexec,nosuid,size=180m"},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForExec([]string{"psql", "-U", "postgres", "-c", "SELECT 1"}),
		),
		Name: containerName, // fixed name for reuse
	}

	// NOTE: Setup from Docker image
	req.Image = "test-postgres-postgis-h3:latest"

	// NOTE: Setup from Docker file
	//
	// _, filename, _, ok := runtime.Caller(0)
	// if !ok {
	// 	return nil, "", nil, fmt.Errorf("failed to get path of setup file")
	// }
	// dockerFilePath := filepath.Join(filepath.Dir(filename), "database/Dockerfile.test")
	// dockerFileDir := filepath.Dir(dockerFilePath)
	// dockerFileName := filepath.Base(dockerFilePath)
	// req.FromDockerfile = testcontainers.FromDockerfile{
	// 	Context:    dockerFileDir,
	// 	Dockerfile: dockerFileName,
	// }

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Reuse:            true,
	})
	if err != nil {
		if container != nil {
			logs, _ := container.Logs(ctx)
			logBytes, _ := io.ReadAll(logs)
			fmt.Printf("Container logs:\n%s\n", string(logBytes))
			_ = logs.Close()
		}
		return nil, "", nil, fmt.Errorf("failed to create generic container for docker compose: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to get postgres container host: %w", err)
	}
	port, err := container.MappedPort(ctx, containerPort)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to get postgres container port: %w", err)
	}

	adminConnString := fmt.Sprintf("postgres://postgres:postgres@%s:%s/postgres?sslmode=disable", host, port.Port())

	dbPool, err := pgxpool.New(ctx, adminConnString)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to create DB pool: %w", err)
	}

	// Cleanup only closes pool - container stays alive for other test packages
	cleanup := func() {
		dbPool.Close()
	}
	return dbPool, adminConnString, cleanup, nil
}

func CreateTemplateDB(ctx context.Context, adminConn *pgxpool.Pool, adminConnString, templateDBName string) error {
	// FIXME: delay is needed to enshure that adminDB is ready
	maxRetries := 50
	for i := 0; i < maxRetries; i++ {
		pingCtx, cancel := context.WithTimeout(ctx, time.Millisecond*50)
		defer cancel()
		err := adminConn.Ping(pingCtx)
		if err == nil {
			break
		}
		if i == maxRetries-1 {
			return fmt.Errorf("failed to ping DB after %d retries: %w", maxRetries, err)
		}
		// Exponential backoff: 100ms, 200ms, 400ms, 800ms, 1600ms
		// sleep := time.Duration(50*(1<<i)) * time.Millisecond
		//
		sleep := 50 * time.Millisecond
		time.Sleep(sleep)
	}

	if templateDBName == "" {
		templateDBName = defaultTemplateDBName
	}
	_, err := adminConn.Exec(ctx, fmt.Sprintf(`CREATE DATABASE %s`, templateDBName))
	// If database already exists (from a previous run), treat as success
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return nil
		}
		return fmt.Errorf("failed to create template for DB")
	}

	// Get path of this file. when calling in tests path would be absolute to this file, not relative to test file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to get path of setup file")
	}
	migrationDir := filepath.Join(filepath.Dir(filename), "../../migrations")

	templateConnUrl, err := convertConnStringToAnotherDB(adminConnString, templateDBName)
	dbMate := dbmate.New(templateConnUrl)
	dbMate.AutoDumpSchema = false
	dbMate.MigrationsDir = []string{migrationDir}
	err = dbMate.Migrate()
	if err != nil {
		return fmt.Errorf("failed to apply migrations to template DB: %w", err)
	}

	_, err = adminConn.Exec(ctx, fmt.Sprintf(`ALTER DATABASE %s IS_TEMPLATE true`, templateDBName))
	if err != nil {
		return fmt.Errorf("failed to mark database as template: %w", err)
	}

	return nil
}

func CreateTestDB(tb testing.TB, ctx context.Context, adminConn *pgxpool.Pool, adminConnString, templateDBName string) (*pgxpool.Pool, string, func()) {
	tb.Helper()

	testDBName := fmt.Sprintf("test_%d_%d", time.Now().UnixNano(), os.Getpid())

	if templateDBName == "" {
		templateDBName = defaultTemplateDBName
	}
	_, err := adminConn.Exec(ctx, fmt.Sprintf(`CREATE DATABASE %s TEMPLATE %s`, testDBName, templateDBName))
	if err != nil {
		tb.Fatalf("Failed to create DB from template: %s", err)
	}

	testConnUrl, err := convertConnStringToAnotherDB(adminConnString, testDBName)
	if err != nil {
		tb.Fatalf("%s", err)
	}

	testConnString := testConnUrl.String()
	testDB, err := pgxpool.New(ctx, testConnString)
	if err != nil {
		tb.Fatalf("Failed to connect to test DB: %s", err)
	}

	cleanup := func() {
		stat := testDB.Stat()
		tb.Logf("acquired connections in db pool: %d", stat.AcquiredConns())
		testDB.Close()

		cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := adminConn.Exec(cleanupCtx, `
			SELECT pg_terminate_backend(pid)
			FROM pg_stat_activity
			WHERE datname = $1 AND pid <> pg_backend_pid()
		`, testDBName)
		if err != nil {
			tb.Logf("cleanup: failed to terminate backends for %s: %v", testDBName, err)
		}
		_, err = adminConn.Exec(cleanupCtx, fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, testDBName))
		if err != nil {
			tb.Logf("cleanup: failed to drop database %s: %v", testDBName, err)
		}
	}

	return testDB, testDBName, cleanup
}

// Function replaces "postgres" which is admin db, with <testDBName> to connect to test DB
//
// Before: postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
// After: postgres://postgres:postgres@localhost:5432/dbName?sslmode=disable
func convertConnStringToAnotherDB(connString string, dbName string) (*url.URL, error) {
	u, err := url.Parse(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse admin string: %v", err)
	}
	u.Path = "/" + dbName
	return u, nil
}
