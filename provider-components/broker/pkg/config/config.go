package config

import (
	"flag"
	"fmt"
)

type Config struct {
	Namespace           string
	KubeConfigLocation  string
	TargetServiceName   string
	TargetServicePort   uint32
	IngressGateway      string
	SaaSProvisionerPort uint
	KymaDomain          string
	SaaSAppServiceName  string
	SaaSAppServicePort  uint32
	CapProvisioningURL  string
}

func ParseConfig() (*Config, error) {
	config := &Config{}
	flag.StringVar(&config.Namespace, "resourceNamespace", "default", "Namespace in which brokered "+
		"resources will be provisioned in.")
	flag.StringVar(&config.KubeConfigLocation, "kubeConfig", "", "Location of kubeconfig file if "+
		"broker is run out of cluster")
	flag.StringVar(&config.TargetServiceName, "targetServiceName", "", "Name of the Service that "+
		"will be exposed via API Rule")
	targetPort := flag.Uint("targetServicePort", 8080, "Port of the Service that will be exposed via "+
		"API Rule")
	flag.StringVar(&config.IngressGateway, "ingressGateway", "kyma-gateway.kyma-system.svc.cluster.local",
		"Kyma gateway used for the ingress")
	flag.UintVar(&config.SaaSProvisionerPort, "saaSProvisionerPort", 8081, "Port to expose "+
		"SaaS Provisioner Callback")
	flag.StringVar(&config.KymaDomain, "kymaDomain", "",
		"Domain of the kyma cluster")
	flag.StringVar(&config.SaaSAppServiceName, "saasAppServiceName", "", "Name of the SaaS App Service that "+
		"will be exposed via API Rule")
	saasPort := flag.Uint("saasAppServicePort", 8080, "Port of the SaaS App Service that will be exposed via "+
		"API Rule")
	flag.StringVar(&config.CapProvisioningURL, "capProvisioningURL", "", "CAP Service Provisioning URl")
	flag.Parse()

	if config.TargetServiceName == "" {
		return nil, fmt.Errorf("mandatory value for targetServiceName is missing")
	}

	if config.SaaSAppServiceName == "" {
		return nil, fmt.Errorf("mandatory value for saasAppServiceName is missing")
	}

	if config.KymaDomain == "" {
		return nil, fmt.Errorf("mandatory value for kymaDomain is missing")
	}

	config.TargetServicePort = uint32(*targetPort)
	config.SaaSAppServicePort = uint32(*saasPort)

	if config.CapProvisioningURL == "" {
		config.CapProvisioningURL = fmt.Sprintf("http://%s:%d", config.TargetServiceName, config.TargetServicePort)
	}
	return config, nil
}

func (c *Config) Print() {
	fmt.Println("Broker configured with: ")
	fmt.Printf("Namespace in which brokered resources will be provisioned in: %s\n", c.Namespace)
	fmt.Printf("Location of kubeconfig file if broker is run out of cluster: %s\n", c.KubeConfigLocation)
	fmt.Printf("Name of the Service that will be exposed via API Rule: %s\n", c.TargetServiceName)
	fmt.Printf("Port of the Service that will be exposed via API Rule: %d\n", c.TargetServicePort)
	fmt.Printf("Kyma gateway used for the ingress: %s\n", c.IngressGateway)
	fmt.Printf("SaaS Provisioner Port: %d\n", c.SaaSProvisionerPort)
	fmt.Printf("Domain of the kyma cluster: %s\n", c.KymaDomain)
	fmt.Printf("Name of the SaaS App Service that will be exposed via API Rule: %s\n", c.SaaSAppServiceName)
	fmt.Printf("Port of the SaaS App Service that will be exposed via API Rule: %d\n", c.SaaSProvisionerPort)
	fmt.Printf("Cap Provisioning URl %s\n", c.CapProvisioningURL)

}
