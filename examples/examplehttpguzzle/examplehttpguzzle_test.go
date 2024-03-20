package examplehttpguzzle

import (
	"fmt"
	"github.com/illuminatingKong/kongming-kit/http/client/samplehttp"
	"testing"
	"time"
)

type testResponse struct {
	Code int `json:"code"`
	Data struct {
		Now     time.Time `json:"now"`
		Version string    `json:"version"`
	} `json:"data"`
	Extra   interface{} `json:"extra"`
	Limit   int         `json:"limit"`
	Message string      `json:"message"`
	Page    int         `json:"page"`
	Total   int         `json:"total"`
}

func TestGet(t *testing.T) {
	s, err := samplehttp.NewProvider("http", "localhost:8081")
	s.SetDebug(true)
	if err != nil {
		panic(err)
	}
	f := &testResponse{}
	out := s.Do("GET", "/health/core/version", &samplehttp.Option{})
	resp, err := s.Json(out, f)
	if err != nil {
		panic(err)
	}
	fmt.Printf("response: %+v", resp)
}
