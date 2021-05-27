# About

The UI component consists of an extended [approuter](https://www.npmjs.com/package/@sap/approuter) that exposes the user interface based on micro frontend framework [Luigi](https://luigi-project.io) (an open-source project maintained by SAP).

File or Folder | Purpose
---------|----------
`approuter/` | The extended [approuter](https://www.npmjs.com/package/@sap/approuter). This folder contains `xs-app.json` to define the destinations and `approuter-start.js` to bootstrap approuter
`approuter/public` | Public-facing JavaScript and HTML content to be executed in the browser
`k8s/` | holds the manifest of the component
`Dockerfile` | to build the image of this component

## Architecture
The approuter bootstrapper `approuter-start.js` also enriches some handy user info (the `/me` endpoint), including the `assignedScopes` array.

The UI is based on Luigi, and it is sitting behind the approuter. The `xs-app.json` defines that all requests shall be routed either to a local directory or to the destinations.
The navigation structure is defined in `public/luigi-config.js` and will be assembled on the fly by intercepting the `/me` lifecycle event.

Most panels of the frontend are built with [Vue.js](https://vuejs.org/) and [Fundamental Vue](https://sap.github.io/fundamental-vue/#/). We also added a traditional [SAP Fiori elements](https://experience.sap.com/fiori-design-web/smart-templates/) interface to demonstrate that Luigi can manage [micro frontends](https://martinfowler.com/articles/micro-frontends.html).


## Register `user-interface` as a Subscription Application

In order to make the `user-interface` available as a consumer app, it needs to be registered with the `saas-registry`
service. This is achieved by creating a service instance of the `saas-registry` service. This service instance requires a callback URI which is invoked when a consumer tries to subscribe from the consumer subaccount. This callback URI is implemented within the [broker](../broker), which is the component that manages the resources needed to create a tenant per consumer subaccount internally. 

For more details on the concept of multi-tenancy on SAP BTP, please refer to the official [documentation](https://help.sap.com/viewer/65de2977205c403bbc107264b8eccf4b/Cloud/en-US/3971151ba22e4faa9b245943feecea54.html).

In a nutshell, to register a multi-tenant application to the SaaS provisioning service, one has to create an instance of the `saas-registry` service. This provisioning requires parameters defining callback to a SaaS subscription handler which in our case is the `broker`.

## API Rule Requirements for approuter in Kyma
To run an approuter successfully in Kyma, API Rules and implicit ingress gateway need to be configured to transport the [`X-FORWARDED-HOST` header](https://blogs.sap.com/2020/10/21/cloud-native-lab-2-comparing-cloud-foundry-and-kyma-manifests/) to the approuter instance. This is done using mutator instructions in the API Rule Resource. 
This sample project uses the [broker](../broker/) to set up the API Rule and fill in the tenant subdomain at the right places. 

Example:
```
apiVersion: gateway.kyma-project.io/v1alpha1
kind: APIRule
metadata:
  labels:
    SaaSSubdomain: faqconsumer
    subscribedTenantId: 1aef26f1-9b80-4576-807e-348a9aa55093
  name: faqconsumer
 spec:
  gateway: kyma-gateway.kyma-system.svc.cluster.local
  rules:
  - accessStrategies:
    - handler: noop
    methods:
    - GET
    - POST
    - PUT
    - PATCH
    - DELETE
    - HEAD
    mutators:
    - config:
        headers:
          x-forwarded-host: faqconsumer.aaa57ed.kyma.shoot.live.k8s-hana.ondemand.com
      handler: header
    path: /.*
  service:
    host: faqconsumer
    name: kyma-faq-ui-shell
    port: 5000
```
