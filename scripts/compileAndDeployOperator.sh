pushd ..
operator-sdk build ironoa/customresource-operator:v0.0.6
docker push ironoa/customresource-operator:v0.0.6
kubectl create -f deploy/operator.yaml
popd
