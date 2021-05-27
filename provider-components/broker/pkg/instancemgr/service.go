package instancemgr

import (
	"context"
	"encoding/json"
	"fmt"
	"broker/pkg/config"

	"code.cloudfoundry.org/lager"
	apiRules "github.com/kyma-incubator/api-gateway/api/v1alpha1"
	oauthclientv1alpha1 "github.com/ory/hydra-maester/api/v1alpha1"
	rulev1alpha1 "github.com/ory/oathkeeper-maester/api/v1alpha1"
	"github.com/pivotal-cf/brokerapi/domain"
	v1 "k8s.io/api/core/v1"

	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	SaaSSubdomainLabel = "SaaSSubdomain"
	InstanceIDLabel    = "instanceID"
	BindingIDLabel     = "bindingID"
	StatusOK           = "OK"
	StatusFailed       = "Failed"
	SubscribedTenantID = "subscribedTenantId"
)

type Service struct {
	k8sClient         client.Client
	logger            lager.Logger
	brokerNamespace   string
	targetServiceName string
	targetServicePort uint32
	ingressGateway    string
	kymaDomain        string
}

type ServiceOperation struct {
	State       domain.LastOperationState `json:"operationState"`
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
}

type Instance struct {
	URL    string
	Status string
}

type Binding struct {
	URL           string
	OAuthTokenURL string
	Scopes        []string
	ClientID      string
	ClientSecret  string
}

func New(l lager.Logger, clnt client.Client, brokerNamespace string,
	targetServiceName string, targetServicePort uint32, ingressGateway string, kymaDomain string) *Service {

	return &Service{
		k8sClient:         clnt,
		logger:            l,
		ingressGateway:    ingressGateway,
		targetServiceName: targetServiceName,
		brokerNamespace:   brokerNamespace,
		targetServicePort: targetServicePort,
		kymaDomain:        kymaDomain,
	}
}

// CreateInstance creates an APIRule which exposes the targetService present in the provider cluster.
// targetService is the faq-ui component at port 8080. faq-ui is the odata service exposing Faq and it's related
// entity data
func (s *Service) CreateInstance(ctx context.Context, instanceID string,
	subdomain string) error {
	inst, err := s.GetInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("error performing existence check: %s", err.Error())
	}
	if inst != nil {
		return fmt.Errorf("instance %s already exists", instanceID)
	}

	apiRule := createAPIRule(instanceID, s.brokerNamespace, s.ingressGateway, s.targetServiceName, s.kymaDomain,
		s.targetServicePort, subdomain)

	err = s.k8sClient.Create(ctx, apiRule)
	if err != nil {
		s.logger.Error(config.ActionCreateServiceInstance, err, lager.Data{
			config.KeyMessage: "Error creating k8s APIRule object in provision",
		})
		return fmt.Errorf("error provisioning instance: %s", err.Error())
	}

	return nil
}

// CreateSaaSInstance creates and APIRule which exposes the saasService present in the provider cluster.
// saasService is the kyma-faq-ui-shell component at port 5000. kyma-faq-ui-shell is the Luigi UI for viewing
// and exporting FAQ Data.
func (s *Service) CreateSaaSInstance(ctx context.Context, kymaDomain, subdomain, serviceName string, servicePort uint32, tenantId string) (string, error) {
	inst, err := s.GetSaaSInstance(ctx, subdomain)
	if err != nil {
		return "", fmt.Errorf("error performing existence check: %s", err.Error())
	}
	if inst != nil {
		return "", fmt.Errorf("instance %s already exists", subdomain)
	}

	apiRule := createNoopAPIRule(s.brokerNamespace, s.ingressGateway, serviceName, kymaDomain, subdomain, servicePort, tenantId)

	err = s.k8sClient.Create(ctx, apiRule)
	if err != nil {
		s.logger.Error(config.ActionCreateSaasInstance, err, lager.Data{
			config.KeyMessage: "Error creating k8s APIRule object in provision",
		})
		return "", fmt.Errorf("error provisioning instance: %s", err.Error())
	}

	return *apiRule.Spec.Service.Host, nil
}

