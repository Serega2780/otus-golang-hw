package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/repository"
	sqlstorage "github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/sql"
	"gopkg.in/yaml.v3"
)

var (
	configFile string
	storage    repository.EventInterface
	sqlStorage = new(sqlstorage.Storage)
)

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	cfg := &config.Config{}
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Could not open a config file %v\n", err)
	} else {
		err = yaml.Unmarshal(yamlFile, cfg)
		if err != nil {
			fmt.Printf("Could not unmarshal a config file %v\n", err)
		}
	}
	if cfg.Logger == nil {
		fmt.Println("Will use default config")
		cfg = config.New()
	}

	log := logger.New(cfg.Logger)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	if cfg.IsInMemory {
		storage = memorystorage.New()
	} else {
		storage = sqlstorage.New(cfg.DB)
		sqlStorage, _ = storage.(*sqlstorage.Storage)
	}

	calendar := app.New(log, storage)

	server := internalhttp.NewServer(cfg.HTTP, log, calendar)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		if err := server.Stop(ctx); err != nil {
			log.Errorf("failed to stop http server %v", err)
		}

		if sqlStorage != nil {
			if err := sqlStorage.Close(); err != nil {
				log.Errorf("failed to close Postgres connection %v", err)
			}
		}
	}()
	if !cfg.IsInMemory {
		if err := sqlStorage.Connect(ctx); err != nil {
			log.Errorf("failed to connect to Postgres %v", err)
			cancel()
			os.Exit(1) //nolint:gocritic
		}
		if len(cfg.DB.Migration) != 0 {
			if err := sqlStorage.Migrate(ctx); err != nil {
				log.Errorf("failed to migrate db %v", err)
				cancel()
				os.Exit(1)
			}
		}
	}
	log.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		log.Errorf("failed to start http server %v", err)
		cancel()
		os.Exit(1)
	}
}
