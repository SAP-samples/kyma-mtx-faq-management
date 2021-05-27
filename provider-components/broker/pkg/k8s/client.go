package k8s

import (
	"fmt"
	apiRules "github.com/kyma-incubator/api-gateway/api/v1alpha1"
	oauthclientv1alpha1 "github.com/ory/hydra-maester/api/v1alpha1"
	rulev1alpha1 "github.com/ory/oathkeeper-maester/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type KymaClient struct {
	client.Client
}

func CreateK8SClient(kubeConfig string) (*KymaClient, error) {
	cfg, err := rest.InClusterConfig()
	if err == rest.ErrNotInCluster {
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			return nil, fmt.Errorf("error obtaining out of cluster kubeconfig (location %s): %s",
				kubeConfig, err.Error())
		}
	} else if err != nil {
		return nil, fmt.Errorf("error obtaining kubeconfig in cluster: %s", err.Error())
	}
	sch, err := createScheme()
	if err != nil {
		return nil, err
	}

	k8sClient, err := client.New(cfg, client.Options{Scheme: sch})
	if err != nil {
		return nil, err
	}

	return &KymaClient{Client: k8sClient}, nil
}

func createScheme() (*runtime.Scheme, error) {
	sch := scheme.Scheme
	var addToSchemes runtime.SchemeBuilder
	addToSchemes = append(addToSchemes, apiRules.AddToScheme)
	addToSchemes = append(addToSchemes, rulev1alpha1.AddToScheme)
	addToSchemes = append(addToSchemes, oauthclientv1alpha1.AddToScheme)
	addToSchemes = append(addToSchemes, v1.AddToScheme)
	err := addToSchemes.AddToScheme(sch)
	if err != nil {
		return nil, err
	}
	return sch, nil
}
