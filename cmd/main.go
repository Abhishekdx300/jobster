package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Abhishekdx300/jobster/internal/api"
	"github.com/Abhishekdx300/jobster/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg := config.LoadConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %v", err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatalf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Could not ping MongoDB: %v", err)
	}
	log.Println("Successfully connected to MongoDB!")

	db := client.Database(cfg.DbName)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: api.RegisterRoutes(db),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
