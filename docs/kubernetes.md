# Working with Kubernetes locally

## 1. Set up kubernetes

Either install edge version of `docker-for-mac` (latest version working properly is `18.03.0-ce-mac58`, mind that `18.03.0-ce-mac62` is broken) or set up minikube.

You'll also require `helm`.

## 2. Get registry up and running

This will start up a local registry server inside kubernetes.

```
$ kubectl create -f https://raw.githubusercontent.com/giantswarm/kubernetes-registry/f81b783/manifests-all.yaml
```

With docker-for-mac registry (localhost:5000) should already be visible to the docker engine.

If needed you can set port forwarding:

```
$ REGISTRY_NAME=$(kubectl get pods --namespace registry -l app=registry,component=main -o template --template '{{(index .items 0).metadata.name}}') \
kubectl --namespace registry port-forward $REGISTRY_NAME 5000:5000
```

You can test it:

```
$ docker pull alpine
$ docker tag alpine localhost:5000/alpine
$ docker push localhost:5000/alpine
```

Alternatively you can set up permanent port forwarding with nodeport to port 30500:

```
kubectl create -f docs/k8s/registry-node-port.yaml
```

## 3. Set up dashboard (not required)

Deploy the dashboard:

```
https://raw.githubusercontent.com/kubernetes/dashboard/ec1d7de4456e6a397c7f931f0a2bfc74a6ca2e9c/src/deploy/recommended/kubernetes-dashboard.yaml
```

Enable port forwarding:

```
$ DASHBOARD_NAME=$(kubectl get pods --namespace kube-system -l k8s-app=kubernetes-dashboard -o template --template '{{(index .items 0).metadata.name}}') \
kubectl port-forward $DASHBOARD_NAME 8443:8443 --namespace=kube-system
```

Alternatevely you can setup permanent node port to port 30443 with an included node port definition

```
kubectl create -f docs/k8s/dashboard-node-port.yaml
```

And test it:

```
$ docker pull alpine
$ docker tag alpine localhost:30500/alpine
$ docker push localhost:30500/alpine
```

## 4. Setup helm

Init helm

```
helm init
```

Create a service account

```
kubectl delete -f docs/k8s/tiller-serviceAccount.yaml
kubectl create -f docs/k8s/tiller-serviceAccount.yaml
```

Upgrade tiller

```
helm init --upgrade
```

## 5. Set up cert-manager

```
helm install --name cert-manager stable/cert-manager
```

Next we need to prepare certificates:

```
# Generate a CA private key
$ openssl genrsa -out ca.key 2048

# Create a self signed Certificate, valid for 10yrs with the 'signing' option set
$ openssl req -x509 -new -nodes -key ca.key -subj "/CN=${COMMON_NAME}" -days 3650 -reqexts v3_req -extensions v3_ca -out ca.crt
```

Create a secret in kubernetes with the generated keys

```
$ kubectl create secret tls ca-key-pair \
   --cert=ca.crt \
   --key=ca.key \
   --namespace=default
```

Create an issuer

```
$ kubectl apply -f docs/k8s/ca-cloudissuer.yaml --namespace cloud
$ kubectl apply -f docs/k8s/ca-localissuer.yaml --namespace local
```

Test it by creating an example certificate:

```
$ kubectl apply -f docs/k8s/example-certificate.yaml
$
```

## 6. How to copy secrets between namespaces

```
kubectl get secrets -n local ca-authsync-local -o yaml | sed 's/namespace: local/namespace: cloud/'  | kubectl create -f -
```

## 7. Set up vault

> Follows these instructions https://github.com/Boostport/kubernetes-vault/blob/master/deployments/quick-start/README.md

First deploy the vault:

```
kubectl apply -f docs/k8s/vault.yaml --namespace kube-system
```

Set up port forwarding:

```
VAULT_NAME=$(kubectl get pods --namespace kube-system -l app=vault -o template --template  '{{(index .items 0).metadata.name}}') \
kubectl port-forward $VAULT_NAME 8200
```
