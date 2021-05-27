package main

import (
	"context"
	"fmt"
	cliConfig "broker/pkg/config"
	"broker/pkg/instancemgr"
	"broker/pkg/k8s"
	"broker/pkg/kyma"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"code.cloudfoundry.org/lager"
)

func main() {
	// Initial Lager instance as the default logger
	logger := lager.NewLogger("faq-broker")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	// Parse Configuration provided as input flags
	config, err := cliConfig.ParseConfig()
	if err != nil {
		logger.Fatal(cliConfig.ActionParseConfig, err)
	}
	config.Print()

	// Create Kubernetes client-go instance to interact with the kube api server
	kymaClient, err := k8s.CreateK8SClient(config.KubeConfigLocation)
	if err != nil {
		logger.Fatal(cliConfig.ActionK8sClientCreation, err, lager.Data{
			cliConfig.KeyMessage: "error creating k8s client",
		})
	}
	// Initialse k8s client delegate to create k8s resources
	k8sService := instancemgr.New(logger, kymaClient, config.Namespace, config.TargetServiceName,
		config.TargetServicePort, config.IngressGateway, config.KymaDomain)


	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	saasProvisionerHandler := kyma.NewKymaHandler(logger, k8sService, config.Namespace, config.KymaDomain,
		config.SaaSAppServiceName, config.SaaSAppServicePort, config.CapProvisioningURL)

	// Create Saas Provisioner handler
	kymaSrv := http.Server{
		Addr:    fmt.Sprintf(":%d", config.SaaSProvisionerPort),
		Handler: saasProvisionerHandler,
	}

	//Start the Saas Provisioner Service
	go func() { // gofunc due to separate http server, as the broker server has middleware enabled
		//kymaAccessHandler := kyma.New(logger, k8sService)

		logger.Info(cliConfig.ActionStartSaasProvisioner, lager.Data{
			cliConfig.KeyMessage: fmt.Sprintf("Ready for APIRule creation on port %d", config.SaaSProvisionerPort),
		})

		if err := kymaSrv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				logger.Fatal(cliConfig.ActionStartSaasProvisioner,
					err, lager.Data{
						cliConfig.KeyMessage: fmt.Sprintf("Error Serving on Port %d", config.SaaSProvisionerPort),
					})
			}
		}
	}()

	terminate := <-interrupt
	switch terminate {
	case os.Interrupt:
		logger.Info(cliConfig.ActionTerminate, lager.Data{
			cliConfig.KeyMessage: "Received Interruption Signal, Terminating...",
		})
	case syscall.SIGTERM:
		logger.Info(cliConfig.ActionTerminate, lager.Data{
			cliConfig.KeyMessage: "Received SIGTERM Signal, Terminating...",
		})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shut down Service Broker
	logger.Info(cliConfig.ActionTerminate, lager.Data{
		cliConfig.KeyMessage: "Shutting down broker handler",
	})
	// Shut down Saas Provisioner
	logger.Info(cliConfig.ActionTerminate, lager.Data{
		cliConfig.KeyMessage: "Shutting down SaaS provisioner handler",
	})
	err = kymaSrv.Shutdown(ctx)
	if err != nil {
		logger.Error(cliConfig.ActionTerminate, err, lager.Data{
			cliConfig.KeyMessage: "Error in terminating SaaS provisioning server",
		})
	}
	logger.Info(cliConfig.ActionTerminate, lager.Data{
		cliConfig.KeyStatus: cliConfig.StatusOk,
	})

}
