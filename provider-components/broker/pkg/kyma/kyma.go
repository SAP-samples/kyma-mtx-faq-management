package kyma

import (
	"bytes"
	"encoding/json"
	"fmt"
	"broker/pkg/config"
	"broker/pkg/instancemgr"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
	"code.cloudfoundry.org/lager"
)

type KymaHandler struct {
	logger             lager.Logger
	service            *instancemgr.Service
	namespace          string
	kymaDomain         string
	serviceName        string
	servicePort        uint32
	capProvisioningURL string
	client             *http.Client
}

func NewKymaHandler(l lager.Logger, k8sService *instancemgr.Service, namespace string, kymaDomain string,
	serviceName string, servicePort uint32, capProvisioningURL string) *KymaHandler {
	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	return &KymaHandler{
		logger:             l,
		service:            k8sService,
		namespace:          namespace,
		serviceName:        serviceName,
		servicePort:        servicePort,
		kymaDomain:         kymaDomain,
		capProvisioningURL: capProvisioningURL,
		client:             client,
	}
}

// saas provisioner service handler, Handles the tenant onboarding/offboarding requests
func (k *KymaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Ignore all the requests besides a PUT request
	if !(r.Method == http.MethodPut || r.Method == http.MethodDelete) {
		writeResponse(http.StatusOK, "message", "Request Ignored", w, r)
		return
	}

	var request map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeResponse(http.StatusBadRequest, "message", "Valid Json request expected", w, r)
		return
	}

	// Extract subdomain attribute from the onboarding request
	subdomain, ok := request["subscribedSubdomain"].(string)
	if !ok || len(subdomain) < 1 {
		writeResponse(http.StatusBadRequest, "message", "Subscriber Subdomain is not present in the request", w, r)
		return
	}

	// Extract tenantID attribute from the onboarding request
	tenantId, ok := request["subscribedTenantId"].(string)
	if !ok || len(tenantId) < 1 {
		writeResponse(http.StatusBadRequest, "message", "Tenant Id is not present in the request", w, r)
		return
	}

	targetCapUrl := fmt.Sprintf("%s%s%s", strings.TrimSuffix(k.capProvisioningURL, "/"), "/mtx/v1/provisioning/tenant/", tenantId)

	reqBytes, _ := json.Marshal(request)
	// Delegate request to targetCapUrl so that a tenant can be created in the CAP application on behalf of the subaccount
	req, err := http.NewRequest(r.Method, targetCapUrl, bytes.NewBuffer(reqBytes))
	if err != nil {
		writeResponse(http.StatusInternalServerError, "message", "Unable to form Cap Provisioning Request", w, r)
		return
	}
	if cType := r.Header.Get("Content-Type"); cType != "" {
		req.Header.Set("Content-Type", cType)
	}
	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	res, err := k.client.Do(req)
	if err != nil {
		writeResponse(http.StatusInternalServerError, "message", "Unable to Provision CAP Application Request", w, r)
		return
	}
	dump, err := httputil.DumpResponse(res, true)
	if err != nil {
		writeResponse(http.StatusInternalServerError, "message", "Unable to Dump Response", w, r)
		return
	}
	fmt.Printf("CAP Response: %q", dump)

	if r.Method == http.MethodPut {
		k.logger.Info(config.ActionCreateSaasInstance, lager.Data{
			config.KeyMessage:   "Creating APIRule for subdomain",
			config.KeySubdomain: subdomain,
		})

		// Create an APIRule which is specific to a service instance, every service instance get their own
		// APIRule. the target service for the APIRule is always the same, which is the kyma-faq-ui-shell service
		hostname, err := k.service.CreateSaaSInstance(r.Context(), k.kymaDomain, subdomain, k.serviceName, k.servicePort, tenantId)
		if err != nil {
			k.logger.Error(config.ActionCreateSaasInstance, err, lager.Data{
				config.KeyMessage:   "Error creating APIRule for subdomain",
				config.KeySubdomain: subdomain,
			})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		k.logger.Info(config.ActionCreateSaasInstance, lager.Data{
			config.KeyStatus: config.StatusOk,
		})
		w.Write([]byte(fmt.Sprintf("https://%s.%s", hostname, k.kymaDomain)))
		w.WriteHeader(http.StatusOK)
	} else if r.Method == http.MethodDelete {
		fmt.Printf("\nrequest Body: %+v\n", request)

		k.logger.Info(config.ActionDeleteSaasInstance, lager.Data{
			config.KeyMessage:   "Deleting APIRule for subdomain",
			config.KeySubdomain: subdomain,
		})

		err := k.service.DeleteSaasInstance(r.Context(), subdomain)
		if err != nil {
			k.logger.Error(config.ActionDeleteSaasInstance, err, lager.Data{
				config.KeyMessage:   "Error Deleting APIRule for subdomain",
				config.KeySubdomain: subdomain,
			})
			return
		}
		messageString := fmt.Sprintf("APIRUle Deleted for subdomain: %s", subdomain)
		k.logger.Info(config.ActionDeleteSaasInstance, lager.Data{
			config.KeyMessage: messageString,
		})
		writeResponse(http.StatusOK, "message", messageString, w, r)
	}
}

// Helper method to write http response with status code and message
func writeResponse(status int, qualifier string, value interface{}, w http.ResponseWriter, r *http.Request) {
	if strings.Contains(strings.ToLower(strings.Join(r.Header["Accept"], " ")), "application/json") {
		result, _ := json.Marshal(map[string]interface{}{qualifier: value})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write(result)
	} else {
		w.Header().Set("message", fmt.Sprintf("%s - %v", qualifier, value))
		w.WriteHeader(status)
	}
}
