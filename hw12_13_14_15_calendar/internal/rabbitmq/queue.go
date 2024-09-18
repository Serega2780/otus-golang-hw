package rabbitmq

import (
	"context"
	"fmt"

	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue struct {
	ctx  context.Context
	log  *logger.Logger
	IP   string
	Port string
	url  string
	Name string
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

func NewQueue(ctx context.Context, conf *config.AMQPConfig, log *logger.Logger) (*Queue, error) {
	fail := func(m string, e error) (*Queue, error) {
		return nil, fmt.Errorf("%s %w", m, e)
	}
	args := amqp.Table{
		amqp.QueueMessageTTLArg: RabbitMessageTTL,
	}
	queue := &Queue{
		ctx: ctx, log: log, IP: conf.IP, Port: conf.Port,
		url: "amqp://guest:guest@" + conf.IP + ":" + conf.Port + "/", Name: conf.QueueName,
	}

	conn, err := amqp.Dial(queue.url)
	if err != nil {
		return fail("Failed to connect to RabbitMQ", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return fail("Failed to open a channel", err)
	}

	_, err = ch.QueueDeclare(
		queue.Name, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		args,       // arguments
	)
	if err != nil {
		return fail("Failed to declare a queue", err)
	}
	queue.Conn = conn
	queue.Ch = ch
	return queue, nil
}
