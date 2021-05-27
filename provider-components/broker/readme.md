# About

The broker spawns up an http server (listening at **saaSProvisionerPort**) as the SaaS provisioning service, responsible for the tenant onboarding/offboarding.

It contains these folders and files, following our recommended project layout:

File or Folder | Purpose
---------|----------
`main.go/` | the implementation of the broker logic
`pkg/` | packages that are consumed by the GO application
`k8s/` | holds the manifest of the component
`Dockerfile` | to build the image of this component


## Configuration

To get the broker up and running, the `kymaDomain` parameter needs to point to the domain of your cluster. We've added the configmap `cluster-domain` before, and the value is taken from there.

