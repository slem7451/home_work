package main

import (
	"flag"

	senderconfig "github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/config/sender" //nolint:depguard
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/logger"                     //nolint:depguard
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/rabbitmq"                   //nolint:depguard
	senderlib "github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/sender"           //nolint:depguard
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/sender_config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config := senderconfig.NewConfig(configFile)
	logg := logger.New(config.Logger.Level)

	rabbit := rabbitmq.New(config.Rabbit, logg)

	if err := rabbit.Connect(); err != nil {
		panic(err)
	}

	defer rabbit.Close()

	sender := senderlib.NewSender(config.Rabbit, rabbit, logg)

	if err := sender.Run(); err != nil {
		panic(err)
	}
}
