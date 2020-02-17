package benchmark_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"event/user/commands"
	"event/user/gokit"
	"event/user/service"
	"event/user/testing/config"
	"event/user/utils/utilsgenerator"
	"testing"
)

var httpHandler http.Handler

func init() {
	config.RunConfig()
	svc, err := service.NewService()
	if err != nil {
		panic(err)
	}
	endpoints := gokit.NewEndPoints(svc, gokit.Logger(), gokit.RequestDuration())
	httpHandler = gokit.NewHTTPServer(context.Background(), endpoints)
}

func BenchmarkAddBookingAPI(b *testing.B) {
	for n := 0; n < b.N; n++ {
		addBookingCommand := commands.AddBookingCommand{
			PatientID:   utilsgenerator.NewID(),
			BranchID:    utilsgenerator.NewID(),
			PoliID:      utilsgenerator.NewID(),
			SubPoliID:   utilsgenerator.NewID(),
			InsuranceID: utilsgenerator.NewID(),
			Tanggal:     "12-12-2019",
			CreateBy:    utilsgenerator.NewID(),
		}
		data, _ := json.Marshal(addBookingCommand)
		_, err := http.NewRequest("POST", "/api/v1/user/create", bytes.NewReader(data))
		if err != nil {
			panic(err)
		}
	}
}
