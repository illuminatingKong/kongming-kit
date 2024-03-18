package samplehttp

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

const (
	GetMethod    = "GET"
	PostMethod   = "POST"
	PutMethod    = "PUT"
	PatchMethod  = "PATCH"
	DeleteMethod = "DELETE"
)

type IOptionFun interface {
	AddHeader(map[string]string) map[string]string
	AddParam(map[string]string) map[string]string
	AddBody(body interface{}) error
	GetHeader() map[string]string
	GetParam() map[string]string
	GetBody() []byte
}

type Option struct {
	Header map[string]string
	Param  map[string]string
	Body   *strings.Reader
}

func (o *Option) GetHeader() map[string]string {
	return o.Header
}

func (o *Option) GetParam() map[string]string {
	return o.Param
}

func (o *Option) GetBody() []byte {
	var data []byte
	var err error
	defer func() {
		if err := recover(); err != nil {
			data = []byte{}
		}
	}()
	data, err = io.ReadAll(o.Body)
	if err != nil {
		return []byte{}
	}
	return data
}

func (o *Option) AddHeader(header map[string]string) map[string]string {
	o.Header = header
	return o.Header
}

func (o *Option) AddParam(param map[string]string) map[string]string {
	o.Param = param
	return o.Param
}

func (o *Option) AddBody(body interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	o.Body = strings.NewReader(string(b))
	return nil
}

func (p *SampleProvider) SetBaseAuth(username, password string) *SampleProvider {
	p.HttpClient.SetBaseAuth(username, password)
	return p
}

func (p *SampleProvider) SetValidateHost(v bool) *SampleProvider {
	p.HttpClient.SetValidateHost(v)
	return p
}

func (p *SampleProvider) SetTransport(transportFn func() *http.Transport) *SampleProvider {
	p.HttpClient.SetTransport(transportFn)
	return p
}

func (p *SampleProvider) SetDebug(v bool) *SampleProvider {
	p.Context = context.WithValue(p.Context, "debug", v)
	return p
}
