package examplehttpguzzle

import (
	"fmt"
	"github.com/illuminatingKong/kongming-kit/http/client/samplehttp"
	"testing"
)

func TestGet(t *testing.T) {
	s, err := samplehttp.New("https", "httpbin.org")
	if err != nil {
		panic(err)
	}
	r := s.Get("/ip")
	fmt.Printf("response: %+v", r)
}
