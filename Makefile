.DEFAULT_GOAL := gen

.PHONY: gen

gen: fmt
	go generate ./...
	sqlc generate
	easyp generate

fmt:
	goimports -w .
