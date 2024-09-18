package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"

	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/api/pb"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/rabbitmq"
	"google.golang.org/protobuf/proto"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/sender_config.yaml",
		"Path to configuration file")
}

func main() {
	flag.Parse()

	cfg := config.ReadConfig(configFile)
	log := logger.New(cfg.Logger)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	rq, err := rabbitmq.NewQueue(ctx, cfg.AMQP, log)
	if err != nil {
		log.Error(err)
		cancel()
	}
	defer rq.Conn.Close()
	defer rq.Ch.Close()

	msgs, err := rq.Ch.Consume(
		rq.Name, // queue
		"",      // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		log.Error(err)
		cancel()
	}

	go func() {
		for d := range msgs {
			p := &pb.Notification{}
			if err := proto.Unmarshal(d.Body, p); err != nil {
				log.Errorf("Failed to parse Person %v", err)
			}
			log.Infof("Received a message %v", p)
		}
	}()

	log.Info(" [*] Waiting for messages. To exit press CTRL+C")

	<-ctx.Done()
}
