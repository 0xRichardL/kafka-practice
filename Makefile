.PHONY: gen-schemas-go

init:
	npm install -g @mermaid-js/mermaid-cli

gen-schemas-go:
	protoc \
		--go_out=paths=source_relative:. \
		schemas/*.proto
