.PHONY: generate
generate:
	rm -rf gen/*
	docker run --rm -it -v ${PWD}:/defs namely/protoc-all:latest \
		-d api \
		-o gen \
		-l go \
		--with-gateway \
		--grpc-gateway_opt generate_unbound_methods=true

.PHONY: clone
clone:
	@if [[ "${dir}" == "" || "${url}" == "" ]]; then echo "make clone dir=/path/on/host/machine url=https://github.com/repo"; exit 1; fi
	docker run -it --rm -v ${dir}:/git alpine/git clone ${url}

.PHONT: extensions
extensions:
	make googleapis
	make grpc-gateway

.PHONY: googleapis
googleapis:
	rm -rf ${PWD}/api/googleapis
	make clone dir=${PWD}/api url=https://github.com/googleapis/googleapis.git

.PHONY: grpc-gateway
grpc-gateway:
	rm -rf ${PWD}/api/grpc-gateway
	make clone dir=${PWD}/api url=https://github.com/grpc-ecosystem/grpc-gateway.git
