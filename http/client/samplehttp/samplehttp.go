package samplehttp

import (
	"bytes"
	"context"
	"github.com/illuminatingKong/kongming-kit/base/logx"
	"github.com/illuminatingKong/kongming-kit/base/logx/logrusx"
	"github.com/illuminatingKong/kongming-kit/http/client/guzzle"
	"io"
	"strings"
)

var NewProvider = func(scheme, instance string) (*SampleProvider, error) {
	return New(scheme, instance)
}

type Response guzzle.Response

type SampleProvider struct {
	Address    string
	Scheme     string
	HttpClient *guzzle.Client
	Log        logx.Logger
	Context    context.Context
}

func New(scheme, address string) (*SampleProvider, error) {
	option := &guzzle.Config{
		Address: address,
		Scheme:  scheme,
		//ValidateHost: true,
	}

	httpClient, err := guzzle.NewClient(option)
	if err != nil {
		panic(err)
	}
	var formatter logrusx.JsonFormatter
	log := logrusx.New(logrusx.WithFormatter(formatter))
	return &SampleProvider{Address: address, Scheme: scheme, HttpClient: httpClient, Log: log, Context: context.Background()}, nil

}

func (p *SampleProvider) Do(method, uri string, option IOptionFun) *Response {
	r := p.HttpClient.DoNewRequest(method, uri)
	if len(option.GetHeader()) > 0 {
		for k, y := range option.GetHeader() {
			r.SetParam("Header", k, y)

		}
	}
	if len(option.GetParam()) > 0 {
		for k, y := range option.GetParam() {
			r.SetParam("Query", k, y)
		}
	}
	if len(option.GetBody()) > 0 {
		r.SetBody(strings.NewReader(string(option.GetBody())))
	}

	duration, httpResp, err := p.HttpClient.NewDoRequest(r)
	out := guzzle.RequireOK(duration, httpResp, err)
	if p.Context.Value("debug") != nil {
		o, _ := io.ReadAll(out.RawResponse.Body)
		p.Log.Infof("debug, method: %s, request: %+v, response: %+v\n", method, r.GetRequest(), r, string(o))
		out.RawResponse.Body = io.NopCloser(bytes.NewBuffer(o))
	}
	return (*Response)(out)
}

func (r *Response) Json(useStruct interface{}) error {
	return (*guzzle.Response)(r).Json(useStruct)
}

func (r *Response) Xml(useStruct interface{}) error {
	return (*guzzle.Response)(r).Xml(useStruct)
}
