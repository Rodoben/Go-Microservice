package main

import (
	"context"
	"log"
	"loggerservice/data"
	"time"
)

type RPCServer struct{}

type RPCpayload struct {
	Name string
	Data string
}

func (r *RPCServer) LogInfo(payload RPCpayload, resp *string) error {

	log.Println("LOG INFO", "payload:", payload)
	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error writing to mongo", err)
		return err
	}

	*resp = "Processed payload via RPC:" + payload.Name
	return nil

}
