package client

import (
	"api/config"
	"flag"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	m "k8s.io/metrics/pkg/client/clientset/versioned"
)

func connK3s() (k3s *kubernetes.Clientset, metrics *m.Clientset) {
	var cfg *rest.Config
	var err error

	if os.Getenv("GIN_MODE") == "release" {
		cfg, err = rest.InClusterConfig()
	} else {
		cfg, err = clientcmd.BuildConfigFromFlags("", *flag.String("kubeconfig", config.ENV.K3sPath, "kubeconfig file location"))
	}

	if err != nil {
		panic(err.Error())
	}

	if k3s, err = kubernetes.NewForConfig(cfg); err != nil {
		panic(err.Error())
	}

	if metrics, err = m.NewForConfig(cfg); err != nil {
		panic(err.Error())
	}

	return
}
