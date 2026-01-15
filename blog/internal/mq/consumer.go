package mq

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type HandlerFunc func(ctx context.Context, d amqp091.Delivery) error

type Consumer struct {
	ch          *amqp091.Channel
	queue       string
	consumerTag string
	prefetch    int
	concurrency int
	log         *zap.Logger
	handler     HandlerFunc

	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewConsumer(ch *amqp091.Channel, queue, consumerTag string, prefetch, concurrency int, log *zap.Logger, handler HandlerFunc) *Consumer {
	if prefetch <= 0 {
		prefetch = 20
	}
	if concurrency <= 0 {
		concurrency = 1
	}
	return &Consumer{
		ch:          ch,
		queue:       queue,
		consumerTag: consumerTag,
		prefetch:    prefetch,
		concurrency: concurrency,
		log:         log,
		handler:     handler,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	if c.ch == nil {
		return errors.New("nil channel")
	}
	if err := c.ch.Qos(c.prefetch, 0, false); err != nil {
		return err
	}
	deliveries, err := c.ch.Consume(
		c.queue,
		c.consumerTag,
		false, // auto-ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	runCtx, cancel := context.WithCancel(ctx)
	c.cancel = cancel

	for i := 0; i < c.concurrency; i++ {
		c.wg.Add(1)
		go func(worker int) {
			defer c.wg.Done()
			for {
				select {
				case <-runCtx.Done():
					return
				case d, ok := <-deliveries:
					if !ok {
						return
					}
					traceID, _ := d.Headers["trace_id"].(string)
					typ, _ := d.Headers["type"].(string)
					c.log.Info("mq_received",
						zap.Int("worker", worker),
						zap.String("type", typ),
						zap.String("trace_id", traceID),
						zap.Int("bytes", len(d.Body)),
					)

					err := c.handler(runCtx, d)
					if err != nil {
						c.log.Error("mq_handle_failed", zap.Error(err))
						_ = d.Nack(false, true) // requeue
						continue
					}
					_ = d.Ack(false)
				}
			}
		}(i + 1)
	}
	return nil
}

func (c *Consumer) Stop() {
	if c.cancel != nil {
		c.cancel()
	}
	c.wg.Wait()
}

func DefaultJSONHandler(log *zap.Logger) HandlerFunc {
	return func(ctx context.Context, d amqp091.Delivery) error {
		var m any
		if err := json.Unmarshal(d.Body, &m); err != nil {
			return err
		}
		log.Info("mq_payload", zap.Any("payload", m))
		return nil
	}
}
