package app

import (
	"context"
	"fmt"
	"os"

	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/service"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/service/memory"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/service/sql"
	memorystorage "github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/repository"
	sqlstorage "github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/sql"
	"gopkg.in/yaml.v3"
)

type App struct { // TODO
	Cfg           *config.Config
	Log           *logger.Logger
	Storage       repository.EventInterface
	EventsService service.EventService
}

func Start(ctx context.Context, configFile string) *App {
	cfg := readConfig(configFile)
	log := logger.New(cfg.Logger)
	app := &App{Cfg: cfg, Log: log}
	if cfg.IsInMemory {
		app.Storage = memorystorage.New()
		app.EventsService = memory.NewEventMemoryService(app.Storage.(*memorystorage.Storage))
	} else {
		app.Storage = sqlstorage.New(cfg.DB)
	}
	if !cfg.IsInMemory {
		err := app.Storage.(*sqlstorage.Storage).Connect(ctx)
		if err != nil {
			app.Log.Errorf("failed to connect to Postgres %v", err)
			os.Exit(1)
		}
		if len(cfg.DB.Migration) != 0 {
			err := app.Storage.(*sqlstorage.Storage).Migrate(ctx)
			if err != nil {
				app.Log.Errorf("failed to migrate db %v", err)
				os.Exit(1)
			}
		}
		app.EventsService = sql.NewEventSQLService(app.Storage.(*sqlstorage.Storage))
	}
	app.Log.Info("calendar is running...")
	return app
}

func readConfig(configFile string) *config.Config {
	conf := &config.Config{}
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Could not open a config file %v\n", err)
	} else {
		err = yaml.Unmarshal(yamlFile, conf)
		if err != nil {
			fmt.Printf("Could not unmarshal a config file %v\n", err)
		}
	}
	if conf.Logger == nil {
		fmt.Println("Will use default config")
		conf = config.New()
	}
	return conf
}
