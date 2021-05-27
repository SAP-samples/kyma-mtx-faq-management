# About

The backend service of the FAQ Management is using SAP's [Cloud Application Programming Model](https://cap.cloud.sap).

It contains these folders and files, following our recommended project layout:

File or Folder | Purpose
---------|----------
`db/` | your domain models and data go here
`srv/` | your service models and code go here
`gen/` | you'll find the generated code here
`package.json` | project metadata and configuration
`k8s/` | holds the manifest of the component
`Dockerfile` | to build the image of this component

## Register `backend` as a Subscription Application

In order to make the `backend` available as a SaaS app, it needs (a) to provide the corresponding endpoint (`server.js`) and (b) be connected with the `saas-registry` service. This callback URI will be called by the [broker](../broker), which is the component that manages the resources needed to create a tenant per consumer subaccount internally. 

Please refer to the official [documentation](https://help.sap.com/viewer/65de2977205c403bbc107264b8eccf4b/Cloud/en-US/3971151ba22e4faa9b245943feecea54.html) for more details on the concept of multi-tenancy on SAP BTP and how to register a multi-tenant application to the SaaS provisioning service.


## Build-in production mode (also for local testing)
- on root directory: `npm install`
- `cds build --production` (this will generate sources in gen/ directory)
- `cd gen`
- `NODE_ENV=production npm start`


## VCAP_SERVICES
The app needs a VCAP_SERVICES env variable or a cds config file.
The VCAP Services need to contain
- SAP HANA Cloud coordinates
- XSUAA coordinates
- Service Manager coordinates

The provided manifest assembles the VCAP_SERVICES variable from the created Service Instances and their respective bindings, plus the secret created from the HANA instance.

## Learn More

Learn more at https://cap.cloud.sap/docs/get-started/.


