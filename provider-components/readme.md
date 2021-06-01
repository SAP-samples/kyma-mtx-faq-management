# Provider Components

The provider subaccount will run all components that need to be exposed by the **provider** of the SaaS solution (hence the name "provider subaccount"). Once deployed, these components will listen for subscription events from other subaccounts (so-called "consumer subaccounts") and perform "subscribe" and "unsubscribe" tasks when called.

As of today, not all required services are available in SAP BTP, Kyma environment. For this reason, we provision a few services in the Cloud Foundry environment and make them accessible to the Kyma environment via service keys.

> This guide only summarizes the mandatory actions you need to complete to get up and running. For details about the individual components, please refer to the readme files of each subdirectory.

## Installation

0. Connect the command-line clients to SAP BTP
    - Download the [`kubeconfig`](https://developers.sap.com/tutorials/cp-kyma-download-cli.html#2d284324-bdd2-4f4b-b786-bab367947689) to connect to your Kyma cluster
    - Run [`cf login`](https://developers.sap.com/tutorials/cp-cf-download-cli.html#1ca87eac-c53f-4ced-a059-c304b1b34cd4) to connect to your Cloud Foundry organization


1. Create a namespace called `project-faq` from the Kyma console. This namespace will hold all resources. 

    ```
    kubectl create namespace project-faq
    ```

2. As a next step, you need to create a Config Map that holds the "Cluster Domain". This domain is referenced by many deployments later on.

    ```
    kubectl create configmap -n project-faq cluster-domain --from-literal cluster-domain="<Your Cluster ID>.kyma.shoot.live.k8s-hana.ondemand.com"
    ```
    > You can find your cluster-ID in the URL of the Kyma console.

2. Add the cluster domain next to the `TODO(For the User)` comments in the components manifest for later:
    - provider-content/user-interface/k8s/manifest.yaml

3. As of the time of creating this documentation SAP HANA Cloud instances can only be provisioned from within a Cloud Foundry context. Hence there is a bit of manual work required. If you haven't done so yet,  [enable Cloud Foundry](https://developers.sap.com/tutorials/hcp-create-trial-account.html) in your account and create a space. Within the context of that space, [create a new SAP HANA Cloud Database](https://developers.sap.com/tutorials/hana-cloud-deploying.html) instance. **And make sure SAP HANA Cloud is started and allows traffic from all IP addresses.**
4. Although Service Manager is available as a brokered service in Kyma, it cannot be used in this context. Instead, the instance of `service-manager` (plan `container`) must be created in the Cloud Foundry space, where the SAP HANA Cloud instance lives. Then, the service key needs to be created and exposed as Kubernetes secret. The backend component will use this secret later on to provision HDI containers on demand.
    - Use the CF CLI to create this instance and download a service key.
    
    Unix:
    ```
    cf create-service service-manager container faq-saas-container
    cf create-service-key faq-saas-container faq-container-key
    cf service-key faq-saas-container faq-container-key | tail -n +3 > faq-container-key.json
    ```
        
    Windows:
    ```
    cf create-service service-manager container faq-saas-container
    cf create-service-key faq-saas-container faq-container-key
    cf service-key faq-saas-container faq-container-key > faq-container-key.json
    # Windows users must remove the first 3 lines of this file manually before proceeding
    ```
    - Create a Kubernetes secret to make this service key available to Kyma
    ```
    kubectl create secret generic -n project-faq sm-credentials --from-file=credentials=faq-container-key.json
    ```
3. Deploy the [backend](backend) component
    ```shell
    kubectl apply -f provider-components/backend/k8s/manifest.yaml 
    ```
3. Deploy the [exporter](exporter) component
    ```shell
    kubectl apply -f provider-components/exporter/k8s/manifest.yaml 
    ```
3. Deploy the [user-interface](user-interface) component to register the application as a Subscription Application
    ```shell
    kubectl apply -f provider-components/user-interface/k8s/manifest.yaml 
    ```
3. Deploy the [service broker](broker) with the following command:
    ```shell
    kubectl apply -f provider-components/broker/k8s/manifest.yaml 
    ```
4. Use the Kyma Console or this command to check when all pods are up
    ```shell
    kubectl -n project-faq get pods
    ```

## Hosted images

Since Kubernetes is based on containers and containers are built from container images, almost all the code in this repo somehow needs to end up in one. 

To simplify your life, we build a set of default images from this repo and store them in our [GitHub Container Registry](https://) for re-use.

If you want to build your own images, that can also be done. Each component that requires a docker image also comes with its own Dockerfile that only requires a `docker build`.

## Deinstallation

Run the following command to remove all Kyma artifacts. Please be aware that you need to restore your namespace when you run these commands.

```
kubectl delete -f provider-components/exporter/k8s/manifest.yaml 
kubectl delete -f provider-components/backend/k8s/manifest.yaml 
kubectl delete -f provider-components/user-interface/k8s/manifest.yaml 
kubectl delete -f provider-components/broker/k8s/manifest.yaml 
```

