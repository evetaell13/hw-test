package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/evetaell13/hw-test/hw12_13_14_15_calendar/internal/app"
	"github.com/evetaell13/hw-test/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/evetaell13/hw-test/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/evetaell13/hw-test/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/evetaell13/hw-test/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig(configFile)
	logg := logger.New(config.Logger.Level)

	var storage app.Storage

	switch strings.ToLower(config.Database.DBimplement) {
	case "inmemory":
		d, err := memorystorage.New(config.Database.FilePath)
		if err != nil {
			logg.Error(fmt.Sprint("create inmemory: ", err))
			os.Exit(1)
		}
		storage = d
	case "pg":
		db, err := sqlstorage.New() // TODO args
		if err != nil {
			logg.Error(fmt.Sprint("create sqlstorage: ", err))
			os.Exit(1)
		}
		storage = db
	default:
		logg.Error(fmt.Sprint("unsupport DBimplement: ", config.Database.DBimplement))
		os.Exit(1)
	}
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
