appname=$(shell basename $(shell pwd)) #note: only works if make run from project root
build_root="_build"
binary=$(build_root)/$(appname)

.PHONY=all,compile,run

all: compile run

compile:
	@mkdir -p $(build_root)
	@go build -o $(binary)

run: compile
	@./$(binary)

fmt:
	@go fmt

test:
	@go test

bench:
	@go test -bench=.
