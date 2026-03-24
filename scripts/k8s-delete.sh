#!/bin/bash
set -e

echo "Deleting Kubernetes resources from Docker Desktop..."
kubectl delete -f k8s/producer.yaml --ignore-not-found
kubectl delete -f k8s/kafka_ui.yaml --ignore-not-found
kubectl delete -f k8s/kafka_connect.yaml --ignore-not-found
kubectl delete -f k8s/schema_registry.yaml --ignore-not-found
kubectl delete -f k8s/kafka.yaml --ignore-not-found
kubectl delete -f k8s/namespace.yaml --ignore-not-found
echo "✓ All Kubernetes resources deleted"
