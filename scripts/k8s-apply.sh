#!/bin/bash
set -e

echo "Applying Kubernetes manifests to Docker Desktop..."
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/kafka.yaml
kubectl apply -f k8s/schema_registry.yaml
kubectl apply -f k8s/kafka_connect.yaml
kubectl apply -f k8s/kafka_ui.yaml
kubectl apply -f k8s/producer.yaml
echo "✓ All Kubernetes objects applied successfully"
