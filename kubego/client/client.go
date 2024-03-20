package client

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/cluster"
)

func NewRestClient(conf []byte) (*restclient.Config, error) {
	restConf, err := clientcmd.RESTConfigFromKubeConfig(conf)
	if err != nil {
		return nil, errors.Wrap(err, ConnectError.Error())
	}
	return restConf, nil
}

func NewRestToClusterObject(restConfig *restclient.Config) (cluster.Cluster, error) {
	scheme := runtime.NewScheme()
	// add all known types
	// if you want to support custom types, call _ = yourCustomAPIGroup.AddToScheme(scheme)
	_ = clientgoscheme.AddToScheme(scheme)

	c, err := cluster.New(restConfig, func(clusterOptions *cluster.Options) {
		clusterOptions.Scheme = scheme
	})
	if err != nil {
		return nil, errors.Wrap(err, ConnectError.Error())
	}

	return c, nil
}

func NewRestToClientSet(restConfig *restclient.Config) (*kubernetes.Clientset, error) {
	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, errors.Wrap(err, ConnectError.Error())
	}
	return clientSet, nil
}

func NewClent(restConfig *restclient.Config, options runtimeclient.Options) (runtimeclient.Client, error) {
	return runtimeclient.New(restConfig, options)
}
