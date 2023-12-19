package main

import (
	"log"
	"os"

	"github.com/nats-io/stan.go"
)

func main() {
	clusterID := "test-cluster"
	clientID := "client-test"

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	filePath := "model.json"

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	subject := "order.pipeline"
	err = sc.Publish(subject, fileContent)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Success send data")
}
