package internalgrpc

import (
	"context"
	"errors"
	"strings"
	"time"

	internalserver "github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/server" //nolint:depguard
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/server/grpc/pb"        //nolint:depguard
	"github.com/slem7451/home_work/hw12_13_14_15_calendar/internal/storage"               //nolint:depguard
	"google.golang.org/protobuf/types/known/timestamppb"
)

type calendarService struct {
	pb.UnimplementedEventServiceServer
	app internalserver.Application
}

func validateEvent(event *pb.Event) error {
	err := ""
	var emptyDate time.Time

	if event.GetTitle() == "" {
		err += "title is required field\n"
	}

	if event.GetEventDate().AsTime() == emptyDate {
		err += "event_date is required field\n"
	}

	if event.GetDateSince().AsTime() == emptyDate {
		err += "date_since is required field\n"
	}

	if event.GetUserId() == 0 {
		err += "user_id is required field\n"
	}

	if err != "" {
		return errors.New(strings.Trim(err, "\n"))
	}

	return nil
}

func eventFromRequest(req *pb.Event) (storage.Event, error) {
	if err := validateEvent(req); err != nil {
		return storage.Event{}, err
	}

	return storage.Event{
		ID:         int(req.GetId()),
		Title:      req.GetTitle(),
		Descr:      req.GetDescr(),
		UserID:     int(req.GetUserId()),
		EventDate:  req.GetEventDate().AsTime(),
		DateSince:  req.GetDateSince().AsTime(),
		NotifyDate: req.GetNotifyDate().AsTime(),
	}, nil
}

func eventToResponse(event storage.Event) *pb.Event {
	return &pb.Event{
		Id:         int64(event.ID),
		Title:      event.Title,
		Descr:      event.Descr,
		UserId:     int64(event.UserID),
		EventDate:  timestamppb.New(event.EventDate),
		DateSince:  timestamppb.New(event.EventDate),
		NotifyDate: timestamppb.New(event.EventDate),
	}
}

func (s *calendarService) CreateEvent(ctx context.Context, req *pb.CreateEventRequest) (
	*pb.CreateEventResponse, error,
) {
	event, err := eventFromRequest(req.GetEvent())
	if err != nil {
		return &pb.CreateEventResponse{ //nolint:nilerr
			Resp: &pb.CreateEventResponse_Error{
				Error: err.Error(),
			},
		}, nil
	}

	id, err := s.app.CreateEvent(ctx, event)
	if err != nil {
		return &pb.CreateEventResponse{ //nolint:nilerr
			Resp: &pb.CreateEventResponse_Error{
				Error: err.Error(),
			},
		}, nil
	}

	return &pb.CreateEventResponse{
		Resp: &pb.CreateEventResponse_Id{
			Id: int64(id),
		},
	}, nil
}

func (s *calendarService) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (
	*pb.UpdateEventResponse, error,
) {
	event, err := eventFromRequest(req.GetEvent())
	if err != nil {
		return &pb.UpdateEventResponse{ //nolint:nilerr
			Error: err.Error(),
		}, nil
	}

	if err := s.app.UpdateEvent(ctx, int(req.GetEventId()), event); err != nil {
		return &pb.UpdateEventResponse{ //nolint:nilerr
			Error: err.Error(),
		}, nil
	}

	return &pb.UpdateEventResponse{}, nil
}

func (s *calendarService) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (
	*pb.DeleteEventResponse, error,
) {
	if err := s.app.DeleteEvent(ctx, int(req.GetEventId())); err != nil {
		return &pb.DeleteEventResponse{ //nolint:nilerr
			Error: err.Error(),
		}, nil
	}

	return &pb.DeleteEventResponse{}, nil
}

func (s *calendarService) FindEventsForDay(ctx context.Context, req *pb.FindEventsDateRequest) (
	*pb.FindEventsResponse, error,
) {
	events, err := s.app.FindEventsForDay(ctx, req.GetDate().AsTime())
	if err != nil {
		return nil, err
	}

	protoEvents := make([]*pb.Event, len(events))
	for i, event := range events {
		protoEvents[i] = eventToResponse(event)
	}

	return &pb.FindEventsResponse{Events: protoEvents}, nil
}

func (s *calendarService) FindEventsForWeek(ctx context.Context, req *pb.FindEventsDateRequest) (
	*pb.FindEventsResponse, error,
) {
	events, err := s.app.FindEventsForWeek(ctx, req.GetDate().AsTime())
	if err != nil {
		return nil, err
	}

	protoEvents := make([]*pb.Event, len(events))
	for i, event := range events {
		protoEvents[i] = eventToResponse(event)
	}

	return &pb.FindEventsResponse{Events: protoEvents}, nil
}

func (s *calendarService) FindEventsForMonth(ctx context.Context, req *pb.FindEventsDateRequest) (
	*pb.FindEventsResponse, error,
) {
	events, err := s.app.FindEventsForMonth(ctx, req.GetDate().AsTime())
	if err != nil {
		return nil, err
	}

	protoEvents := make([]*pb.Event, len(events))
	for i, event := range events {
		protoEvents[i] = eventToResponse(event)
	}

	return &pb.FindEventsResponse{Events: protoEvents}, nil
}

func (s *calendarService) FindEventsBetweenDates(ctx context.Context, req *pb.FindEventsBetweenRequest) (
	*pb.FindEventsResponse, error,
) {
	events, err := s.app.FindEventsBetweenDates(ctx, req.GetStart().AsTime(), req.GetEnd().AsTime())
	if err != nil {
		return nil, err
	}

	protoEvents := make([]*pb.Event, len(events))
	for i, event := range events {
		protoEvents[i] = eventToResponse(event)
	}

	return &pb.FindEventsResponse{Events: protoEvents}, nil
}