// Callback on Delete Service Instance. DeleteInstance will delete the APIRule created in CreateInstance
func (s *Service) DeleteInstance(ctx context.Context, instanceID string) error {

	inst, err := s.GetInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("error performing existence check: %s", err.Error())
	}
	//Nothing to do, no instance exists
	if inst == nil {
		return nil
	}

	err = s.k8sClient.DeleteAllOf(ctx, &apiRules.APIRule{},
		client.MatchingLabels{InstanceIDLabel: instanceID},
		client.InNamespace(s.brokerNamespace))
	if err != nil {
		s.logger.Error(config.ActionDeleteServiceInstance, err, lager.Data{
			config.KeyMessage: "Error deleting k8s APIRule object in deprovision",
		})
		return fmt.Errorf("error deprovisioning instance %s: %s", instanceID, err.Error())
	}

	return nil
}

// DeleteSaasInstance will delete APIRule created in CreateSaaSInstance for tenant offboarding.
func (s *Service) DeleteSaasInstance(ctx context.Context, subDomain string) error {
	inst, err := s.GetSaaSInstance(ctx, subDomain)
	if err != nil {
		return fmt.Errorf("error performing existence check: %s", err.Error())
	}
	//Nothing to do, no instance exists
	if inst == nil {
		return nil
	}

	err = s.k8sClient.DeleteAllOf(ctx, &apiRules.APIRule{},
		client.MatchingLabels{SaaSSubdomainLabel: subDomain},
		client.InNamespace(s.brokerNamespace))
	if err != nil {
		s.logger.Error(config.ActionDeleteSaasInstance, err, lager.Data{
			config.KeyMessage: "Error deleting k8s apirule object in deprovision",
		})
		return fmt.Errorf("error deprovisioning saas instance %s: %s", subDomain, err.Error())
	}

	return nil
}

// Retrieve Service Instance via instanceID Label, internnally performs a retrival of an API Rule
func (s *Service) GetInstance(ctx context.Context, instanceID string) (*Instance, error) {
	return s.getInstanceViaNamedLabel(ctx, InstanceIDLabel, instanceID)
}

func (s *Service) GetSaaSInstance(ctx context.Context, subdomain string) (*Instance, error) {
	return s.getInstanceViaNamedLabel(ctx, SaaSSubdomainLabel, subdomain)
}

func (s *Service) getInstanceViaNamedLabel(ctx context.Context, labelName, labelValue string) (*Instance, error) {
	apiRuleList := &apiRules.APIRuleList{}
	err := s.k8sClient.List(ctx, apiRuleList, client.MatchingLabels{
		labelName: labelValue,
	}, client.InNamespace(s.brokerNamespace))
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	if len(apiRuleList.Items) > 1 {
		return nil, fmt.Errorf("More than one API Rule for label - %s [namespace: %s]", labelValue, s.brokerNamespace)
	} else if len(apiRuleList.Items) == 0 {
		return nil, nil
	}

	status := StatusOK
	if apiRuleList.Items[0].Status.APIRuleStatus.Code != apiRules.StatusOK {
		status = StatusFailed
	}

	return &Instance{
		URL:    *apiRuleList.Items[0].Spec.Service.Host,
		Status: status,
	}, nil
}

