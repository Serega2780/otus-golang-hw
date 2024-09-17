package rabbitmq

import (
	"fmt"

	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/api/pb"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

const (
	ContentType      = "text/plain"
	RabbitMessageTTL = 31560000000
)

func (q *Queue) Send(n *pb.Notification) error {
	fail := func(m string, e error) error {
		return fmt.Errorf("%s %w", m, e)
	}
	out, err := proto.Marshal(n)
	if err != nil {
		return fail("Failed to encode event:", err)
	}
	body := out
	err = q.Ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: ContentType,
			Body:        body,
		})
	if err != nil {
		return fail("Failed to publish a message:", err)
	}
	return nil
}
