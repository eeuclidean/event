package gokit

import (
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

func NewHTTPServer(ctx context.Context, endpoints Endpoints) http.Handler {
	r := mux.NewRouter()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}
	r.Methods("POST").Path("/api/v1/user/create").Handler(httptransport.NewServer(
		endpoints.CreateBookingEndPoint,
		decodeCreateBookingRequest,
		encodeCreateBookingResponse,
		options...,
	))
	r.Methods("PUT").Path("/api/v1/user/call").Handler(httptransport.NewServer(
		endpoints.CallBookingEndPoint,
		decodeCallBookingRequest,
		encodeSuccessResponse,
		options...,
	))
	r.Methods("PUT").Path("/api/v1/user/pay").Handler(httptransport.NewServer(
		endpoints.PayBookingEndPoint,
		decodePayBookingRequest,
		encodeSuccessResponse,
		options...,
	))
	r.Methods("PUT").Path("/api/v1/user/cancel").Handler(httptransport.NewServer(
		endpoints.CancelBookingEndPoint,
		decodeCancelBookingRequest,
		encodeSuccessResponse,
		options...,
	))
	r.Methods("PUT").Path("/api/v1/user/checkin_loket").Handler(httptransport.NewServer(
		endpoints.CheckinLoketBookingEndPoint,
		decodeCheckinLoketBookingRequest,
		encodeSuccessResponse,
		options...,
	))
	r.Methods("PUT").Path("/api/v1/user/checkin_poli").Handler(httptransport.NewServer(
		endpoints.CheckinPoliBookingEndPoint,
		decodeCheckinPoliRequest,
		encodeSuccessResponse,
		options...,
	))
	r.Methods("PUT").Path("/api/v1/user/antrian").Handler(httptransport.NewServer(
		endpoints.UpdateAntrianEndPoint,
		decodeUpdateAntrianRequest,
		encodeSuccessResponse,
		options...,
	))
	return r
}
