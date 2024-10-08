package apolloconfig

import (
	"fmt"
	"github.com/illuminatingKong/kongming-kit/http/client/guzzle"
)

// appNavtree 通过appid 查询到 env，cluster
func (p *Provider) appNavtree(appID string) (*RespClusters, error) {
	uri := fmt.Sprintf("%s/apps/%s/navtree", p.EndPoint, appID)
	r := p.HttpClient.DoNewRequest("GET", uri)
	r.SetParam("Header", "Content-Type", "application/json")
	r.SetParam("Header", "Cookie", p.Cookie)
	out := guzzle.RequireOK(p.HttpClient.NewDoRequest(r))
	f := &RespClusters{}
	if err := out.Json(f); err != nil {
		return nil, err
	}
	return f, nil
}
