package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	connStr := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to RabbitMQ: %v", err)
		return
	}
	defer conn.Close()

	fmt.Println("Successfully connected to RabbitMQ")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open a RabbitMQ channel: %v", err)
		return
	}

	fmt.Println("Sending a message to the exchange")

	err = pubsub.PublishJSON(
		ctx,
		ch,
		routing.ExchangePerilDirect,
		routing.PauseKey,
		routing.PlayingState{
			IsPaused: true,
		},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to publish message: %v", err)
		return
	}

	fmt.Println("Message sent")

	<-ctx.Done()

	fmt.Println("Shutting down the server...")
}
