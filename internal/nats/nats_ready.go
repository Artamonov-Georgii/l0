package nats_checker

import (
	"encoding/json"
	"errors"
	"github.com/nats-io/stan.go"
    "github.com/Artamonov-Georgii/l0/internal/server" 
)

func IsValidOrderMsg(msg *stan.Msg) (server.Order, error) {
	var order server.Order
	err := json.Unmarshal(msg.Data, &order)
	
	if err != nil {
        return order, err
    } else if order.OrderUID != "" {
		return order, nil
	} else {
		return order, errors.New("Not a valid order")
	}
}