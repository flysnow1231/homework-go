package mq

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	ch         *amqp091.Channel
	exchange   string
	routingKey string
}

func NewProducer(ch *amqp091.Channel, exchange, routingKey string) *Producer {
	return &Producer{ch: ch, exchange: exchange, routingKey: routingKey}
}

func (p *Producer) PublishJSON(ctx context.Context, typ string, traceID string, payload []byte) error {
	if p == nil || p.ch == nil {
		return nil
	}
	headers := amqp091.Table{
		"type":     typ,
		"trace_id": traceID,
		"id":       uuid.NewString(),
		"ts":       time.Now().Unix(),
	}
	msg := amqp091.Publishing{
		ContentType: "application/json",
		Body:        payload,
		Headers:     headers,
		Timestamp:   time.Now(),
	}
	return p.ch.PublishWithContext(ctx, p.exchange, p.routingKey, false, false, msg)
}
