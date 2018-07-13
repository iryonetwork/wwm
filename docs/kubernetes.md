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
$ kubectl create -f docs/k8s/registry-node-port.yaml
```

## 3. Publish latest container images to local registry

```
$ DOCKER_TAG=v0.0.1 DOCKER_REGISTRY="localhost:5000/" make publish
```

## 4. Set up dashboard (not required)

Deploy the dashboard:

```
kubectl create -f https://raw.githubusercontent.com/kubernetes/dashboard/ec1d7de4456e6a397c7f931f0a2bfc74a6ca2e9c/src/deploy/recommended/kubernetes-dashboard.yaml
```

Enable port forwarding:

```
$ DASHBOARD_NAME=$(kubectl get pods --namespace kube-system -l k8s-app=kubernetes-dashboard -o template --template '{{(index .items 0).metadata.name}}')
$ kubectl port-forward $DASHBOARD_NAME 8443:8443 --namespace=kube-system
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

To get the token allowing to access Kubernetes Dashboard run:

```
kubectl -n kube-system describe secrets kubernetes-dashboard-token
```

## 5. Setup helm

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

## 6. Set up cert-manager

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

The root certificate has to be added to your computer for you to be able to use `curl` without the `-k (insecure)` flag or to open `https://iryo.k8s.local` in the browser without the `This page is not secure` warnings. To install it:

```bash
security add-trusted-cert -k $HOME/Library/Keychains/login.keychain ca.crt
```

Finally, you should create a secret in kubernetes with the generated keys in namespaces: default, kube-system, local & cloud so all components can access it:

```
$ kubectl create secret tls ca-key-pair \
   --cert=ca.crt \
   --key=ca.key \
   --namespace=default
$ kubectl create namespace local
$ kubectl create secret tls ca-key-pair \
   --cert=ca.crt \
   --key=ca.key \
   --namespace=local
$ kubectl create namespace cloud
$ kubectl create secret tls ca-key-pair \
   --cert=ca.crt \
   --key=ca.key \
   --namespace=cloud
```

Next step is to create cert-manager issuers in each relevant namespace:

```
$ kubectl apply -f docs/k8s/ca-clusterissuer.yaml --namespace default
$ kubectl apply -f docs/k8s/ca-cloudissuer.yaml --namespace cloud
$ kubectl apply -f docs/k8s/ca-localissuer.yaml --namespace local
```

Test it by creating an example certificate:

```
$ kubectl apply -f docs/k8s/example-certificate.yaml
```

## 8. Install IRYO cloud & local charts

> Download IRYO helm charts from here: https://github.com/iryonetwork/charts
> Install the charts with values files for local deplyment:

*  Local IRYO node chart values: `docs/k8s/local-development-values.yaml`
*  Cloud IRYO node chart values:: `docs/k8s/cloud-development-values.yaml`

```
$ helm dependency update local && helm install --namespace local --values $GOPATH/src/github.com/iryonetwork/wwm/docs/k8s/local-development-values.yaml local -n local
$ helm dependency update cloud && helm install --namespace cloud --values $GOPATH/src/github.com/iryonetwork/wwm/docs/k8s/cloud-development-values.yaml cloud -n cloud
```

## 9. Copy secrets between namespaces

To ensure that sync of auth data and storage data is working you need to copy secrets from local namespace to cloud namespace.

```
$ kubectl get secrets -n local ca-authsync-local -o yaml | sed 's/namespace: local/namespace: cloud/'  | kubectl create -f -
$ kubectl get secrets -n local ca-storagesync-local -o yaml | sed 's/namespace: local/namespace: cloud/'  | kubectl create -f -
$ kubectl get secrets -n local ca-batchstoragesync-local -o yaml | sed 's/namespace: local/namespace: cloud/'  | kubectl create -f -
```

## 10. Provision Disocvery's database

For now you need to manually provision discovery data.
First you need to make it possible to access PostgreSQL DB.

Enable port forwarding for local-postgresql:

```
$ PSQL_POD_NAME=$(kubectl get pods --namespace local -l app=postgresql -o template --template '{{(index .items 0).metadata.name}}')
$ kubectl port-forward $PSQL_POD_NAME 5432:5432 --namespace=local
```

Connect to PostgreSQL using psql:

```
$ psql -U postgres -h localhost -p 5432
```

Initialize users with statements from `docs/k8s/localDiscoveryInit.sql` then connect to localdiscvoery db (`\c localdiscovery`) and create schema and data using statements from `docs/k8s/discoverySchemaAndData.sql`.

Do the same for cloudDiscovery replacing `local` with `cloud` in commands where applicable.

Enable port forwarding for cloud-postgresql:

```
$ PSQL_POD_NAME=$(kubectl get pods --namespace cloud -l app=postgresql -o template --template '{{(index .items 0).metadata.name}}')
$ kubectl port-forward $PSQL_POD_NAME 5432:5432 --namespace=cloud
```

Connect to PostgreSQL using psql:

```
$ psql -U postgres -h localhost -p 5432
```

Initialize users with statements from `docs/k8s/cloudDiscoveryInit.sql` then connect to localdiscvoery db (`\c clouddiscovery`) and create schema and data using statements from `docs/k8s/discoverySchemaAndData.sql`.

## 11. Install Traefik ingress controller

> Download IRYO helm charts from here: https://github.com/iryonetwork/charts
> Install `traefik` ingress controller chart in kube-system namespace with values for local deployment.

```
$ helm dependency update traefik && helm install --namespace default--values $GOPATH/src/github.com/iryonetwork/wwm/docs/k8s/traefik-development-values.yaml traefik -n traefik
```

## 12. Update `/etc/hosts` on your machine

```
127.0.0.1 iryo.k8s.local iryo.k8s.cloud traefik-dashboard.k8s.local
```

## 13. You should be able now to:

*  Access cloud dashboard frontend at `https://iryo.k8s.cloud` and cloud APIs at `https://iryo.k8s.cloud/api/v1/*` from your host;
*  Access clinic frontend at `https://iryo.k8s.local` and clinic APIs at `https://iryo.k8s.local/api/v1/*` from your host.

## 14. Upgrading

Publish new images to local registry and run:

```
$ helm dependency update local && helm upgrade --namespace local --values $GOPATH/src/github.com/iryonetwork/wwm/docs/k8s/local-development-values.yaml local local
$ helm dependency update cloud && helm upgrade --namespace cloud --values $GOPATH/src/github.com/iryonetwork/wwm/docs/k8s/cloud-development-values.yaml cloud cloud
```


