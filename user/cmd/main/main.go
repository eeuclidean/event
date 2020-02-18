package main

import (
	"context"
	"event/user/event/consumer"
	"event/user/gokit"
	"event/user/service"
	"flag"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	svc, err := service.NewService()
	if err != nil {
		Log("NewService", err.Error())
		return
	}
	if err := startEventConsumer(svc, ctx); err != nil {
		return
	}
	startMetricServer()
	startServer(svc)
}

func Logger() log.Logger {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)
	logger = log.With(logger, "namespace", "event")
	logger = log.With(logger, "service", "user")
	return logger
}

func Histogram() metrics.Histogram {
	duration := prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "project",
		Subsystem: "booking",
		Name:      "request_duration_seconds",
		Help:      "Request duration in seconds.",
	}, []string{"method", "success"})
	return duration
}

func Log(functionName string, msg string) {
	log.With(Logger(), "method", functionName).Log("msg", msg, "took", 0)
}

func startEventConsumer(svc service.Service, ctx context.Context) error {
	eventConsumers, err := consumer.NewEventConsumer(svc, Log)
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
		endpoints := gokit.NewEndPoints(svc, Logger(), Histogram())
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
