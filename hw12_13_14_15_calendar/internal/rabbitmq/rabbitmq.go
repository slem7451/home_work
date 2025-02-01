package rabbitmq

import (
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/config" //nolint:depguard
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/logger" //nolint:depguard
	"github.com/streadway/amqp"                                            //nolint:depguard
)

type RabbitClient struct {
	url       string
	queueName string
	log       *logger.Logger
	conn      *amqp.Connection
	channel   *amqp.Channel
}

type Config struct {
	Port     int
	Host     string
	User     string
	Password string
	Exchange string
	Queue    string
}

func New(config config.RabbitConf, log *logger.Logger) *RabbitClient {
	return &RabbitClient{
		url: config.URL(),
		log: log,
	}
}

func (r *RabbitClient) Connect() error {
	conn, err := amqp.Dial(r.url)
	if err != nil {
		return err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return err
	}

	r.conn = conn
	r.channel = channel

	return nil
}

func (r *RabbitClient) Close() {
	if r.channel != nil {
		r.channel.Close()
	}

	if r.conn != nil {
		r.conn.Close()
	}
}

func (r *RabbitClient) ExchangeDeclare(exchangeName string) error {
	return r.channel.ExchangeDeclare(exchangeName, "fanout", true, false, false, false, nil)
}

func (r *RabbitClient) QueueDeclare(queueName string) error {
	_, err := r.channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	r.queueName = queueName

	return nil
}

func (r *RabbitClient) QueueBind(exchangeName string) error {
	return r.channel.QueueBind(r.queueName, "", exchangeName, false, nil)
}

func (r *RabbitClient) Publish(exchangeName string, message []byte) error {
	qMessage := amqp.Publishing{
		ContentType: "application/json",
		Body:        message,
	}

	return r.channel.Publish(exchangeName, r.queueName, false, false, qMessage)
}

func (r *RabbitClient) Consume() (<-chan amqp.Delivery, error) {
	return r.channel.Consume(r.queueName, "", false, false, false, false, nil)
}
