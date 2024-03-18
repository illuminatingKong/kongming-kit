package examplehttpguzzle

import (
	"fmt"
	"github.com/illuminatingKong/kongming-kit/http/client/samplehttp"
	"testing"
)

func TestGet(t *testing.T) {
	s, err := samplehttp.NewProvider("https", "httpbin.org")
	if err != nil {
		panic(err)
	}
	type dome struct {
		Origin string `json:"origin"`
	}
	req := s.Do("GET", "/ip", &samplehttp.Option{})
	out := req.Json(&dome{})
	fmt.Printf("response: %+v", out)
}
