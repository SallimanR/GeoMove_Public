package e2e

import (
	"context"
	"fmt"
	"testing"
	"time"

	"monolith/cmd/server"
	driverCommands "monolith/internal/domains/driver/application/commands"
	"monolith/internal/domains/geolocation/infrastructure/db/postgres"
	"monolith/internal/domains/geolocation/infrastructure/db/sqlc"
	geoWS "monolith/internal/domains/geolocation/interfaces/websocket"
	"monolith/internal/websockethub"
	"monolith/pkg/logger"
	"monolith/test/testutils"
	testDomains "monolith/test/testutils/domains"
	"monolith/test/testutils/network"
)

func TestE2E_DriversRealtimeOnMaps(t *testing.T) {
	t.Parallel()

	// ------------- 1. Setup server ---------------
	logger := logger.NewDebugLogger()

	ctx := context.Background()
	db, dbName, dbCleaunp := testutils.CreateTestDB(t, ctx, adminConn, adminConnString, templateDBName)
	t.Logf("dbName: %s", dbName)

	publisherRole := "tow_driver"
	subscriberRole := "tow_subscriber"
	wsOptions := websockethub.WebsocketServerOptions{
		Roles: []string{publisherRole, subscriberRole},
	}
	wsServer := websockethub.NewWebsocketServer(wsOptions)
	wsPort := network.GetFreePort(t)

	srv, err := server.NewServer(
		server.WithWebsocketServer(wsServer),
		server.WithWebsocketAddr(fmt.Sprintf("127.0.0.1:%d", wsPort)),
	)
	if err != nil {
		t.Fatalf("Failed to setup server: %s", err)
	}
	wsAddrs := srv.GetWebSocketAddrs()

	geoRepo := postgres.NewGeolocationRepository(sqlc.New(db))
	// updateGPSRelatime := commands.NewUpdateGPSRealtimeHandler(geoRepo)
	// _, channelName, err := geoWS.NewGPSRealtimeChannel(wsServer, []string{publisherRole, subscriberRole}, updateGPSRelatime, logger)
	_, channelName, err := geoWS.NewGPSRealtimeChannel(wsServer, []string{publisherRole, subscriberRole}, geoRepo, logger)
	if err != nil {
		t.Fatalf("failed to register the channel %s", err)
	}

	err = srv.Start()
	if err != nil {
		t.Fatalf("Failed to start server: %s", err)
	}
	time.Sleep(100 * time.Millisecond)
	wsServer.Run()

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err := srv.Shutdown(ctx)
		if err != nil {
			t.Logf("failed to Shutdown server: %s", err)
		}
		dbCleaunp()
	})

	// ------------- 2. Seed DB data ---------------
	driverDomain := testDomains.NewDriverDomain(t, db)
	createDriverCmd := driverCommands.CreateDriverCommand{
		Latitude:  55.5,
		Longitude: 38.4,
	}
	driver, err := driverDomain.Commands.CreateDriver.Handle(ctx, createDriverCmd)
	if err != nil {
		t.Fatalf("failed to create driver: %s", err)
	}
	err = geoRepo.CreateDriverRealtime(ctx, uint32(driver.ID))
	if err != nil {
		t.Fatalf("failed to create driver in realtime table: %s", err)
	}
}
