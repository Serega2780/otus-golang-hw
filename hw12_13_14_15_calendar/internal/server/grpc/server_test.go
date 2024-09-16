package grpc

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/api/pb"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/config"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/model"
	service "github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/service/memory"
	"github.com/Serega2780/otus-golang-hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_Events(t *testing.T) {
	ctx := context.Background()
	l := logger.New(&config.LoggerConf{Level: "info", Format: "json", LogToFile: false, LogToConsole: true})
	repo := memory.New()
	service := service.NewEventMemoryService(repo)

	client, closer := server(ctx, l, service)
	defer closer()

	t.Run("not found", func(t *testing.T) {
		out, err := client.FindEvent(ctx, &pb.FindEventRequest{Id: "c39d8763-34ad-4759-954f-b288649dac73"})
		require.NoError(t, err)
		require.Nil(t, out.Event)
	})
	t.Run("insert", func(t *testing.T) {
		out, err := client.CreateEvent(ctx, getCreateEventRequest())
		require.NoError(t, err)
		require.Greater(t, len(out.Event.Id), 0)

		_, err = client.RemoveEvent(ctx, &pb.RemoveEventRequest{Id: out.Event.Id})
		require.NoError(t, err)
	})
	t.Run("update", func(t *testing.T) {
		out, err := client.CreateEvent(ctx, getCreateEventRequest())
		require.NoError(t, err)
		event := out.Event
		event.Title = "UPDATED"
		out, err = client.UpdateEvent(ctx, &pb.UpdateEventRequest{Event: event})
		require.NoError(t, err)
		require.Equal(t, "UPDATED", out.Event.Title)

		_, err = client.RemoveEvent(ctx, &pb.RemoveEventRequest{Id: out.Event.Id})
		require.NoError(t, err)
	})

	t.Run("findAll", func(t *testing.T) {
		e := getCreateEventRequest()
		e2 := getCreateEventRequest()
		e3 := getCreateEventRequest()
		e2.Event.StartTime = timestamppb.New(e2.Event.StartTime.AsTime().Add(2 * time.Hour))
		e2.Event.EndTime = timestamppb.New(e2.Event.StartTime.AsTime().Add(1 * 30 * time.Minute))
		e3.Event.StartTime = timestamppb.New(e3.Event.StartTime.AsTime().Add(3 * time.Hour))
		e3.Event.EndTime = timestamppb.New(e3.Event.StartTime.AsTime().Add(1 * 30 * time.Minute))
		_, _ = client.CreateEvent(ctx, e)
		_, _ = client.CreateEvent(ctx, e2)
		_, _ = client.CreateEvent(ctx, e3)

		out, err := client.FindEvents(ctx, &emptypb.Empty{})
		require.NoError(t, err)
		require.Equal(t, 3, len(out.Events))
	})
}

func server(ctx context.Context, l *logger.Logger, s *service.EventMemoryService) (pb.EventServiceClient, func()) {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)
	h := NewEventsHandler(l, s)
	baseServer := grpc.NewServer()
	pb.RegisterEventServiceServer(baseServer, h)
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "", //nolint:staticcheck
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server: %v", err)
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		baseServer.Stop()
	}

	client := pb.NewEventServiceClient(conn)

	return client, closer
}

func getEvent() model.Event {
	return model.Event{
		Title:       "4th event",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(1 * time.Hour),
		Description: "long description for 4th event",
		UserID:      uuid.New().String(),
	}
}

func getCreateEventRequest() *pb.CreateEventRequest {
	event := getEvent()
	return &pb.CreateEventRequest{
		Event: &pb.EventNew{
			Title:       event.Title,
			StartTime:   timestamppb.New(event.StartTime),
			EndTime:     timestamppb.New(event.EndTime),
			Description: event.Description,
			UserId:      event.UserID,
		},
	}
}
