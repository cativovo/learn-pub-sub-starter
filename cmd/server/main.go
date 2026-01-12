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

	go func() {
		defer cancel()

		gamelogic.PrintServerHelp()

		for {
			input := gamelogic.GetInput()
			if len(input) == 0 {
				continue
			}

			switch input[0] {
			case "pause":
				fmt.Println("Sending pause message...")
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
					fmt.Println("failed to publish message,", err)
				}
			case "resume":
				fmt.Println("Sending pause resume...")
				err = pubsub.PublishJSON(
					ctx,
					ch,
					routing.ExchangePerilDirect,
					routing.PauseKey,
					routing.PlayingState{
						IsPaused: false,
					},
				)
				if err != nil {
					fmt.Println("Failed to publish message,", err)
				}
			case "quit":
				fmt.Println("Quitting...")
				return
			default:
				fmt.Println("Unable to process the command")
			}
		}
	}()

	<-ctx.Done()

	fmt.Println("Shutting down the server...")
}
