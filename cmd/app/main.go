package main

import (
	"log"
	"time"

	"github.com/g0dm0d/wbtest/internal/config"
	"github.com/g0dm0d/wbtest/internal/server"
	"github.com/g0dm0d/wbtest/internal/service"
	"github.com/g0dm0d/wbtest/internal/store/postgres"
	"github.com/g0dm0d/wbtest/pkg/cache"
	"github.com/nats-io/stan.go"
)

func main() {
	config, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	// Init db conn
	db, err := postgres.New(config.Postgres.DSN)
	if err != nil {
		log.Fatal(err)
	}

	// Init stores
	orderStore := postgres.NewOrderStore(db)

	// Init request cache store
	cacheMap := cache.NewCacheMap()

	// Run simple gc every 10 minutes. So that the cache does not contain too much unnecessary data
	go func() {
		for {
			cacheMap.GCCollector()
			time.Sleep(10 * time.Minute)
		}
	}()

	// Init nats conn
	sc, err := stan.Connect(config.Nats.ClusterID, config.Nats.ClientID, stan.NatsURL(config.Nats.URL))
	if err != nil {
		log.Fatal(err)
	}

	// Init services
	services := service.New(orderStore, cacheMap)

	// Init server
	server := server.NewServer(&server.Config{
		Addr:     config.App.Addr,
		Port:     config.App.Port,
		Service:  services,
		StanConn: sc,
	})

	server.SetupRouter()

	log.Print("Server is up and running.")

	err = server.RunServer()
	if err != nil {
		log.Panic(err)
	}
}
