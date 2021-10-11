package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/devinmarder/go-qa-service/event"
	"github.com/devinmarder/go-qa-service/repository"
)

func Test_updateHandler(t *testing.T) {
	repo = &repository.LocalRepository{}

	eventChan := make(chan string)
	go event.RunEventProducer(eventChan)

	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name     string
		args     args
		expected string
	}{
		{
			"add test",
			args{
				httptest.NewRecorder(),
				httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`{"payload": {"service_name": "test", "coverage": 75}}`)),
			},
			"service name: test \ncoverage: 75",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateHandler(tt.args.w, tt.args.r)
			res := tt.args.w.Result()
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("expected error to be nil got %v", err)
			}
			if string(data) != tt.expected {
				t.Errorf("expected %v got %v", tt.expected, string(data))
			}
		})
	}
}