// Callback on Create Service Binding. CreateBinding creates an OAuth client for an APIRule, with scope defined as
// instanceId, so that every service instance gets it's own set of client credentials are confined to their scope,
// which will restrict their privileges to invokoe other APIRules, in case of collision.
func (s *Service) CreateBinding(ctx context.Context, instanceID, bindingID string) (*Binding, error) {
	oauthClient, err := s.getOAuth2Client(ctx, instanceID, bindingID)
	if err != nil {
		s.logger.Error(config.ActionCreateServiceBinding, err, lager.Data{
			config.KeyMessage: "error reading binding",
		})
		return nil, fmt.Errorf("error reading OAuth2 client: %s", err.Error())
	}
	if oauthClient != nil {
		s.logger.Error(config.ActionCreateServiceBinding,
			fmt.Errorf("oauth2 client with name %s in namespace %s already exists",
				createK8SBindingName(bindingID), s.brokerNamespace))
		return nil, fmt.Errorf("oauth2 client with name %s in namespace %s already exists",
			createK8SBindingName(bindingID), s.brokerNamespace)
	}

	secret := createOauth2Secret(instanceID, bindingID, s.brokerNamespace)
	err = s.k8sClient.Create(ctx, secret)
	if err != nil {
		s.logger.Error(config.ActionCreateServiceBinding, err, lager.Data{
			config.KeyMessage: "Error creating OAuth2 Secret in cluster",
		})
		return nil, fmt.Errorf("error provisioning binding: %s", err.Error())
	}

	oauthClient = createOAuth2Client(instanceID, bindingID, s.brokerNamespace)
	err = s.k8sClient.Create(ctx, oauthClient)
	if err != nil {
		s.logger.Error(config.ActionCreateServiceBinding, err, lager.Data{
			config.KeyMessage: "Error creating OAuth2 Client in cluster"})
		return nil, fmt.Errorf("error provisioning binding: %s", err.Error())
	}

	return &Binding{
		URL:           fmt.Sprintf("https://%s.%s/", createInstanceHost(instanceID), s.kymaDomain),
		OAuthTokenURL: fmt.Sprintf("https://oauth2.%s/oauth2/token", s.kymaDomain),
		Scopes: []string{
			instanceID,
		},
		ClientID:     string(secret.StringData["client_id"]),
		ClientSecret: string(secret.StringData["client_secret"]),
	}, nil
}

// Delete the OAuth Client created in CreateBinding
func (s *Service) DeleteBinding(ctx context.Context, instanceID, bindingID string) error {
	oauthClient, err := s.getOAuth2Client(ctx, instanceID, bindingID)
	if err != nil {
		s.logger.Error(config.ActionDeleteServiceBinding, err)
		return fmt.Errorf("error reading OAuth2 client: %s", err.Error())
	}
	// since it is gone already
	if oauthClient == nil {
		return nil
	}

	err = s.k8sClient.DeleteAllOf(ctx, &oauthclientv1alpha1.OAuth2Client{},
		client.MatchingLabels{InstanceIDLabel: instanceID, BindingIDLabel: bindingID},
		client.InNamespace(s.brokerNamespace))
	if err != nil {
		s.logger.Error(config.ActionDeleteServiceBinding, err, lager.Data{
			config.KeyMessage: "error deleting OAuth2 client",
		})
		return fmt.Errorf("error deleting binding %s (instance %s, namespace %s): %s",
			bindingID, instanceID, s.brokerNamespace, err.Error())
	}

	secret, err := s.getSecret(ctx, oauthClient.Spec.SecretName)
	if err != nil {
		s.logger.Error(config.ActionDeleteServiceBinding, err)
		return fmt.Errorf("error reading secret: %s", err.Error())
	}
	// since it is gone already
	if oauthClient == nil {
		return nil
	}

	if err := s.k8sClient.Delete(ctx, secret); err != nil {
		s.logger.Error(config.ActionDeleteServiceBinding, err)
		return fmt.Errorf("error deleting secret: %s", err.Error())
	}

	return nil
}

