package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/stan.go"
)

func main() {
	clusterID := "test-cluster"
	clientID := "test-subscriber"
	subject := "my-subject"

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Fatal(err)
	}

	defer sc.Close()

	log.Println("Connected to NATS Streaming server")

	sub, err := sc.Subscribe(subject, func(msg *stan.Msg) {
		log.Printf("Received message on subject %s: %s\n", msg.Subject, string(msg.Data))
	})
	if err != nil {
		log.Fatal(err)
	}

    defer sub.Unsubscribe()
    
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c

	log.Println("Received interrupt signal, shutting down...")
}
