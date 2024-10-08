package apolloconfig

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const (
	CreateItem = "create"
)

type IApolloConfig interface {
	SessionID() ([]*http.Cookie, error)
	GetApps() (*[]RespApp, error)
	GetCluster(appid string) (*RespClusters, error)
	Namespaces(appID, env, cluster string) (*[]RespNamespaces, error)
	AddPropertyItem(appID, env, cluster, namespace, key, value string) error
	PublishRelease(appID, env, cluster, namespace, releaseTitle, releaseComment string) error
}

func (p *Provider) SessionID() ([]*http.Cookie, error) {
	cookieMutex.Lock()
	defer cookieMutex.Unlock()
	var err error
	cookieInstance, err := p.sessionID()
	return cookieInstance, err
}

func (p *Provider) sessionID() ([]*http.Cookie, error) {
	login := fmt.Sprintf("%s/signin", p.EndPoint)
	j, _ := cookiejar.New(nil)
	client := &http.Client{Jar: j}
	defer client.CloseIdleConnections()
	_, err := client.PostForm(login, url.Values{
		"username": {p.Username},
		"password": {p.Password},
	})

	urlX, _ := url.Parse(p.EndPoint)
	return client.Jar.Cookies(urlX), err
}

func (p *Provider) GetApps() (*[]RespApp, error) {
	return p.apps()
}

func (p *Provider) GetCluster(appid string) (*RespClusters, error) {
	return p.appNavtree(appid)
}

func (p *Provider) Namespaces(appID, env, namespace string) (*[]RespNamespaces, error) {
	return p.namespaces(appID, env, namespace)
}

func (p *Provider) AddPropertyItem(appID, env, cluster, namespace string, key, value string) error {
	return p.addPropertyItem(appID, env, cluster, namespace, key, value)
}

func (p *Provider) PublishRelease(appID, env, cluster, namespace, releaseTitle, releaseComment string) error {
	return p.publishRelease(appID, env, cluster, namespace, releaseTitle, releaseComment)
}
