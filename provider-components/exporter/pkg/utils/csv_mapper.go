package utils

import (
	"bytes"
	"encoding/csv"
	"github.com/sirupsen/logrus"
)

func GetCSVBytes() []byte {
	record := [][]string{
		{"Question", "Answer"},
		{"How much does Kyma cost?", "The Open Source Project is free. The SAP Cloud Platform service can be calculated here: https://discovery-center.cloud.sap/serviceCatalog/kyma-runtime\n"},
		{"Does Kyma contain a Service Mesh?\n", "The Istio Service Mesh is part of every Kyma Deployment\n"},
		{"Is Kyma SAP's version of Kubernetes?\n", "Kyma is based on Kubernetes and integrates other relatedd Open Source projects into a coherent offering for building enterprise ready extensions and business applications.\n"},
	} // just some test data to use for the wr.Writer() method below.

	b := &bytes.Buffer{}
	wr := csv.NewWriter(b)
	for _, data := range record {
		if err := wr.Write(data); err != nil {
			logrus.Warning("Error in writing CSV")
		}
	}
	wr.Flush()

	return b.Bytes()
}


