if test -z "$DIR"
then
      DIR=$(dirname "${BASH_SOURCE[0]}")/..
      cd "$DIR" || exit
fi

IMAGE_METRICS=ironoa/polkadot-metrics:v0.0.1 # define your favourite

pushd .. >/dev/null 2>&1
docker build -t $IMAGE_METRICS ./metrics
docker push $IMAGE_METRICS
popd >/dev/null 2>&1 || exit
