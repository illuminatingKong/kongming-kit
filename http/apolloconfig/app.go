package apolloconfig

import (
	"fmt"
	"github.com/illuminatingKong/kongming-kit/http/client/guzzle"
)

// Apps 获取所有的 apollo appID
func (p *Provider) apps() (*[]RespApp, error) {
	uri := fmt.Sprintf("%s/apps/", p.EndPoint)
	r := p.HttpClient.DoNewRequest("GET", uri)
	r.SetParam("Header", "Content-Type", "application/json")
	r.SetParam("Header", "Cookie", p.Cookie)
	out := guzzle.RequireOK(p.HttpClient.NewDoRequest(r))
	f := &[]RespApp{}
	if err := out.Json(f); err != nil {
		return nil, err
	}
	return f, nil
}
