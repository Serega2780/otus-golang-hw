package main

import (
	"context"
	"flag"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/api/pb"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/mapper"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/rabbitmq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	configFile string
	period     int
)

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/scheduler_config.yaml",
		"Path to configuration file")
}

func main() {
	flag.Parse()

	var opts []grpc.DialOption
	cfg := config.ReadConfig(configFile)
	log := logger.New(cfg.Logger)
	period = cfg.SchedulerPeriodSeconds

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

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.NewClient(strings.Join([]string{cfg.GRPC.IP, cfg.GRPC.Port}, ":"), opts...)
	if err != nil {
		log.Error(err)
		cancel()
	}
	client := pb.NewEventServiceClient(conn)

	ticker := time.NewTicker(time.Duration(period) * time.Second)

	log.Info(" [*] Scanning for notifications. To exit press CTRL+C")

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			eRes, err := client.FindForNotify(ctx, &emptypb.Empty{})
			if err != nil {
				log.Error(err)
			} else if eRes != nil && len(eRes.Events) > 0 {
				notifies := mapper.Notifications(eRes)
				processNotifications(ctx, rq, client, log, notifies)
			}
		}
	}
}

func processNotifications(ctx context.Context, rq *rabbitmq.Queue, client pb.EventServiceClient, log *logger.Logger,
	notifies []*pb.Notification,
) {
	for _, n := range notifies {
		err := rq.Send(n)
		if err != nil {
			log.Error(err)
		}
		log.Infof("Sent notify %v", n)
		_, err = client.SetNotified(ctx, &pb.SetNotifiedRequest{Id: n.Id})
		if err != nil {
			log.Error(err)
		}
	}
}
