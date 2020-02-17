package gokit

import (
	"event/user/service"
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type errorer interface {
	error() error
}

type bodyResponse struct {
	Message interface{} `json:"message,omitempty"`
}

func encodeSuccessResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(bodyResponse{
		Message: "success",
	})
}

func encodeCreateBookingResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(bodyResponse{
		Message: response,
	})
}

func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})

}

func codeFrom(err error) int {
	if strings.Contains(err.Error(), service.ERR_PREFIX) {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
