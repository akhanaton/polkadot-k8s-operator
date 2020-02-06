pushd ..
kubectl create -f deploy/service_account.yaml
kubectl create -f deploy/role.yaml
kubectl create -f deploy/role_binding.yaml
kubectl create -f deploy/crds/cache.example.com_customresources_crd.yaml
popd
source ./compileAndDeployOperator.sh
source ./deployCR.sh