package client

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/rest"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"testing"
)

func TestNewClusterObject(t *testing.T) {
	var cfg *rest.Config
	c, err := NewRestToClusterObject(cfg)
	Describe("NewClusterObject", func() {
		It("should return an error if there is no Config", func() {
			Expect(c).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("must specify Config"))
		})
	})

}

func TestNewRestClient(t *testing.T) {
	config := createValidTestConfig()
	b, err := transformTestConfigToBytes(config)
	if err != nil {
		panic(err)
	}
	_, err = NewRestClient(b)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

}

func transformTestConfigToBytes(config *clientcmdapi.Config) ([]byte, error) {
	return json.Marshal(config)
}

func createValidTestConfig() *clientcmdapi.Config {
	const (
		server = "https://anything.com:8080"
		token  = "the-token"
	)

	config := clientcmdapi.NewConfig()
	config.Clusters["clean"] = &clientcmdapi.Cluster{
		Server: server,
	}
	config.AuthInfos["clean"] = &clientcmdapi.AuthInfo{
		Token: token,
	}
	config.Contexts["clean"] = &clientcmdapi.Context{
		Cluster:  "clean",
		AuthInfo: "clean",
	}
	config.CurrentContext = "clean"

	return config
}
