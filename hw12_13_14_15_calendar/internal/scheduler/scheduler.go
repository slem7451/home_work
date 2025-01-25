package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/config"
	schedulerconfig "github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/config/scheduler"
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/logger"
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/storage"
)

type Scheduler struct {
	config          config.RabbitConf
	schedulerConf schedulerconfig.SchedulerConf
	rabbit *rabbitmq.RabbitClient
	storage      app.Storage
	log          *logger.Logger
}

func NewScheduler(config config.RabbitConf, schedulerConf schedulerconfig.SchedulerConf,  rabbit *rabbitmq.RabbitClient, storage app.Storage, log *logger.Logger) *Scheduler {
	return &Scheduler{
		config:       config,
		schedulerConf: schedulerConf,
		rabbit: rabbit,
		storage:      storage,
		log:          log,
	}
}

func (s *Scheduler) Run(ctx context.Context) error {
	if err := s.rabbit.ExchangeDeclare(s.config.GetExchange()); err != nil {
		return err
	}
	
	if err := s.rabbit.QueueDeclare(s.config.GetQueue()); err != nil {
		return err
	}

	notifyDuration, err := time.ParseDuration(s.schedulerConf.Update)
	if err != nil {
		return err
	}

	notifyTicker := time.NewTicker(notifyDuration)
	defer notifyTicker.Stop()

	removeDuration, err := time.ParseDuration(s.schedulerConf.Remove)
	if err != nil {
		return err
	}

	removeTicker := time.NewTicker(removeDuration)
	defer removeTicker.Stop()

	s.log.Info("scheduler is running...")

	for {
		select {
		case <-notifyTicker.C:
			events, err := s.storage.FindEventsForNotify(ctx)
			if err != nil {
				s.log.Error(err.Error())
				continue
			}

			if err := s.sendNotifications(ctx, events); err != nil {
				s.log.Error(err.Error())
				continue
			}
		case <-removeTicker.C:
			if err := s.storage.RemoveOldEvents(ctx); err != nil {
				s.log.Error(err.Error())
				continue
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (s *Scheduler) sendNotifications(ctx context.Context, events []storage.Event) error {
	for _, event := range events {
		body, err := json.Marshal(event)
		if err != nil {
			return err
		}
		
		if err := s.rabbit.Publish(s.config.GetExchange(), body); err != nil {
			return err
		}

		s.log.Info(fmt.Sprintf("event %d is in scheduler now", event.ID))

		s.storage.MarkSendedEvent(ctx, event.ID)
	}

	return nil
}