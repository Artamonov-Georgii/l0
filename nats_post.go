package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/nats-io/nats.go"
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

	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatalf("Error connecting to NATS server: %v", err)
	}
	defer nc.Close()

	// Connect to NATS Streaming
	sc, err := stan.Connect("test-cluster", "test-client", stan.NatsConn(nc))
	if err != nil {
		log.Fatal(err, "STAN")
		return
	}
	defer sc.Close()

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

	for i := 1; i <= 100; i++ {
		order.OrderUID = fmt.Sprintf("%s-%d", order.OrderUID, i)
		msg, err := json.Marshal(order)

		if err != nil {
			log.Fatal(err)
		}

		sc.Publish("orders", msg)
	}

	println("\n \n the publishing is complete")

}
