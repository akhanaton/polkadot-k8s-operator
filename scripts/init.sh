pushd ..
kubectl create -f deploy/service_account.yaml
kubectl create -f deploy/role.yaml
kubectl create -f deploy/role_binding.yaml
kubectl create -f deploy/crds/polkadot.swisscomblockchain.com_polkadots_crd.yaml
popd
source ./compileAndDeployOperator.sh
source ./deployCR.sh