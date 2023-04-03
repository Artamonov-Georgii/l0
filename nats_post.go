package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"os"
    "os/signal"
    "syscall"

	"github.com/nats-io/stan.go"
)

type Order struct {
	OrderUID  string   `json:"order_uid"`
	TrackNum  string   `json:"track_number"`
	Entry     string   `json:"entry"`
	Delivery  Delivery `json:"delivery"`
	Payment   Payment  `json:"payment"`
	Items     []Item   `json:"items"`
	Locale    string   `json:"locale"`
	Customer  string   `json:"customer_id"`
	Service   string   `json:"delivery_service"`
	ShardKey  string   `json:"shardkey"`
	SMID      int      `json:"sm_id"`
	CreatedAt string   `json:"date_created"`
	OOFShard  string   `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type Item struct {
	ChrtID     int    `json:"chrt_id"`
	TrackNum   string `json:"track_number"`
	Price      int    `json:"price"`
	RID        string `json:"rid"`
	Name       string `json:"name"`
	Sale       int    `json:"sale"`
	Size       string `json:"size"`
	TotalPrice int    `json:"total_price"`
	NmID       int    `json:"nm_id"`
	Brand      string `json:"brand"`
	Status     int    `json:"status"`
}

func main() {

	clusterID := "test-cluster"
	clientID := "test-publisher"
	subject := "my-subject"

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Fatal(err)
	}

	// Read JSON file into order struct
	order := &Order{}
	file, err := ioutil.ReadFile("model.json")

	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
		return
	}

	if err := json.Unmarshal(file, order); err != nil {
		log.Fatal(err)
		fmt.Println(err)
		return
	}

	fmt.Println(*order)

	for i := 1; i <= 10; i++ {
		order.OrderUID = fmt.Sprintf("%s-%d", order.OrderUID, i)
		msg, err := json.Marshal(order)

		if err != nil {
			log.Fatal(err)
		}

		if err := sc.Publish(subject, msg); err != nil {
            log.Fatal(err)
        }
	}

	println("\n \n the publishing is complete")

	signalChan := make(chan os.Signal, 1)
    signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
    <-signalChan
    fmt.Println("Received signal, shutting down...")
    os.Exit(0)
}
