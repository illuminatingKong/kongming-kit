package example_crd_redis

import (
	"github.com/illuminatingKong/kongming-kit/kubego/client"
	rolloutyaml "github.com/illuminatingKong/kongming-kit/rollout/yaml"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
)

var (
	kubeconfigTest = `apiVersion: v1
kind: Config
clusters:
- name: "sit"
  cluster:
    server: "https://qarancher.hwwt2.com/k8s/clusters/c-88nvp"
users:
- name: "sit"
  user:
    token: "kubeconfig-u-r8fnh.c-88nvp:zmv89q6q8sx5g4jrpf2xfg8fjf5v5dsxzsb66bwrq8sf9b56gr2m4s"


contexts:
- name: "sit"
  context:
    user: "sit"
    cluster: "sit"

current-context: "sit"`
	yamlContentTest = `apiVersion: cache.tongdun.net/v1alpha1
kind: RedisStandby
metadata:
  name: redis-standby-srv-name1
  namespace: redis
spec:
  app: standby-customer
  capacity: 512
  dc: sh
  env: production
  image: harbor.hwwt2.com/ops/op-redis-5.0.9:v2
  monitorimage: harbor.hwwt2.com/ops/redis-exporter:paas1.0
  netmode: ClusterIP
  realname: srv-name1
  secret: ''
  sentinelimage: harbor.hwwt2.com/ops/sentinel-standby:v1
  storageclass: ''
  vip: ''
  labels:             
    plt: customer
    sid: "10001"
    srv: "srv-name1"`
)

func TestCreateRedis(t *testing.T) {
	yamlContent := yamlContentTest
	d := &rolloutyaml.Deploy{}
	kubeconfig := []byte(kubeconfigTest)
	c, err := client.NewRestClient(kubeconfig)
	if err != nil {
		panic(err)
	}

	cl, err := client.NewClient(c, runtimeclient.Options{})
	if err != nil {
		panic(err)
	}

	createErr := d.CreateOrPatchCustomResource(yamlContent, cl)
	if createErr != nil {
		panic(createErr)
	}
}
