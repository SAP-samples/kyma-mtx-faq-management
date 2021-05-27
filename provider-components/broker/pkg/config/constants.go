package config

const (
	KeyMessage = "message"

	KeySubdomain = "subdomain"

	KeyStatus = "status"

	StatusOk = "Success"

	ActionParseConfig = "parse-config"

	ActionK8sClientCreation = "k8s-client-creation"

	ActionStartSaasProvisioner = "start-saas-provisioner"

	ActionStartServiceBroker = "start-service-broker"

	ActionCreateSaasInstance = "create-saas-instance"

	ActionDeleteSaasInstance = "delete-saas-instance"

	ActionCreateServiceInstance = "create-service-instance"

	ActionDeleteServiceInstance = "delete-service-instance"

	ActionGetServiceInstance = "get-service-instance"

	ActionUpdateServiceInstance = "update-service-instance"

	ActionGetLastOperation = "get-last-operation"

	ActionSetLastOperation = "set-last-operation"

	ActionGetLastBindingOperation = "get-last-binding-operation"

	ActionSetLastBindingOperation = "set-last-binding-operation"

	ActionGetServices = "get-services"

	ActionGetServiceBinding = "get-service-binding"

	ActionCreateServiceBinding = "create-service-binding"

	ActionDeleteServiceBinding = "delete-service-binding"

	ActionCreateAPIRule = "create-api-rule"

	ActionDeleteAPIRule = "delete-api-rule"

	ActionTerminate = "shut-down"
)
