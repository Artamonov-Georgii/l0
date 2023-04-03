package main

import (
	"log"
	"time"

	"github.com/nats-io/stan.go"
)

func main() {
	clusterID := "test-cluster"
	clientID := "test-publisher"
	subject := "my-subject"

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		msg := []byte("Hello NATS Streaming!")
		if err := sc.Publish(subject, msg); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
	}

	sc.Close()
}
