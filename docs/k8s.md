# Kubernetes Local Setup (Docker Desktop)

This guide explains how to run the Kafka practice stack on Docker Desktop Kubernetes using a local NGINX ingress controller.

## 1. Install the NGINX Ingress Controller

Docker Desktop does not install the ingress controller automatically, so install it manually:

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.9.4/deploy/static/provider/cloud/deploy.yaml
```

Wait until the ingress controller is ready:

```bash
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=180s
```

If you want to remove the ingress controller later, run:

```bash
kubectl delete -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.9.4/deploy/static/provider/cloud/deploy.yaml
```

## 2. Update `/etc/hosts` with Ingress Hosts

The Kubernetes ingress in this repo uses the following hostnames:

- `kafka-ui.k8s.local`
- `producer.k8s.local`

Add them to your local `/etc/hosts` file so your browser and CLI use the ingress hostnames:

```text
127.0.0.1 kafka-ui.k8s.local producer.k8s.local
```

On macOS, edit the file with sudo privileges:

```bash
sudo nano /etc/hosts
```

Save the file and then verify the entries with:

```bash
ping -c 1 kafka-ui.k8s.local
ping -c 1 producer.k8s.local
```

## 3. Build & Apply the Kubernetes Resources

The repo includes a convenience script to apply the manifests in the correct order.

```bash
chmod +x scripts/k8s-apply.sh
./scripts/k8s-apply.sh
```

This script applies:

- `k8s/namespace.yaml`
- `k8s/kafka.yaml`
- `k8s/schema_registry.yaml`
- `k8s/kafka_connect.yaml`
- `k8s/kafka_ui.yaml`
- `k8s/producer.yaml`
- `k8s/ingress.yaml`

After the script completes, verify the ingress is available:

```bash
kubectl get ingress -n app
```

## 4. Shutdown / Tear Down the Cluster Resources

Use the repo shutdown script to delete the Kubernetes resources cleanly:

```bash
chmod +x scripts/k8s-delete.sh
./scripts/k8s-delete.sh
```

This script deletes the same resources created by the apply script.

If you also installed the NGINX ingress controller and want to remove it, delete its manifest separately:

```bash
kubectl delete -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.9.4/deploy/static/provider/cloud/deploy.yaml
```

## Notes

- Make sure Kubernetes is enabled in Docker Desktop before running these steps.
- If you change hostnames in `k8s/ingress.yaml`, update `/etc/hosts` accordingly.
- Use `kubectl get pods -n app` to confirm the app pods are running after applying the manifest.
