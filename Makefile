.PHOHY: import build deploy-gcp

GOROOT:=$(shell go env GOROOT)

import:
	cp "$(GOROOT)/misc/wasm/wasm_exec.js" static/

build:
	GOOS=js GOARCH=wasm go build -o static/age.wasm

deploy-gcp:
	gcloud app deploy