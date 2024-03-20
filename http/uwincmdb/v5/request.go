package v5

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type CMDBHttpRequester struct {
	Uri          string
	Method       string
	AK           string
	SK           string
	Param        map[string]string
	Body         *strings.Reader
	bodyMd5      string
	Host         string
	StringToSign string
	Expires      string
	Err          error
	HeaderHost   string
}

// NewCMDBHttpRequesterClient ip or domain
func NewCMDBHttpRequesterClient(Host string) *CMDBHttpRequester {
	// TODO ValidateHost is True
	r := &CMDBHttpRequester{Host: Host}
	return r
}

func (request *CMDBHttpRequester) WithAK(AK string) *CMDBHttpRequester {
	request.AK = AK
	return request
}

func (request *CMDBHttpRequester) WithSK(SK string) *CMDBHttpRequester {
	request.SK = SK
	return request
}

func (request *CMDBHttpRequester) WithUri(Uri string) *CMDBHttpRequester {
	request.Uri = Uri
	return request
}

func (request *CMDBHttpRequester) WithBody(body interface{}) *CMDBHttpRequester {
	b, err := json.Marshal(body)
	if err != nil {
		request.Err = err
	}
	play := strings.NewReader(string(b))
	m := md5.New()
	m.Write(b)
	request.bodyMd5 = fmt.Sprintf("%x", m.Sum(nil))
	request.Body = play
	return request
}

func (request *CMDBHttpRequester) WithBodyByte(body []byte) *CMDBHttpRequester {
	play := strings.NewReader(string(body))
	m := md5.New()
	m.Write(body)
	request.bodyMd5 = fmt.Sprintf("%x", m.Sum(nil))
	request.Body = play
	return request
}

func (request *CMDBHttpRequester) WantPost() {
	request.Method = "POST"
	request.Sign(time.Now())

}

func (request *CMDBHttpRequester) genSignature(requestime time.Time) *CMDBHttpRequester {
	var urlParams = ""
	var contentType = "application/json"

	str_sign := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s", request.Method, request.Uri, urlParams, contentType, request.bodyMd5, fmt.Sprintf("%d", requestime.Unix()), request.AK)
	sigHash := hmac.New(sha1.New, []byte(request.SK))
	sigHash.Write([]byte(str_sign))
	StringToSign := hex.EncodeToString(sigHash.Sum(nil))
	request.StringToSign = StringToSign
	return request
}

func (request *CMDBHttpRequester) Sign(expires time.Time) {
	t := fmt.Sprintf("%d", expires.Unix())
	request.Expires = string(t)
	request.genSignature(expires)

}

func (request *CMDBHttpRequester) WithHeaderHost(host string) *CMDBHttpRequester {
	request.HeaderHost = host
	return request
}
