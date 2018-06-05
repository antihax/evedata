curl -fsSL https://raw.githubusercontent.com/appscode/voyager/7.0.0/hack/deploy/voyager.sh \
    | bash -s -- --provider=baremetal --rbac

kubectl apply -f evedata.yaml