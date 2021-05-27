package client

import (
	"encoding/json"
	"fmt"
	"github.com/sap-samples/csv-exporter/pkg/api"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type ODataClient interface {
	GetData(header http.Header) *api.FaqResponse
	HandleCSVExport(writer http.ResponseWriter, request *http.Request)
}

func NewOdataClient(config *api.EnvConfig) ODataClient {
	client := &http.Client{
		Timeout: 45 * time.Second,
	}
	mySecondString := "state eq 'answered'"
	t := &url.URL{Path: mySecondString}
	mySecondEncodedString := t.String()

	return &odataClient{
		url: fmt.Sprintf("%s%s%s", strings.TrimSuffix(config.OdataUrl, "/"),
			"/admin/Faqs?$filter=", mySecondEncodedString),
		client: client,
	}
}

type odataClient struct {
	url    string
	client *http.Client
}

func (o *odataClient) HandleCSVExport(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		response := o.GetData(request.Header)
		if response != nil {
			writer.Header().Set("Content-Type", "text/csv")
			writer.Header().Set("Content-Disposition", "attachment;filename=data.csv")
			if _, err := writer.Write(response.GetCSVBytes()); err != nil {
				logrus.Warningf("Error while writing CSV.\n err: %v", err)
			}
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
		}
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (o *odataClient) GetData(header http.Header) *api.FaqResponse {
	req, err := http.NewRequest(http.MethodGet, o.url, nil)

	req.Header.Set("Authorization", header.Get("Authorization"))

	res, err := o.client.Do(req)
	if err != nil {
		logrus.Warningf("Error Retrieving data from the service, err: %+v", err)
		return nil
	}
	dump, _ := httputil.DumpResponse(res, true)
	fmt.Printf("Response Dump : %q\n", dump)

	var responseData api.FaqResponse

	if err := json.NewDecoder(res.Body).Decode(&responseData); err != nil {
		logrus.Warningf("Error in Unmarshalling")
	}
	return &responseData
}
