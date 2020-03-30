IMAGE_OPERATOR=ironoa/customresource-operator:v0.0.9 #define your favourite
# The above parameter have to match with the ones in the deployed resource defined in the deploy/operator.yaml file

K8S_OPERATOR=operator.yaml
K8S_CR=polkadot.swisscomblockchain.com_v1alpha1_polkadot_cr.yaml
K8S_CRD=polkadot.swisscomblockchain.com_polkadots_crd.yaml
K8S_SERVICE_ACCOUNT=service_account.yaml
K8S_ROLE=role.yaml
K8S_ROLE_BINDING=role_binding.yaml