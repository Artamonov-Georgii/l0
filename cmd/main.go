package main

import (
	"fmt"
	"log"
	"os"
    "os/signal"
    "syscall"

	"github.com/Artamonov-Georgii/l0/internal/nats"
	"github.com/Artamonov-Georgii/l0/internal/pgsql"
	"github.com/Artamonov-Georgii/l0/internal/server"
	
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
)

var err_cache error

func main() {

	// Устанавливаем соединение с Постгре

	pgsql.StartSql()
	defer pgsql.Db.Close()

	fmt.Println("DB connection established")
	
	// Make nonnil map
	
	server.CacheMsg, err_cache = pgsql.GetAllOrders(pgsql.Db)
	fmt.Println("Cache Len:", len(server.CacheMsg))

	if err_cache != nil {
        fmt.Println(err_cache)
    }

	fmt.Println("Cache operations are done")

	// Устанавливаем соединение с Натс-Стриминг

	clusterID := "test-cluster"
	clientID := "test-subscriber"
	subject := "my-subject"

	sc, err_stan := stan.Connect(clusterID, clientID, stan.NatsURL("nats://localhost:4222"))
	
	if err_stan != nil {
		log.Fatal(err_stan)
	}

	defer sc.Close()

	// Запускаем http

	go server.Run()

	// Теперь сабскрайбим на тему в натсе и вставляем функцию обработки сообщений

	sub, err_sub := sc.Subscribe(subject, nats_checker.MessageHandler)

	if err_sub != nil {
		log.Fatal(err_sub)
	}

    defer sub.Unsubscribe()
		

	signalChan := make(chan os.Signal, 1)
    signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
    <-signalChan
    fmt.Println("\n Received signal, shutting down...")
    os.Exit(0)

}
