package main

import (
	"database/sql"
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

func main() {

	// Устанавливаем соединение с Постгре
	db, err_db := sql.Open("postgres", "postgres://postgres:postgrespw@localhost:32771/postgres?sslmode=disable")

	if err_db != nil {
		log.Fatal("Not possible to connect to the database")
	}

	fmt.Println("DB connection established")

	server.CacheMsg, _ = pgsql.GetAllOrders(db)

	fmt.Println("Cache operations are done")

	defer db.Close()

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

	// Теперь сабскрайбим на тему в натсе

	sub, err_sub := sc.Subscribe(subject, func(msg *stan.Msg) {
		log.Printf("Received message on subject")
		if order, error_mess := nats_checker.IsValidOrderMsg(msg); error_mess == nil {
			server.CacheMsg = append(server.CacheMsg, order)
			err_inst := pgsql.InsertOrder(db, order) 
			fmt.Println(server.CacheMsg)
			if err_inst != nil {
                fmt.Println(err_inst)
			}	
			fmt.Printf("Receive successfully")
		} else {
			fmt.Println("message is not correct")
		}})

	if err_sub != nil {
		log.Fatal(err_sub)
	}

    defer sub.Unsubscribe()
		

	signalChan := make(chan os.Signal, 1)
    signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
    <-signalChan
    fmt.Println("Received signal, shutting down...")
    os.Exit(0)

}
