.PHONY: proto
proto:
	rm -rf api/proto
	docker run --rm -w /workdir -v ${PWD}/api:/workdir openapitools/openapi-generator-cli generate \
		-i service.yml \
		-g protobuf-schema \
		-o proto
	mkdir api/proto/plugins
	make clone dir=${PWD}/api/proto/plugins url=https://github.com/googleapis/googleapis.git
	make clone dir=${PWD}/api/proto/plugins url=https://github.com/grpc-ecosystem/grpc-gateway.git

.PHONY: protoc
protoc:
	make build name=protoc
	rm -rf gen/pbgen
	mkdir gen/pbgen
	docker run --rm -it -v ${PWD}:/workdir chaseisabelle/protoc:local protoc \
		-I api/proto \
		-I api/proto/services \
		-I api/proto/models \
		-I api/proto/plugins/googleapis \
		-I api/proto/plugins/grpc-gateway \
		--go_out gen/pbgen \
		--go_opt paths=source_relative \
		--go_opt=Mpbgen \
		--go-grpc_out gen/pbgen \
		--go-grpc_opt paths=source_relative,require_unimplemented_servers=false \
		--grpc-gateway_out gen/pbgen \
		--grpc-gateway_opt paths=source_relative,generate_unbound_methods=true \
		api/proto/services/*.proto

.PHONY: clone
clone:
	@docker run -it --rm -v ${dir}:/git alpine/git clone ${url}

.PHONY: image
image:
	@echo "chaseisabelle/${name}:local"

.PHONY: built
built:
	@docker image inspect $(shell make image name="${name}") >/dev/null 2>&1

.PHONY: build
build:
	@make built image="${name}" || docker build --no-cache --target "${name}" -t $(shell make image name=${name}) "images/${name}"

.PHONY: rebuild
rebuild:
	@docker rmi $(shell make image name=${name})
	@make build name="${name}"