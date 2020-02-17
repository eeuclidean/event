package gokit

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

func RequestDuration() metrics.Histogram {
	duration := prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "project",
		Subsystem: "booking",
		Name:      "request_duration_seconds",
		Help:      "Request duration in seconds.",
	}, []string{"method", "success"})
	return duration
}

func InstrumentingLatencyMiddleware(duration metrics.Histogram) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				if err == nil {
					duration.With("success", fmt.Sprint(err == nil)).Observe(time.Since(begin).Seconds())
				}
			}(time.Now())
			return next(ctx, request)
		}
	}
}
