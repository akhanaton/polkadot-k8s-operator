if test -z "$DIR"
then
      DIR=$(dirname "${BASH_SOURCE[0]}")
fi
cd "$DIR" || exit

source ./utils/wipeOperator.sh
source ./utils/wipeCR.sh

pushd .. >/dev/null 2>&1
kubectl delete -f deploy/crds/polkadot.swisscomblockchain.com_polkadots_crd.yaml
kubectl delete -f deploy/role_binding.yaml
kubectl delete -f deploy/role.yaml
kubectl delete -f deploy/service_account.yaml
popd >/dev/null 2>&1 || exit