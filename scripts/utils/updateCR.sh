if test -z "$DIR"
then
      DIR=$(dirname "${BASH_SOURCE[0]}")/..
      cd "$DIR" || exit
fi

pushd .. >/dev/null 2>&1
kubectl apply -f deploy/crds/polkadot.swisscomblockchain.com_v1alpha1_polkadot_cr.yaml
popd >/dev/null 2>&1 || exit