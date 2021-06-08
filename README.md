[![REUSE status](https://api.reuse.software/badge/github.com/SAP-samples/kyma-mtx-faq-management)](https://api.reuse.software/info/github.com/SAP-samples/kyma-mtx-faq-management)


# Multi-tenant FAQ Management Application on SAP BTP, Kyma Runtime

[![Video here](https://img.youtube.com/vi/hnD7Lr_2464/0.jpg)](https://www.youtube.com/watch?v=hnD7Lr_2464)

## About

Kubernetes is the baseline of cloud-native infrastructures. The code (and documentation in this repo) are going to take you on a fast-forward journey towards a business service live and kicking on Kubernetes ([SAP BTP, Kyma Runtime](https://discovery-center.cloud.sap/#/serviceCatalog/kyma-runtime?region=all)). This sample project shows the integration points into the SaaS Application Registry to be consumed from within the BTP Cockpit.

This scenario does not cover the steps required to commercialize the service. However, it provides the necessary preconditions to start the process within SAP.

## Scenario Description

The FAQ management service is built on SAP BTP using the [Cloud Application Programming Model (CAP)](https://cap.cloud.sap/). It allows customers to manage a list of Frequently Asked Questions through an SAP Fiori interface. Overall, this is a simple and yet powerful multi-tenant service that allows each subscriber to configure roles and add users who can access the isolated data of the respective tenant. 


## Requirements

Make sure you have the following tools installed and services provisioned:

> You can complete all steps in production as well as in trial accounts.

* Create an [SAP BTP Trial account](https://developers.sap.com/tutorials/hcp-create-trial-account.html) in the region EU10
* Deploy [SAP HANA Cloud](https://developers.sap.com/tutorials/hana-cloud-deploying.html)
* Enable [SAP BTP, Kyma Runtime](https://developers.sap.com/tutorials/cp-kyma-getting-started.html)​
* Install the [Kubernetes Command Line Interface (CLI)](https://developers.sap.com/tutorials/cp-kyma-download-cli.html)
* Install the [Cloud Foundry CLI](https://developers.sap.com/tutorials/cp-cf-download-cli.html)

We will use a few SAP BTP services during the deployment. Please make sure you see the following services in the SAP BTP Cockpit service marketplace:

* SAP HANA Schemas & HDI Containers (`hdi-shared` plan)
* SaaS Provisioning (`container` and `application` plans)
* XSUAA: Authorization & Trust Management (`broker` plan)


## Optional

Install the following tools if you would like to run the individual components on your local machine and build your own images:
* Install [Docker](https://docs.docker.com/get-docker/) and log in with your [Docker ID](https://docs.docker.com/docker-id/) (or use another Docker registry of your choosing)
* Install [Golang](http://golang.org/) 
* Install [Node.js](https://nodejs.org/en/download/) 
* Install [@sap/cds-dk](https://cap.cloud.sap/docs/get-started/) 

## Download and Installation

1. [Set up the provider subaccount, which will offer the service to other subaccounts.](provider-components/readme.md)
2. [Subscribe to the SaaS solution in the consumer subaccount.](consumer/readme.md)

## Limitations

This sample project has been created to sketch specific ideas to demonstrate how to build multi-tenant apps on SAP BTP, Kyma runtime.
Please be aware that this project still contains a few workarounds and takes "technical shortcuts" like:
* SAP HANA Cloud is currently not natively available in the Kyma runtime. Here we use the workaround and provision the instance in the Cloud Foundry environment.
* As of now, this project doesn't enforce access restrictions on the admin endpoint of the services. To keep things simple, we hide some UI elements when the user doesn't have sufficient scope. The only data separation in this sample happens on a tenant level.

## Support

Please use the GitHub bug tracking system to post questions and bug reports.


## Knows Issues
None so far


## Contributing
This project is only updated by SAP employees and only accepts bug reports but no other contribution. 

## License

Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved. This project is licensed under the Apache Software License, version 2.0, except as noted otherwise in the [LICENSE](LICENSES/Apache-2.0.txt) file.