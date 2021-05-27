module broker

go 1.16

require (
	code.cloudfoundry.org/lager v2.0.0+incompatible
	github.com/drewolson/testflight v1.0.0 // indirect
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/kyma-incubator/api-gateway v0.0.0-20201127140450-8af556cde95f
	github.com/ory/hydra-maester v0.0.19
	github.com/ory/oathkeeper-maester v0.1.0
	github.com/pivotal-cf/brokerapi v6.4.2+incompatible
	k8s.io/api v0.19.4
	k8s.io/apimachinery v0.19.4
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	sigs.k8s.io/controller-runtime v0.6.4
)

replace k8s.io/client-go => k8s.io/client-go v0.19.4 // incompatibility fix
