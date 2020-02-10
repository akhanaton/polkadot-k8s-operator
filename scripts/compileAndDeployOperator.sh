pushd ..
operator-sdk build ironoa/customresource-operator:v0.0.5
docker push ironoa/customresource-operator:v0.0.5
kubectl create -f deploy/operator.yaml
popd
