package sender

import (
	"encoding/json"
	"fmt"

	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/config"
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/logger"
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/streadway/amqp"
)

type Sender struct {
	config          config.RabbitConf
	rabbit *rabbitmq.RabbitClient
	log          *logger.Logger
}

func NewSender(config config.RabbitConf, rabbit *rabbitmq.RabbitClient, logg *logger.Logger) *Sender {
	return &Sender{
		config: config,
		rabbit: rabbit,
		log: logg,
	}
}

func (s *Sender) Run() error {
	if err := s.rabbit.ExchangeDeclare(s.config.GetExchange()); err != nil {
		return err
	}

	if err := s.rabbit.QueueDeclare(s.config.GetQueue()); err != nil {
		return err
	}

	if err := s.rabbit.QueueBind(s.config.GetExchange()); err != nil {
		return err
	}

	for {
		messageCh, err := s.rabbit.Consume()
		if err != nil {
			return err
		}

		s.log.Info("sender is running...")

		for message := range messageCh {
			if err := s.processMessage(message); err != nil {
				s.log.Error(err.Error())
				message.Nack(false, true)
				continue
			}

			message.Ack(false)
		}
	}
}

func (s *Sender) processMessage(msg amqp.Delivery) error {
	var event storage.Event
	
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		return err
	}

	s.log.Info(fmt.Sprintf("event is sended: %v", event))

	return nil
}