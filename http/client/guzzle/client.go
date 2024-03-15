package guzzle

import (
	"fmt"
	"net/http"
	"os"

	"strconv"
	"strings"
)

func DefaultConfig() *Config {
	return defaultConfig(DefaultPooledTransport)
}

func defaultConfig(transportFn func() *http.Transport) *Config {
	config := &Config{
		Address:   "127.0.0.1:8080",
		Scheme:    "http",
		Transport: transportFn(),
	}

	if addr := os.Getenv(HTTPAddrEnvName); addr != "" {
		config.Address = addr
	}

	if auth := os.Getenv(HTTPAuthEnvName); auth != "" {
		var username, password string
		if strings.Contains(auth, ":") {
			split := strings.SplitN(auth, ":", 2)
			username = split[0]
			password = split[1]
		} else {
			username = auth
		}

		config.HttpAuth = &HttpBasicAuth{
			Username: username,
			Password: password,
		}
	}

	if ssl := os.Getenv(HTTPSSLEnvName); ssl != "" {
		enabled, err := strconv.ParseBool(ssl)
		if err != nil {
			panic(fmt.Sprintf("[WARN] client: could not parse %s: %s", HTTPSSLEnvName, err))
		}

		if enabled {
			config.Scheme = "https"
		}
	}

	return config
}

// NewHttpClient returns an http client configured with the given Transport and TLS
// config.
func NewHttpClient(transport *http.Transport, tlsConf TLSConfig) (*http.Client, error) {
	client := &http.Client{
		Transport: transport,
	}

	// TODO (slackpad) - Once we get some run time on the HTTP/2 support we
	// should turn it on by default if TLS is enabled. We would basically
	// just need to call http2.ConfigureTransport(transport) here. We also
	// don't want to introduce another external dependency on
	// golang.org/x/net/http2 at this time. For a complete recipe for how
	// to enable HTTP/2 support on a transport suitable for the API client
	// library see agent/http_test.go:TestHTTPServer_H2.

	if transport.TLSClientConfig == nil {
		tlsClientConfig, err := SetupTLSConfig(&tlsConf)

		if err != nil {
			return nil, err
		}

		transport.TLSClientConfig = tlsClientConfig
	}

	return client, nil
}

func (r *request) AddBulkHeader(value map[string]string) {
	for k, v := range value {
		r.header.Set(k, v)
	}
}

// setQueryOptions is used to annotate the request with
// additional query options
func (r *request) SetParam(position, name, value string) {
	if len(value) > 0 {
		switch position {
		case "Header":
			r.header.Set(name, value)
		case "Query":
			r.params.Set(name, value)
		}
	}
}

func (r *request) SetBody(body *strings.Reader) {
	r.body = body
}
