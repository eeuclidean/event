package gokit

import (
	"event/user/commands"
	"context"
	"encoding/json"
	"net/http"
)

func decodeCreateBookingRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req commands.AddBookingCommand
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeCancelBookingRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req commands.CancelBookingCommand
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodePayBookingRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req commands.PayBookingCommand
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeCallBookingRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req commands.CallBookingCommand
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeCheckinLoketBookingRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req commands.LoketCheckinBookingCommand
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
func decodeCheckinPoliRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req commands.PoliCheckinBookingCommand
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
func decodeUpdateAntrianRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req commands.UpdateAntrianCommand
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
