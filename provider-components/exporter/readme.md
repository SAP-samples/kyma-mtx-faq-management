# About

This service exports all the FAQs into a CSV File. The `user-interface` component calls this service.

It contains these folders and files, following our recommended project layout:

File or Folder | Purpose
---------|----------
`cmd/csv-service/main.go` | the implementation of the exporter
`pkg/` | packages that are consumed by the `main.go` application
`k8s/` | holds the manifest of the component
`Dockerfile` | to build the image of this component

## Usage
This exporter needs an environment variable (`CSV_ODATA_URL`) that points to the `backend` service to retrieve FAQs. The exporter will use this information to export the FAQs to a CSV file. 
The `exporter` service always needs to be called with a valid,  XSUAA-issued `Authorization` token, and this token is propagated to the `backend`.


```shell script
export CSV_ODATA_URL=<ODATA_SERVICE_URL>
```

### Run Locally 
```shell script
go run ./cmd/csv-service/main.go
```

### Test
If running locally:
```shell script
curl --location --request GET 'localhost:8081/getCSV' \
--header 'Authorization: Bearer Token'
```

If running on Kyma:
```shell script
curl --location --request GET 'https://csv-api.<CLUSTER_DOMAIN>/getCSV' \
--header 'Authorization: Bearer Token'
```
