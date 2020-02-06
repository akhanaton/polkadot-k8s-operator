pushd ..
operator-sdk build ironoa/customresource-operator:v0.0.4
docker push ironoa/customresource-operator:v0.0.4
kubectl create -f deploy/operator.yaml
popd
