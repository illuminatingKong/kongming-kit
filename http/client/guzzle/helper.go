package guzzle

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/hashicorp/go-rootcerts"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

// encodeBody is used to encode a request body
func encodeBody(obj interface{}) (io.Reader, error) {
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(obj); err != nil {
		return nil, err
	}
	return buf, nil
}

// durToMsec converts a duration to a millisecond specified string. If the
// user selected a positive value that rounds to 0 ms, then we will use 1 ms
// so they get a short delay, otherwise Consul will translate the 0 ms into
// a huge default delay.
func durToMsec(dur time.Duration) string {
	ms := dur / time.Millisecond
	if dur > 0 && ms == 0 {
		ms = 1
	}
	return fmt.Sprintf("%dms", ms)
}

func SetupTLSConfig(tlsConfig *TLSConfig) (*tls.Config, error) {
	tlsClientConfig := &tls.Config{
		InsecureSkipVerify: tlsConfig.InsecureSkipVerify,
	}

	if tlsConfig.Address != "" {
		server := tlsConfig.Address
		hasPort := strings.LastIndex(server, ":") > strings.LastIndex(server, "]")
		if hasPort {
			var err error
			server, _, err = net.SplitHostPort(server)
			if err != nil {
				return nil, err
			}
		}
		tlsClientConfig.ServerName = server
	}

	if tlsConfig.CertFile != "" && tlsConfig.KeyFile != "" {
		tlsCert, err := tls.LoadX509KeyPair(tlsConfig.CertFile, tlsConfig.KeyFile)
		if err != nil {
			return nil, err
		}
		tlsClientConfig.Certificates = []tls.Certificate{tlsCert}
	}

	if tlsConfig.CAFile != "" || tlsConfig.CAPath != "" {
		rootConfig := &rootcerts.Config{
			CAFile: tlsConfig.CAFile,
			CAPath: tlsConfig.CAPath,
		}
		if err := rootcerts.ConfigureTLS(tlsClientConfig, rootConfig); err != nil {
			return nil, err
		}
	}

	return tlsClientConfig, nil
}

// decodeBody is used to JSON decode a body
func DecodeBody(resp *http.Response, out interface{}) error {
	dec := json.NewDecoder(resp.Body)
	return dec.Decode(out)
}

func XmlDecodeBody(resp *http.Response, out interface{}) error {
	dec := xml.NewDecoder(resp.Body)
	return dec.Decode(out)
}

func TextDecodeBody(resp *http.Response, out string) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	out = string(body)
	return nil
}