// Geth the OAuthClient for an APiRule which is per ServiceInstance.
func (s *Service) GetBinding(ctx context.Context, instanceID, bindingID string) (*Binding, error) {
	oauthClient, err := s.getOAuth2Client(ctx, instanceID, bindingID)
	if err != nil {
		s.logger.Error(config.ActionGetServiceBinding, err)
		return nil, fmt.Errorf("error reading OAuth2 client: %s", err.Error())
	}
	if oauthClient == nil {
		return nil, nil
	}
	secret, err := s.getSecret(ctx, oauthClient.Spec.SecretName)
	if err != nil {
		s.logger.Error(config.ActionGetServiceBinding, err)
		return nil, fmt.Errorf("error reading client credential secret: %s", err.Error())
	}
	if secret == nil {
		return nil, nil
	}

	return &Binding{
		URL:           fmt.Sprintf("https://%s.%s/", createInstanceHost(instanceID), s.kymaDomain),
		OAuthTokenURL: fmt.Sprintf("https://oauth2.%s/oauth2/token", s.kymaDomain),
		Scopes: []string{
			instanceID,
		},
		ClientID:     string(secret.Data["client_id"]),
		ClientSecret: string(secret.Data["client_secret"]),
	}, nil
}

// helper function to retrieve oauthclientv1alpha1.OAuth2Client Custom Resource.
func (s *Service) getOAuth2Client(ctx context.Context, instanceID, bindingID string) (*oauthclientv1alpha1.OAuth2Client,
	error) {

	oAuth2ClientList := &oauthclientv1alpha1.OAuth2ClientList{}
	err := s.k8sClient.List(ctx, oAuth2ClientList, client.MatchingLabels{
		InstanceIDLabel: instanceID,
		BindingIDLabel:  bindingID,
	}, client.InNamespace(s.brokerNamespace))
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	if len(oAuth2ClientList.Items) > 1 {
		return nil, fmt.Errorf("more than one Client for binding Id - %s (instance %s, namespace: %s)",
			bindingID, instanceID, s.brokerNamespace)
	} else if len(oAuth2ClientList.Items) == 0 {
		return nil, nil
	}

	return &oAuth2ClientList.Items[0], nil
}

// helper function to retrieve a v1.Secret Custom Resource.
func (s *Service) getSecret(ctx context.Context, name string) (*v1.Secret,
	error) {

	secret := &v1.Secret{}
	err := s.k8sClient.Get(ctx, client.ObjectKey{Name: name, Namespace: s.brokerNamespace}, secret)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return secret, nil
}

// helper function to create a v1.Secret Custom Resource.
func createOauth2Secret(instanceID string, bindingID string, namespace string) *v1.Secret {
	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      createK8SBindingName(bindingID),
			Labels: map[string]string{
				InstanceIDLabel: instanceID,
				BindingIDLabel:  bindingID,
			},
		},
		StringData: map[string]string{
			"client_id":     uuid.New().String(),
			"client_secret": uuid.New().String(),
		},
	}
}
// helper function to create a oauthclientv1alpha1.OAuth2Client Custom Resource.
func createOAuth2Client(instanceID string, bindingID string, namespace string) *oauthclientv1alpha1.OAuth2Client {

	return &oauthclientv1alpha1.OAuth2Client{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      createK8SBindingName(bindingID),
			Labels: map[string]string{
				InstanceIDLabel: instanceID,
				BindingIDLabel:  bindingID,
			},
		},
		Spec: oauthclientv1alpha1.OAuth2ClientSpec{
			GrantTypes: []oauthclientv1alpha1.GrantType{
				oauthclientv1alpha1.GrantType("client_credentials"),
			},
			ResponseTypes: []oauthclientv1alpha1.ResponseType{
				oauthclientv1alpha1.ResponseType("token"),
			},
			Scope:      instanceID,
			SecretName: createK8SBindingName(bindingID),
			HydraAdmin: oauthclientv1alpha1.HydraAdmin{},
		},
	}
}

