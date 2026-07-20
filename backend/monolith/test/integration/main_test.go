package integration

import (
	"context"
	"log"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"monolith/test/testutils"
)

var (
	adminConn       *pgxpool.Pool
	adminConnString string
)

const templateDBName = ""

func TestMain(m *testing.M) {
	ctx := context.Background()
	var err error
	var cleanup func()
	adminConn, adminConnString, cleanup, err = testutils.CreateAdminDBReusable(ctx)
	if err != nil {
		log.Fatalf("Failed to create admin DB: %e", err)
	}
	defer cleanup()

	err = testutils.CreateTemplateDB(ctx, adminConn, adminConnString, templateDBName)
	if err != nil {
		log.Fatalf("Failed to create template DB: %v", err)
	}

	m.Run()
}
