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

//type Response guzzle.Response

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
	}

	httpClient, err := guzzle.NewClient(option)
	if err != nil {
		panic(err)
	}
	var formatter logrusx.JsonFormatter
	log := logrusx.New(logrusx.WithFormatter(formatter))
	return &SampleProvider{Address: address, Scheme: scheme, HttpClient: httpClient, Log: log,
		Context: context.TODO()}, nil

}

func (p *SampleProvider) Do(method, uri string, option IOptionFun) *guzzle.Response {
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

	return out
}

func (p *SampleProvider) Byte(out *guzzle.Response) ([]byte, error) {
	o, err := io.ReadAll(out.RawResponse.Body)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (p *SampleProvider) Json(out *guzzle.Response, useStruct interface{}) (interface{}, error) {
	err := out.Json(useStruct)
	if err != nil {
		return nil, err
	}
	return useStruct, nil
}
func (p *SampleProvider) Xml(out *guzzle.Response, useStruct interface{}) (interface{}, error) {
	err := out.Xml(useStruct)
	if err != nil {
		return nil, err
	}
	return useStruct, nil
}
