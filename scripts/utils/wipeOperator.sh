if test -z "$DIR"
then
      DIR=$(dirname "${BASH_SOURCE[0]}")/..
      cd "$DIR" || exit
fi

pushd .. >/dev/null 2>&1
kubectl delete -f deploy/operator.yaml
popd >/dev/null 2>&1 || exit
