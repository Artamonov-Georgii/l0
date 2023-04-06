package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"encoding/base64"

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

func randomJSON() []byte {
	data := make(map[string]interface{})
	data["name"] = generateRandomString(8)
	data["age"] = rand.Intn(100)
	data["email"] = fmt.Sprintf("%s@example.com", generateRandomString(8))
	data["is_active"] = rand.Float32() < 0.5
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return jsonData
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

func main() {

	clusterID := "test-cluster"
	clientID := "test-publisher"
	subject := "my-subject"

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Fatal(err)
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

	for i := 1; i <= 10; i++ {

		if i%2 == 0 {
			msg := randomJSON()
			if err := sc.Publish(subject, msg); err != nil {
				log.Fatal(err)
				
			}
			continue
		}

		order.OrderUID = generateRandomString(10)
		msg, err := json.Marshal(order)

		if err != nil {
			log.Fatal(err)
		}

		if err := sc.Publish(subject, msg); err != nil {
            log.Fatal(err)
        }
	}

	println("\n \n the publishing is complete")

}
