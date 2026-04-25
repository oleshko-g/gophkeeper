.DEFAULT_GOAL := gen

.PHONY: gen

gen: fmt
	go generate ./...
	sqlc generate

fmt:
	goimports -w .
