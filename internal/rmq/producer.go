package rmq

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

var ErrChanNotDeclared = errors.New("channel is not declared")

type Connection interface {
	Channel() (*amqp.Channel, error)
}

type QMessage struct {
	Type     string `json:"type"`
	SlotID   int64  `json:"slotId"`
	BannerID int64  `json:"bannerId"`
	GroupID  int64  `json:"groupId"`
	Date     int64  `json:"date"`
}

type Producer struct {
	name    string
	conn    Connection
	channel *amqp.Channel
}

func NewRMQProducer(name string, conn Connection) *Producer {
	return &Producer{name: name, conn: conn}
}

func (p *Producer) Connect() error {
	ch, err := p.conn.Channel()
	if err != nil {
		return fmt.Errorf("rmq get channel -> %w", err)
	}

	p.channel = ch

	_, err = ch.QueueDeclare(
		p.name,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("rmq declare queue -> %w", err)
	}

	return nil
}

func (p *Producer) Publish(message QMessage) error {
	if p.channel != nil {
		b, err := json.Marshal(message)
		if err != nil {
			return fmt.Errorf("rmq marshall message -> %w", err)
		}

		err = p.channel.Publish(
			"",     // exchange
			p.name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				MessageId:    uuid.NewString(),
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         b,
			})

		if err != nil {
			return fmt.Errorf("rmq publish message -> %w", err)
		}

		return nil
	}

	return ErrChanNotDeclared
}
