package bootstrap

import (
	"time"

	"blog/internal/config"
	"github.com/rabbitmq/amqp091-go"
)

type Rabbit struct {
	Conn *amqp091.Connection
	Ch   *amqp091.Channel
}

func InitRabbitMQ(cfg config.RabbitMQConfig) (*Rabbit, error) {
	conn, err := amqp091.Dial(cfg.URL)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}

	if err := ch.ExchangeDeclare(cfg.Exchange, "topic", true, false, false, false, nil); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, err
	}
	_, err = ch.QueueDeclare(cfg.Queue, true, false, false, false, nil)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, err
	}
	if err := ch.QueueBind(cfg.Queue, cfg.RoutingKey, cfg.Exchange, false, nil); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, err
	}
	return &Rabbit{Conn: conn, Ch: ch}, nil
}

func PingRabbitMQ(cfg config.RabbitMQConfig) error {
	dialer := amqp091.Config{Dial: amqp091.DefaultDial(3 * time.Second)}
	conn, err := amqp091.DialConfig(cfg.URL, dialer)
	if err != nil {
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return err
	}
	_ = ch.Close()
	_ = conn.Close()
	return nil
}
