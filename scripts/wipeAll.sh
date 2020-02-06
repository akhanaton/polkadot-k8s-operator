source ./wipeOperator.sh
source ./wipeCR.sh
pushd ..
kubectl delete -f deploy/role_binding.yaml
kubectl delete -f deploy/role.yaml
kubectl delete -f deploy/service_account.yaml
popd