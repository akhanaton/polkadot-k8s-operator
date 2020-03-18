if test -z "$DIR"
then
      DIR=$(dirname "${BASH_SOURCE[0]}")/..
      cd "$DIR" || exit
      source ./config/config.sh
fi

pushd .. >/dev/null 2>&1
operator-sdk build "$IMAGE_OPERATOR"
docker push "$IMAGE_OPERATOR"
kubectl create -f deploy/"$K8S_OPERATOR"
popd >/dev/null 2>&1 || exit
