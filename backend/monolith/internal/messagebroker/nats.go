package messagebroker

import (
	"log"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

func ConnectToCluster() (*nats.Conn, error) {
	opts := []nats.Option{
		nats.Name("websocket-hub"),
		nats.Timeout(10 * time.Second),
		nats.PingInterval(30 * time.Second),
		nats.MaxPingsOutstanding(5),
		nats.ReconnectWait(2 * time.Second),
		nats.MaxReconnects(-1),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Printf("Nats disconnected %v", err)
		}),
	}

	servers := []string{
		"nats://localhost:4022",
	}

	return nats.Connect(strings.Join(servers, ","), opts...)
}
