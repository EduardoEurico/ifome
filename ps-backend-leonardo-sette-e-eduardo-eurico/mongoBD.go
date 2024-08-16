package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

var (
	Client *mongo.Client
	DB     *mongo.Database
)

func connectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("Erro ao fazer ping no MongoDB: %v", err)
	}

	Client = client
	DB = client.Database("test")

	log.Println("Conexão ao MongoDB estabelecida com sucesso.")
}

func disconnectDB() {
	if err := Client.Disconnect(context.TODO()); err != nil {
		log.Fatalf("Erro ao desconectar do MongoDB: %v", err)
	}
	fmt.Println("Desconectado do MongoDB!")
}
