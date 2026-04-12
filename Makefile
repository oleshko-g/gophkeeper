.DEFAULT_GOAL := gen

.PHONY: gen

gen:
	go generate ./...
