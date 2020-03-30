if test -z "$DIR"
then
      DIR=$(dirname "${BASH_SOURCE[0]}")
      cd "$DIR" || exit
      source ./config/config.sh
fi

source ./utils/wipeCR.sh
source ./utils/wipeOperator.sh

pushd .. >/dev/null 2>&1
kubectl delete -f deploy/crds/"$K8S_CRD"
kubectl delete -f deploy/"$K8S_ROLE_BINDING"
kubectl delete -f deploy/"$K8S_ROLE"
kubectl delete -f deploy/"$K8S_SERVICE_ACCOUNT"
popd >/dev/null 2>&1 || exit