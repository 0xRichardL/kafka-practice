.PHONY: gen-schemas-go

gen-schemas-go:
	protoc \
		--go_out=paths=source_relative:. \
		schemas/*.proto
