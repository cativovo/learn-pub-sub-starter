package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")

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

	name, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get the client's name: %v", err)
		return
	}

	ch, queue, err := pubsub.DeclareAndBind(
		conn,
		"peril_direct",
		fmt.Sprintf("%s.%s", routing.PauseKey, name),
		routing.PauseKey,
		pubsub.SimpleQueueTypeTransient,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to declare and bind queue: %v", err)
		return
	}

	_ = ch
	_ = queue

	<-ctx.Done()

	fmt.Println("Shutting down the client...")
}
