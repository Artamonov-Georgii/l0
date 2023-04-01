package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Artamonov-Georgii/l0/internal/postgreSQL"
	"github.com/Artamonov-Georgii/l0/internal/server"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

var CacheMsg []server.Order

func main() {

	// Устанавливаем соединение с Постгре
	db, err := sql.Open("postgres", "host=localhost:32769 user=postgres password=postgrespw dbname=postgres")

	if err != nil {
		log.Fatal("Not possible to connect to the database")
	}

	CacheMsg := postgreSQL.GetAllOrders(conn * sql.DB)

	if err != nil {
		log.Fatal("Not possible to restore cache")
	}

	defer db.Close()

	// Устанавливаем соединение с Натс-Стриминг
	nc, err := nats.Connect("nats://localhost:4222")

	if err != nil {
		log.Fatal("Not possible to connect to the nats-streaming server")
	}
	defer nc.Close()

	// Запускаем http

	server.Run()

	// Теперь сабскрайбим на тему в натсе

	sub, err_subs := nc.Subscribe("orders", func(msg *nats.Msg) {
		fmt.Printf("Получено сообщение: %s\n", string(msg.Data))
		order, err := nats_checker.IsValidOrderMsg(msg)

		if err == nil {
			CacheMsg = append(CacheMsg, order)
			pgsql.InsertOrderMsg(*db, order)
		} else {
			fmt.Println("Полученное сообщение не удалось сохранить")
		}
	})

	if err_subs != nil {
		log.Fatal("Not possible to subscribe to the topic")
	}

}
