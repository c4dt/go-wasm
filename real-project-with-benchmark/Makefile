GOROOT := $(shell go env GOROOT)
GOPATH := $(shell go env GOPATH)

.PHONY: build
build: dist/main.wasm dist/main.js dist/index.html dist/wasm_exec.js

node_modules/.bin/webpack:
	npm i

dist:
	mkdir $@
dist/main.wasm: src/main.go | dist
	GOARCH=wasm GOOS=js go build -o $@ ./$(<D)
dist/main.js: tsconfig.json $(wildcard src/*.ts) | dist node_modules/.bin/webpack
	npx webpack
dist/index.html: src/index.html | dist
	cp $< $@
dist/wasm_exec.js: $(GOROOT)/misc/wasm/wasm_exec.js | dist
	cp $< $@

.PHONY: serve
$(GOPATH)/bin/goexec:
	go get github.com/shurcooL/goexec
serve: | $(GOPATH)/bin/goexec
serve: private PATH := $(PATH):$(GOPATH)/bin
serve: | build
	goexec 'http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`dist`)))'
