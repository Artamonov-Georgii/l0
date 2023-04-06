package nats_checker

import (
	"encoding/json"
	"log"
	"fmt"
	"errors"
	"github.com/nats-io/stan.go"
    "github.com/Artamonov-Georgii/l0/internal/server" 
	"github.com/Artamonov-Georgii/l0/internal/pgsql" 
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

func MessageHandler(msg *stan.Msg) {
	log.Printf("Received message on subject \n")
		if order, error_mess := IsValidOrderMsg(msg); error_mess == nil {
			server.CacheMsg[order.OrderUID] = order
			err_inst := pgsql.InsertOrder(pgsql.Db, order) 
			if err_inst != nil {
                fmt.Println(err_inst)
			}	
			fmt.Printf("Received successfully \n")
		} else {
			fmt.Println("Message is not correct")
		}
}