// Helper function to create apiRules.APIRule Custom Resource, with OAUth.
func createAPIRule(instanceID, namespace, ingressGateway,
	k8sServiceName, kymaDomain string, targetServicePort uint32, subdomain string) *apiRules.APIRule {

	instanceHost := createInstanceHost(instanceID)

	addedHeaders := map[string]string{
		"x-forwarded-host": fmt.Sprintf("%s.%s", instanceHost, kymaDomain),
	}

	if subdomain != "" {
		addedHeaders["subdomain"] = subdomain
	}
	addedHeadersRaw, _ := json.Marshal(map[string]interface{}{"headers": addedHeaders})

	apiRule := apiRules.APIRule{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      fmt.Sprintf("inst-%s", instanceID),
			Labels: map[string]string{
				InstanceIDLabel: instanceID,
			},
		},
		Spec: apiRules.APIRuleSpec{
			Rules: []apiRules.Rule{
				{
					Methods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
					Path:    "/.*",
					AccessStrategies: []*rulev1alpha1.Authenticator{
						{
							Handler: &rulev1alpha1.Handler{
								Name:   "oauth2_introspection",
								Config: &runtime.RawExtension{Raw: createRequiredScopes([]string{instanceID})},
							},
						},
					},
					Mutators: []*rulev1alpha1.Mutator{
						{
							Handler: &rulev1alpha1.Handler{
								Name:   "header",
								Config: &runtime.RawExtension{Raw: addedHeadersRaw},
							},
						},
					},
				},
			},
			Service: &apiRules.Service{
				Name: &k8sServiceName,
				Port: &targetServicePort,
				Host: &instanceHost,
			},
			Gateway: &ingressGateway,
		},
	}
	return &apiRule
}

// Helper function to create apiRules.APIRule Custom Resource, without OAUth.
func createNoopAPIRule(namespace, ingressGateway, serviceName, kymaDomain, subdomain string, servicePort uint32, tenantId string) *apiRules.APIRule {
	instanceHost := createSaaSInstanceHost(subdomain)
	apiRule := apiRules.APIRule{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      instanceHost,
			Labels: map[string]string{
				SaaSSubdomainLabel: subdomain,
				SubscribedTenantID: tenantId,
			},
		},
		Spec: apiRules.APIRuleSpec{
			Rules: []apiRules.Rule{
				{
					Methods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
					Path:    "/.*",
					AccessStrategies: []*rulev1alpha1.Authenticator{
						{
							Handler: &rulev1alpha1.Handler{
								Name: "noop",
							},
						},
					},
					Mutators: []*rulev1alpha1.Mutator{
						{
							Handler: &rulev1alpha1.Handler{
								Name:   "header",
								Config: &runtime.RawExtension{Raw: createMutator(instanceHost, kymaDomain)},
							},
						},
					},
				},
			},
			Service: &apiRules.Service{
				Name: &serviceName,
				Port: &servicePort,
				Host: &instanceHost,
			},
			Gateway: &ingressGateway,
		},
	}
	return &apiRule
}

func createRequiredScopes(scopes []string) []byte {
	requiredScopes := struct {
		// Array of required scopes
		RequiredScope []string `json:"required_scope"`
	}{
		RequiredScope: scopes,
	}
	requiredScopesJSON, _ := json.Marshal(requiredScopes)

	return requiredScopesJSON
}

func createMutator(name, domain string) []byte {
	headers := struct {
		// Array of required scopes
		Header struct {
			XForwardedHost string `json:"x-forwarded-host"`
		} `json:"headers"`
	}{
		Header: struct {
			XForwardedHost string `json:"x-forwarded-host"`
		}{
			XForwardedHost: fmt.Sprintf("%s.%s", name, domain),
		},
	}
	headerJSON, _ := json.Marshal(headers)

	return headerJSON
}

func createK8SBindingName(bindingID string) string {
	return fmt.Sprintf("bnd-%s", bindingID)
}

func createInstanceHost(instanceId string) string {
	return fmt.Sprintf("svc-%s", instanceId)
}

func createSaaSInstanceHost(instanceId string) string {
	return fmt.Sprintf("%s", instanceId)
}
