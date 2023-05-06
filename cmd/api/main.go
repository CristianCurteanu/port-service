package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/CristianCurteanu/koken-api/internal/domains/ports"
	"github.com/CristianCurteanu/koken-api/internal/infra/http"
	"github.com/CristianCurteanu/koken-api/internal/infra/storage/database"
	"github.com/CristianCurteanu/koken-api/internal/infra/storage/inmemory"
)

var (
	port        *int
	mongoDbUrl  *string
	mongoDbName *string
)

func init() {
	port = flag.Int("port", 8080, "Port on which server will listen for requests")
	mongoDbName = flag.String("mongo-db-name", "", "The database name for MongoDB storage")
	mongoDbUrl = flag.String("mongo-db-url", "", "The URL for MongoDB storage")
}

func main() {
	flag.Parse()

	app, err := http.BuildApp(*port,
		http.PortHandlers(createPortService()),
	)
	if err != nil {
		panic(err)
	}

	go func() {
		err = app.Run()
		if err != nil {
			log.Printf("[server error]: %s\n", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	defer close(shutdown)

	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	err = app.Close()
	if err != nil {
		log.Fatalf("[server-force-shutdown]:\n\tError %q", err.Error())
	}

	log.Println("[server-exit]: OK")
}

func createPortService() ports.PortService {
	if mongoDbUrl == nil || mongoDbName == nil {
		return nil
	}
	// TODO: use a DI container, like wire
	portsDbStorage, err := database.NewMongoDB(context.Background(), *mongoDbUrl, *mongoDbName, "ports")
	if err != nil {
		panic(err)
	}
	return ports.NewPortService(ports.NewPortRepositories(inmemory.NewInMemoryStorage(), portsDbStorage))
}
