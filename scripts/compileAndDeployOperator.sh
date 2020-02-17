pushd ..
operator-sdk build ironoa/customresource-operator:v0.0.7 # define your favourite
docker push ironoa/customresource-operator:v0.0.7 #define your favourite
kubectl create -f deploy/operator.yaml
popd
