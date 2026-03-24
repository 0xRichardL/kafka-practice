.PHONY: gen-schemas-go init build-images k8s-apply k8s-delete k8s-deploy

init:
	npm install -g @mermaid-js/mermaid-cli

gen-schemas-go:
	protoc \
		--go_out=paths=source_relative:. \
		schemas/*.proto

# Docker build commands
build-images:
	@echo "Building local Docker images..."
	docker build -t kafka-practice/producer:v1 ./producer

# Kubernetes commands for Docker Desktop
k8s-apply:
	@bash scripts/k8s-apply.sh

k8s-delete:
	@bash scripts/k8s-delete.sh

k8s-deploy: build-images k8s-apply
	@echo "✓ Deployment complete! Using Docker Desktop UI to monitor."
