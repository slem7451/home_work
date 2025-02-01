package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"

	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/app"                              //nolint:depguard
	schedulerconfig "github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/config/scheduler" //nolint:depguard
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/logger"                           //nolint:depguard
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/rabbitmq"                         //nolint:depguard
	schedulerlib "github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/scheduler"           //nolint:depguard
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/scheduler_config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config := schedulerconfig.NewConfig(configFile)
	logg := logger.New(config.Logger.Level)

	storage := app.NewStorage(config)

	rabbit := rabbitmq.New(config.Rabbit, logg)

	if err := rabbit.Connect(); err != nil {
		panic(err)
	}

	defer rabbit.Close()

	scheduler := schedulerlib.NewScheduler(config.Rabbit, config.Scheduler, rabbit, storage, logg)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	if err := scheduler.Run(ctx); err != nil {
		panic(err)
	}
}
