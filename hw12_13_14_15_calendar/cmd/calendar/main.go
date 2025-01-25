package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/app"                            //nolint:depguard
	calendarconfig "github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/config/calendar" //nolint:depguard
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/logger"                         //nolint:depguard
	serverbuilder "github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/server/builder"   //nolint:depguard
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/calendar_config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := calendarconfig.NewConfig(configFile)
	logg := logger.New(config.Logger.Level)

	storage := app.NewStorage(config)
	calendar := app.New(logg, storage)

	servers := serverbuilder.NewServers(logg, calendar, config)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	for _, server := range servers {
		go func() {
			<-ctx.Done()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()

			if err := server.Stop(ctx); err != nil {
				logg.Error(fmt.Sprintf("failed to stop %s server: %s", server.Whoami(), err.Error()))
			}
		}()
	}

	logg.Info("calendar is running...")

	for _, server := range servers {
		go func() {
			if err := server.Start(ctx); err != nil {
				logg.Error(fmt.Sprintf("failed to start %s server: %s", server.Whoami(), err.Error()))
				cancel()
				os.Exit(1)
			}
		}()
	}

	<-ctx.Done()
}
