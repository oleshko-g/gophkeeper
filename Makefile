.DEFAULT_GOAL := gen

.PHONY: gen

gen: fmt
	easyp generate
	go generate ./...
	sqlc generate

fmt:
	goimports -w .
