package gokit

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
)

func instrumentingLatencyMiddleware(duration metrics.Histogram) endpoint.Middleware {
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
