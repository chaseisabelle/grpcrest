.PHONY: plugins
plugins:
	rm -rf api/googleapis api/grpc-gateway
	make clone url=https://github.com/googleapis/googleapis.git dir=api/googleapis
	make clone url=https://github.com/grpc-ecosystem/grpc-gateway.git dir=api/grpc-gateway

.PHONY: protoc
protoc:
	rm -rf gen/pb api/service.swagger.yaml
	mkdir gen/pb
	make run what=protoc

.PHONY: sqlboiler
sqlboiler:
	make run what=sqlboiler

.PHONY: test
test:
	make run what=tester

.PHONY: cover
cover:
	make run what=coverer

.PHONY: vet
vet:
	make run what=vetter

.PHONY: clone
clone:
	make run what=git command="clone ${url} ${dir}"

.PHONY: run
run:
	docker compose run --rm -it ${what} ${command}