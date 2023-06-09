package server

import (
	"encoding/json"
	"fmt"
	"os/signal"
	"syscall"
	"os"
	"net/http"
)

type Order struct {
	OrderUID  string   `json:"order_uid"`
	TrackNum  string   `json:"track_number"`
	Entry     string   `json:"entry"`
	IntSig    string   `json:"internal_signature"`
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

var CacheMsg = make(map[string]Order)

func getOrder(w http.ResponseWriter, r *http.Request) {
	orderUID := r.URL.Query().Get("order_uid")

	o, ok := CacheMsg[orderUID]
	if !ok {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	jsonBytes, err := json.MarshalIndent(o, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, string(jsonBytes))
}

func Run() {
    http.HandleFunc("/order_uid", getOrder)

    fmt.Println("Starting server on :8080")

    server := &http.Server{Addr: ":8080"}

    go func() {
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            panic(fmt.Sprintf("Error starting server: %v", err))
        }
    }()

    // Wait for a signal to shutdown the server
    sigint := make(chan os.Signal, 1)
    signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
    <-sigint
    fmt.Println("Shutting down server...")
}


