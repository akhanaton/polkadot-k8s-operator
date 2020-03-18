if test -z "$DIR"
then
      DIR=$(dirname "${BASH_SOURCE[0]}")/..
      cd "$DIR" || exit
      source ./config/config.sh
fi

pushd .. >/dev/null 2>&1
docker build -t "$IMAGE_METRICS" ./metrics
docker push "$IMAGE_METRICS"
popd >/dev/null 2>&1 || exit
