package gokit

import (
	"context"
	"event/user/commands"
	"event/user/service"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

func NewEndPoints(svc service.Service, logger log.Logger, duration metrics.Histogram) Endpoints {
	var createBookingEndPoint endpoint.Endpoint
	{
		createBookingEndPoint = makeCreateBookingEndPoint(svc)
		createBookingEndPoint = loggingMiddleware(log.With(logger, "method", "CreateBooking"))(createBookingEndPoint)
		createBookingEndPoint = instrumentingLatencyMiddleware(duration.With("method", "CreateBooking"))(createBookingEndPoint)
	}
	var cancelBookingEndPoint endpoint.Endpoint
	{
		cancelBookingEndPoint = makeCancelBookingEndPoint(svc)
		cancelBookingEndPoint = loggingMiddleware(log.With(logger, "method", "CancelBooking"))(cancelBookingEndPoint)
		cancelBookingEndPoint = instrumentingLatencyMiddleware(duration.With("method", "CancelBooking"))(cancelBookingEndPoint)
	}
	var payBookingEndPoint endpoint.Endpoint
	{
		payBookingEndPoint = makePayBookingEndPoint(svc)
		payBookingEndPoint = loggingMiddleware(log.With(logger, "method", "PayBooking"))(payBookingEndPoint)
		payBookingEndPoint = instrumentingLatencyMiddleware(duration.With("method", "PayBooking"))(payBookingEndPoint)
	}
	var callBookingEndPoint endpoint.Endpoint
	{
		callBookingEndPoint = makeCallBookingEndPoint(svc)
		callBookingEndPoint = loggingMiddleware(log.With(logger, "method", "CallBooking"))(callBookingEndPoint)
		callBookingEndPoint = instrumentingLatencyMiddleware(duration.With("method", "CallBooking"))(callBookingEndPoint)
	}
	var loketCheckinBookingEndPoint endpoint.Endpoint
	{
		loketCheckinBookingEndPoint = makeCheckinLoketEndPoint(svc)
		loketCheckinBookingEndPoint = loggingMiddleware(log.With(logger, "method", "LoketCheckinBooking"))(loketCheckinBookingEndPoint)
		loketCheckinBookingEndPoint = instrumentingLatencyMiddleware(duration.With("method", "LoketCheckinBooking"))(loketCheckinBookingEndPoint)
	}
	var poliCheckinBookingEndPoint endpoint.Endpoint
	{
		poliCheckinBookingEndPoint = makeCheckinPoliEndPoint(svc)
		poliCheckinBookingEndPoint = loggingMiddleware(log.With(logger, "method", "PoliCheckinBooking"))(poliCheckinBookingEndPoint)
		poliCheckinBookingEndPoint = instrumentingLatencyMiddleware(duration.With("method", "PoliCheckinBooking"))(poliCheckinBookingEndPoint)
	}
	var updateAntrianEndPoint endpoint.Endpoint
	{
		updateAntrianEndPoint = makeUpdateAntrianEndPoint(svc)
		updateAntrianEndPoint = loggingMiddleware(log.With(logger, "method", "UpdateAntrian"))(updateAntrianEndPoint)
		updateAntrianEndPoint = instrumentingLatencyMiddleware(duration.With("method", "UpdateAntrian"))(updateAntrianEndPoint)
	}
	return Endpoints{
		CreateBookingEndPoint:       createBookingEndPoint,
		PayBookingEndPoint:          payBookingEndPoint,
		CancelBookingEndPoint:       cancelBookingEndPoint,
		CallBookingEndPoint:         callBookingEndPoint,
		CheckinLoketBookingEndPoint: loketCheckinBookingEndPoint,
		CheckinPoliBookingEndPoint:  poliCheckinBookingEndPoint,
		UpdateAntrianEndPoint:       updateAntrianEndPoint,
	}
}

type Endpoints struct {
	CreateBookingEndPoint       endpoint.Endpoint
	PayBookingEndPoint          endpoint.Endpoint
	CancelBookingEndPoint       endpoint.Endpoint
	CallBookingEndPoint         endpoint.Endpoint
	CheckinLoketBookingEndPoint endpoint.Endpoint
	CheckinPoliBookingEndPoint  endpoint.Endpoint
	UpdateAntrianEndPoint       endpoint.Endpoint
}

func makeCreateBookingEndPoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, rq interface{}) (interface{}, error) {
		request := rq.(commands.AddBookingCommand)
		return svc.CreateBooking(request)
	}
}

func makeCancelBookingEndPoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, rq interface{}) (interface{}, error) {
		request := rq.(commands.CancelBookingCommand)
		return nil, svc.CancelBooking(request)
	}
}

func makeCallBookingEndPoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, rq interface{}) (interface{}, error) {
		request := rq.(commands.CallBookingCommand)
		return nil, svc.CallBooking(request)
	}
}

func makePayBookingEndPoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, rq interface{}) (interface{}, error) {
		request := rq.(commands.PayBookingCommand)
		return nil, svc.PayBooking(request)
	}
}

func makeCheckinLoketEndPoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, rq interface{}) (interface{}, error) {
		request := rq.(commands.LoketCheckinBookingCommand)
		return nil, svc.CheckInLoketBooking(request)
	}
}

func makeCheckinPoliEndPoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, rq interface{}) (interface{}, error) {
		request := rq.(commands.PoliCheckinBookingCommand)
		return nil, svc.CheckInPoliBooking(request)
	}
}

func makeUpdateAntrianEndPoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, rq interface{}) (interface{}, error) {
		request := rq.(commands.UpdateAntrianCommand)
		return nil, svc.UpdateAntrian(request)
	}
}
