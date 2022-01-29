package rmq

import (
	"banners-rotator/internal/logger"
	"context"
	"fmt"
)

type Consumer struct {
	name   string
	conn   Connection
	logger *logger.Logger
}

func NewRMQConsumer(name string, conn Connection, logg *logger.Logger) *Consumer {
	return &Consumer{
		name:   name,
		conn:   conn,
		logger: logg,
	}
}

type Message struct {
	Ctx  context.Context
	Data []byte
}

func (c *Consumer) Consume(ctx context.Context, queue string) (<-chan Message, error) {
	messages := make(chan Message)

	ch, err := c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("rmq open channel: %w", err)
	}

	go func() {
		<-ctx.Done()
		if err := ch.Close(); err != nil {
			c.logger.Error("rmq close channel: " + err.Error())
		}
	}()

	deliveries, err := ch.Consume(
		queue,
		c.name,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("rmq start consuming: %w", err)
	}

	go func() {
		defer func() {
			close(messages)
			c.logger.Info("rmq close messages channel")
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case del := <-deliveries:
				if err := del.Ack(false); err != nil {
					c.logger.Error("rmq deliver message: " + err.Error())
				}

				msg := Message{
					Ctx:  context.Background(),
					Data: del.Body,
				}

				select {
				case <-ctx.Done():
					return
				case messages <- msg:
				}
			}
		}
	}()

	return messages, nil
}
