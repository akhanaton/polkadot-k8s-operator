if test -z "$DIR"
then
      DIR=$(dirname "${BASH_SOURCE[0]}")/..
      cd "$DIR" || exit
fi

pushd .. >/dev/null 2>&1
operator-sdk build ironoa/customresource-operator:v0.0.8 # define your favourite
docker push ironoa/customresource-operator:v0.0.8 #define your favourite
kubectl create -f deploy/operator.yaml
popd >/dev/null 2>&1 || exit
