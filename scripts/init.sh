if test -z "$DIR"
then
      DIR=$(dirname "${BASH_SOURCE[0]}")
      cd "$DIR" || exit
      source ./config/config.sh
      source ./wipeAll.sh 2>/dev/null
fi

pushd .. >/dev/null 2>&1
kubectl create -f deploy/"$K8S_SERVICE_ACCOUNT"
kubectl create -f deploy/"$K8S_ROLE"
kubectl create -f deploy/"$K8S_ROLE_BINDING"
kubectl create -f deploy/crds/"$K8S_CRD"
popd >/dev/null 2>&1 || exit

source ./utils/compileAndDeployOperator.sh
source ./utils/deployCR.sh