package main

import (
	"context"
	"flag"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/server/http"
)

var (
	configFile string
	wg         sync.WaitGroup
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

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	calendar := app.Start(ctx, configFile)

	httpServer := http.NewServer(ctx, calendar)
	grpcServer := grpc.NewServer(ctx, calendar)

	wg.Add(1)
	go func() {
		defer wg.Done()
		httpServer.Start(ctx)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		grpcServer.Start(ctx)
	}()
	wg.Wait()

	if calendar.SQLStorage != nil {
		if err := calendar.SQLStorage.Close(); err != nil {
			calendar.Log.Errorf("failed to close Postgres connection %v", err)
		} else {
			calendar.Log.Info("Postgres connection closed")
		}
	}
}
