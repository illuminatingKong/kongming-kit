package apollotesting

import (
	"fmt"
	"github.com/illuminatingKong/kongming-kit/http/apolloconfig"
	"os"
	"testing"
	"time"
)

func TestRichProcess(t *testing.T) {
	data, err := newTestRich()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}

type TestCaseApolloInfo struct {
	Host        string
	Scheme      string
	Port        string
	UserName    string
	PassWord    string
	Cluster     string
	AppID       string
	ENV         string
	Namespace   string
	TestKeyName string
	TestValue   string
}

func newTestRich() (interface{}, error) {
	var err error
	info := newTestCaseApolloInfo()
	provider, err := apolloconfig.NewProvider(info.UserName, info.PassWord, info.Scheme, info.Host, info.Port)
	if err != nil {
		return nil, err
	}
	data, err := provider.GetApps()
	if err != nil {
		return nil, err
	}

	err = provider.AddPropertyItem(info.AppID, info.ENV, info.Cluster, info.Namespace, info.TestKeyName, info.TestValue)

	releaseTitle := fmt.Sprintf("test-%s", time.Now().String())
	releaseComment := "this is a test"
	err = provider.PublishRelease(info.AppID, info.ENV, info.Cluster, info.Namespace, releaseTitle, releaseComment)
	return data, err
}

func newTestCaseApolloInfo() *TestCaseApolloInfo {
	host := os.Getenv("APOLLO_HOST")
	Scheme := os.Getenv("APOLLO_SCHEME")
	Port := os.Getenv("APOLLO_PORT")
	UserName := os.Getenv("APOLLO_USERNAME")
	PassWord := os.Getenv("APOLLO_PASSWORD")
	AppID := os.Getenv("APOLLO_APPID")
	Env := "PRO"
	Cluster := "default"
	Namespace := os.Getenv("APOLLO_NAMESPACE")
	info := &TestCaseApolloInfo{
		Host:        host,
		Scheme:      Scheme,
		Port:        Port,
		UserName:    UserName,
		PassWord:    PassWord,
		Cluster:     Cluster,
		AppID:       AppID,
		ENV:         Env,
		Namespace:   Namespace,
		TestKeyName: "APOLLO_TESTKEY",
		TestValue:   "APOLLO_TEST",
	}
	return info
}
