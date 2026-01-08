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

type SimpleQueueType int

const (
	SimpleQueueTypeDurable SimpleQueueType = iota
	SimpleQueueTypeTransient
)

// DeclareAndBind declares a queue and binds it to an exchange
func DeclareAndBind(
	conn *amqp.Connection,
	exchange string,
	queueName string,
	key string,
	queueType SimpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("pubsub: open channel: %w", err)
	}

	durable := queueType == SimpleQueueTypeDurable
	autoDelete := queueType == SimpleQueueTypeTransient
	exclusive := queueType == SimpleQueueTypeTransient

	queue, err := ch.QueueDeclare(
		queueName,
		durable,
		autoDelete,
		exclusive,
		false,
		nil,
	)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("pubsub: queue declare: %w", err)
	}

	err = ch.QueueBind(
		queueName,
		key,
		exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("pubsub: queue bind: %w", err)
	}

	return ch, queue, nil
}
