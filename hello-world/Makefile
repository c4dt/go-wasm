GOROOT := $(shell go env GOROOT)
GOPATH := $(shell go env GOPATH)

.PHONY: build
build: main.wasm wasm_exec.js

main.wasm: main.go
	GOARCH=wasm GOOS=js go build -o $@

wasm_exec.js: $(GOROOT)/misc/wasm/wasm_exec.js
	cp $< $@

.PHONY: serve
$(GOPATH)/bin/goexec:
	go get github.com/shurcooL/goexec
serve: | $(GOPATH)/bin/goexec
serve: private PATH := $(PATH):$(GOPATH)/bin
serve: | build
	goexec 'http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`.`)))'
