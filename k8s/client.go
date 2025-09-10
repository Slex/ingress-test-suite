package k8s

import (
	"ingress-test-suite/logger"
	"ingress-test-suite/pkg/messages"
	"os"
	"path/filepath"

	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var log = logger.GetLogger()

func SetupK8SClient() *kubernetes.Clientset {
	client := getConfig()

	log.Infof(messages.CreteK8SClient, client.AppsV1())

	return client
}

func getConfig() *kubernetes.Clientset {
	var cfg *rest.Config
	var err error

	cfg, err = rest.InClusterConfig()

	if err != nil {
		var kubeConfig string
		if home := homedir.HomeDir(); home != "" {
			kubeConfig = filepath.Join(home, ".kube", "config")
		} else if os.Getenv("KUBECONFIG") != "" {
			kubeConfig = os.Getenv("KUBECONFIG")
		} else {
			log.Errorf(messages.ErrorLoadKubeConfig, err)
			os.Exit(31)
		}

		cfg, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			log.Errorf(messages.ErrorBuildKubeConfig, err)
			os.Exit(31)
		}
	}

	clientSet, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Errorf(messages.FailedCreateKubeClient, err)
		os.Exit(31)
	}

	return clientSet
}
