package apolloconfig

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/illuminatingKong/kongming-kit/http/client/guzzle"
	"sync"
)

var cookieMutex sync.Mutex

type Provider struct {
	HttpClient    *guzzle.Client
	Username      string `validate:"required"`
	Password      string `validate:"required"`
	GatewayHost   string `validate:"required"`
	GatewayPort   string `validate:"required"`
	GatewayScheme string `validate:"required"`
	Cookie        string
	EndPoint      string
}

func NewProvider(username, password, gatewayScheme, gatewayHost, gatewayPort string) (IApolloConfig, error) {
	var err error
	option := &guzzle.Config{
		Address:      fmt.Sprintf("%s:%s", gatewayHost, gatewayPort),
		Scheme:       gatewayScheme,
		ValidateHost: true,
	}
	httpClient, err := guzzle.NewClient(option)
	if err != nil {
		return nil, err
	}
	ed := fmt.Sprintf("%s://%s:%s", gatewayScheme, gatewayHost, gatewayPort)

	p := &Provider{
		HttpClient:    httpClient,
		Username:      username,
		Password:      password,
		GatewayHost:   gatewayHost,
		GatewayPort:   gatewayPort,
		GatewayScheme: gatewayScheme,
		EndPoint:      ed,
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(p)
	err = p.newWebCookie()
	return p, err
}

func (p *Provider) newWebCookie() error {
	var err error
	cookies, err := p.sessionID()
	if err != nil {
		return err
	}
	p.Cookie = toCookieString(cookies)
	return err
}
