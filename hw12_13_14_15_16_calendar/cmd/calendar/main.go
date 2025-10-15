package main

import (
	"context"
	"errors"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DGKot/otus_hw/hw12_13_14_15_16_calendar/internal/app"
	"github.com/DGKot/otus_hw/hw12_13_14_15_16_calendar/internal/logger"
	internalhttp "github.com/DGKot/otus_hw/hw12_13_14_15_16_calendar/internal/server/http"
	memorystorage "github.com/DGKot/otus_hw/hw12_13_14_15_16_calendar/internal/storage/memory"
	sqlstorage "github.com/DGKot/otus_hw/hw12_13_14_15_16_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./../../configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig(configFile)
	logg := logger.New(logger.LogDeps{Level: config.Logger.Level})

	storage, err := createStorage(config.DB)
	if err != nil {
		logg.Error(err.Error())
		os.Exit(1)
	}

	calendar := app.New(logg, storage)

	serverDeps := &internalhttp.ServerDeps{
		Host:         config.Server.Host,
		Port:         config.Server.Port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	server := internalhttp.NewServer(logg, calendar, *serverDeps)

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
	logg.Info("END")
}

func createStorage(cfg DBConf) (app.Storage, error) {
	switch cfg.Type {
	case "memory":
		return memorystorage.New(), nil
	case "sql":
		return sqlstorage.New(sqlstorage.DBDeps{
			DSN:            cfg.DSN(),
			MigrationsPath: cfg.MigrationsPath,
		})
	default:
		return nil, errors.New("unknown db type")
	}
}
