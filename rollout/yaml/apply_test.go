package yaml

import (
	"github.com/illuminatingKong/kongming-kit/kubego/client"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
)

func TestDemo(t *testing.T) {
	yamlContent := `apiVersion: v1
kind: Namespace
metadata:
  name: ops
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment-takumi
  namespace: ops
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx
        ports:
        - containerPort: 80`
	d := &Deploy{}
	kubeconfig := []byte(``)
	c, err := client.NewRestClient(kubeconfig)
	if err != nil {
		t.Fatal(err)
	}

	cl, err := client.NewClent(c, runtimeclient.Options{})
	if err != nil {
		t.Fatal(err)
	}
	createErr := d.CreateOrPatch(yamlContent, cl, Option{NameSpace: "ops"})
	if createErr != nil {
		panic(createErr)
	}
}
