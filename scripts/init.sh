if test -z "$DIR"
then
      DIR=$(dirname "${BASH_SOURCE[0]}")
fi
cd "$DIR" || exit

pushd .. >/dev/null 2>&1
kubectl create -f deploy/service_account.yaml
kubectl create -f deploy/role.yaml
kubectl create -f deploy/role_binding.yaml
kubectl create -f deploy/crds/polkadot.swisscomblockchain.com_polkadots_crd.yaml
popd >/dev/null 2>&1 || exit

source ./utils/buildAndDeployMetrics.sh
source ./utils/compileAndDeployOperator.sh
source ./utils/deployCR.sh