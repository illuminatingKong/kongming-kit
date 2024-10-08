package apolloconfig

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/illuminatingKong/kongming-kit/http/client/guzzle"
	"strings"
	"time"
)

// Namespaces get the Namespace configuration, need to provide the appID, env, and namespace
// http://localhost:8070/apps/SampleApp/envs/DEV/clusters/default/namespaces
func (p *Provider) namespaces(appID, env, cluster string) (*[]RespNamespaces, error) {
	uri := fmt.Sprintf("%s/apps/%s/envs/%s/clusters/%s/namespaces", p.EndPoint, appID, env, cluster)
	r := p.HttpClient.DoNewRequest("GET", uri)
	r.SetParam("Header", "Cookie", p.Cookie)
	r.SetParam("Header", "Content-Type", "application/json;charset=UTF-8")
	out := guzzle.RequireOK(p.HttpClient.NewDoRequest(r))
	f := &[]RespNamespaces{}
	if err := out.Json(f); err != nil {
		return nil, err
	}
	return f, nil
}

// publishRelease 发布配置
func (p *Provider) publishRelease(appID, env, cluster, namespace, releaseTitle, releaseComment string) error {
	uri := fmt.Sprintf("%s/apps/%s/envs/%s/clusters/%s/namespaces/%s/releases", p.EndPoint, appID, env, cluster, namespace)
	r := p.HttpClient.DoNewRequest("POST", uri)
	r.SetParam("Header", "Cookie", p.Cookie)
	r.SetParam("Header", "Content-Type", "application/json;charset=UTF-8")
	if len(releaseComment) == 0 {
		currentTime := time.Now()
		formattedTime := currentTime.Format("2006-01-02 15:04:05.000")
		releaseTitle = formattedTime
	}
	body := ReqReleaseBody{
		ReleaseTitle:       releaseTitle,
		ReleaseComment:     releaseComment,
		IsEmergencyPublish: false,
	}
	bj, err := json.Marshal(body)
	if err != nil {
		return err
	}
	play := strings.NewReader(string(bj))
	r.SetBody(play)
	out := guzzle.RequireOK(p.HttpClient.NewDoRequest(r))
	f := &RespPublishedRelease{}
	if err := out.Json(f); err != nil {
		em := fmt.Sprintf("request: %+v error: %+v", r.GetRequest(), err.Error())
		return errors.New(em)
	}
	return nil
}

func (p *Provider) getReverseHistoryReleases(appID, env, cluster, namespace, page, size string) (*[]RespHistoryRelease, error) {
	uri := fmt.Sprintf("%s/apps/%s/envs/%s/clusters/%s/namespaces/%s/releases/histories", p.EndPoint, appID, env, cluster, namespace)
	r := p.HttpClient.DoNewRequest("GET", uri)
	r.SetParam("Query", "page", page)
	r.SetParam("Query", "size", size)
	r.SetParam("Header", "Cookie", p.Cookie)
	r.SetParam("Header", "Content-Type", "application/json;charset=UTF-8")
	out := guzzle.RequireOK(p.HttpClient.NewDoRequest(r))
	f := &[]RespHistoryRelease{}
	if err := out.Json(f); err != nil {
		em := fmt.Sprintf("request: %+v error: %+v", r.GetRequest(), err.Error())
		return nil, errors.New(em)
	}
	return f, nil
}

func (p *Provider) addPropertyItem(appID, env, cluster, namespace, key, value string) error {
	var err error
	body := &ReqAddPropertyItem{
		TableViewOperType:  CreateItem,
		Key:                key,
		Value:              value,
		AddItemBtnDisabled: true,
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(body)
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s/apps/%s/envs/%s/clusters/%s/namespaces/%s/item", p.EndPoint, appID, env, cluster, namespace)
	r := p.HttpClient.DoNewRequest("POST", uri)
	r.SetParam("Header", "Cookie", p.Cookie)
	r.SetParam("Header", "Content-Type", "application/json;charset=UTF-8")
	bj, err := json.Marshal(body)
	if err != nil {
		return err
	}
	play := strings.NewReader(string(bj))
	r.SetBody(play)
	out := guzzle.RequireOK(p.HttpClient.NewDoRequest(r))
	f := &RespItemChangeBody{}
	if err := out.Json(f); err != nil {
		em := fmt.Sprintf("request: %+v error: %+v", r.GetRequest(), err.Error())
		return errors.New(em)
	}
	if f.Id == 0 {
		return errors.New("add property item failed")
	}
	return nil
}
