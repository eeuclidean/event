package main

import (
	"event/user/event"
	"event/user/event/handlers"
	"event/user/gokit"
	"event/user/service"
	"context"
	"flag"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Log(functionName string, msg string) {
	log.With(gokit.Logger(), "method", functionName).Log("msg", msg, "took", 0)
}

func startEventConsumer(svc service.Service, ctx context.Context) error {

	eventConsumers, err := event.NewRedisEventConsumer(handlers.NewEventHandlers(svc, Log), Log)
	if err != nil {
		Log("NewRedisEventConsumer", err.Error())
		return err
	}
	err = eventConsumers.Run(ctx)
	if err != nil {
		Log("EventConsumer Run", err.Error())
		return err
	}
	return nil
}

func startMetricServer() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":9090", nil)
		if err := http.ListenAndServe(":9090", nil); err != nil {
			Log("http.ListenAndServe(\":9090\", nil)", err.Error())
			return
		}
	}()
}

func startServer(svc service.Service) {
	var server *http.Server
	{
		endpoints := gokit.NewEndPoints(svc, gokit.Logger(), gokit.RequestDuration())
		httpHandlers := gokit.NewHTTPServer(context.Background(), endpoints)
		httpAddr := flag.String("http.addr", ":3210", "HTTP listen address")
		flag.Parse()
		server = &http.Server{
			Addr:    *httpAddr,
			Handler: httpHandlers,
		}
	}
	if err := server.ListenAndServe(); err != nil {
		Log("server.ListenAndServe()", err.Error())
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	svc, err := service.NewService()
	if err != nil {
		Log("NewService", err.Error())
	}
	if err := startEventConsumer(svc, ctx); err != nil {
		return
	}
	startMetricServer()
	startServer(svc)
}
