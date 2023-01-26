.PHONY: plugins
plugins:
	make clone dir=${PWD}/api url=https://github.com/googleapis/googleapis.git
	make clone dir=${PWD}/api url=https://github.com/grpc-ecosystem/grpc-gateway.git

.PHONY: protoc
protoc:
	make build name=protoc
	rm -rf gen/pb api/service.swagger.yaml
	mkdir gen/pb
	docker run --rm -it -v ${PWD}:/workdir $(shell make image name=protoc) protoc \
		-I api \
		-I api/googleapis \
		-I api/grpc-gateway \
		--go_out gen/pb \
		--go_opt paths=source_relative \
		--go-grpc_out gen/pb \
		--go-grpc_opt paths=source_relative,require_unimplemented_servers=false \
		--grpc-gateway_out gen/pb \
		--grpc-gateway_opt paths=source_relative,generate_unbound_methods=true \
		--openapiv2_out api \
		--openapiv2_opt logtostderr=true,use_go_templates=true,output_format=yaml \
		api/service.proto

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