package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	connStr := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connStr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	_ = conn

	fmt.Println("Successfully connected to RabbitMQ")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	<-ctx.Done()

	fmt.Println("Shutting down the server...")
}
