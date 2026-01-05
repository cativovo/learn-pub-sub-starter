package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// PublishJSON marshals the value and publishes it to exchange in the server.
func PublishJSON[T any](ctx context.Context, ch *amqp.Channel, exchange string, key string, val T) error {
	b, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("pubsub: marshal val: %w", err)
	}

	err = ch.PublishWithContext(
		ctx,
		exchange,
		key,
		false,
		false,
		amqp.Publishing{ContentType: "application/json", Body: b},
	)
	if err != nil {
		return fmt.Errorf("pubsub: publish: %w", err)
	}

	return nil
}
