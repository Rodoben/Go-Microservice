package main

import (
	"context"
	"fmt"
	"log"
	"loggerservice/data"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Model data.Models
}

func main() {

	log.Println("!11111111111111111111111111111111111111")

	mgoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	log.Println("MGOCLIENT:", mgoClient)
	client = mgoClient
	log.Println("!1111111111111111111222222222222222211111111111111")
	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	app := Config{Model: data.New(client)}

	err = rpc.Register(new(RPCServer))
	log.Println("!111111111111113333333333333311111111111111111111111", err)
	go app.rpcListen()
	go app.gRPCListen()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connectToMongo() (*mongo.Client, error) {
	log.Println("Connecting to mongo database")
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	log.Println("CLIENTOPTIONS:", clientOptions)

	c, err := mongo.Connect(context.TODO(), clientOptions)
	log.Println("CONNECTION:", c)
	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}
	log.Println("Connected to mongo!")
	return c, nil
}

func (app *Config) rpcListen() error {
	log.Println("Starting RPC Server on port:", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))

	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}

}
