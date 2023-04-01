package nats_checker

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
    "github.com/Artamonov-Georgii/l0/server" 
)

func IsValidOrderMsg(msg *nats.Msg) server.Order, error {
	var order server.Order
	err := json.Unmarshal(msg.Data, &order)
	
	if err != nil {
        return nil, err
    } else {
		return order, nil
	}
}