package samplehttp

import (
	"bytes"
	"github.com/illuminatingKong/kongming-kit/base/logx"
	"github.com/illuminatingKong/kongming-kit/base/logx/logrusx"
	"github.com/illuminatingKong/kongming-kit/http/client/guzzle"
	"net/http"
)

type SampleProvider struct {
	Address    string
	Scheme     string
	HttpClient *guzzle.Client
	Log        logx.Logger
}

type ResponseWrapper struct {
	StatusCode int
	Body       string
	Header     http.Header
}

func New(scheme, address string) (*SampleProvider, error) {
	option := &guzzle.Config{
		Address:      address,
		Scheme:       scheme,
		ValidateHost: true,
	}

	httpClient, err := guzzle.NewClient(option)
	if err != nil {
		panic(err)
	}
	var formatter logrusx.JsonFormatter
	log := logrusx.New(logrusx.WithFormatter(formatter))
	return &SampleProvider{Address: address, Scheme: scheme, HttpClient: httpClient, Log: log}, nil

}

func (p *SampleProvider) Get(uri string) ResponseWrapper {
	r := p.HttpClient.DoNewRequest("GET", uri)
	r.SetParam("Header", "User-Agent", "sample-provider")
	timeout, httpResp, err := p.HttpClient.NewDoRequest(r)
	out := guzzle.RequireOK(timeout, httpResp, err)
	wrapper := ResponseWrapper{StatusCode: 0, Body: "", Header: make(http.Header)}
	wrapper.StatusCode = out.StatusCode
	buf := new(bytes.Buffer)
	buf.ReadFrom(out.RawResponse.Body)
	respBytes := buf.String()
	wrapper.Body = string(respBytes)
	wrapper.Header = out.Header
	return wrapper
}
