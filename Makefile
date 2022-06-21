BINARY_NAME=stockup
GOROOT=$(shell go env GOROOT)

SITE_ASSETS=site/assets
JS_ASSETS=$(SITE_ASSETS)/js
WASM_ASSETS=$(SITE_ASSETS)/il

dir:
	mkdir -p $(SITE_ASSETS)
	mkdir -p $(JS_ASSETS)
	mkdir -p $(WASM_ASSETS)

runtime: $(GOROOT)/misc/wasm/wasm_exec.js
	cp $< $(JS_ASSETS)/runtime.js

clean:
	go clean
	rm -rf $(WASM_ASSETS)/*.wasm $(JS_ASSETS)/*.js

build: dir clean runtime
	GOOS=js GOARCH=wasm go build -o $(WASM_ASSETS)/$(BINARY_NAME).wasm wasm/main.go
