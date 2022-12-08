generate:
	rm -rf gen/*
	docker run --rm -it -v ${PWD}:/defs namely/protoc-all:latest \
		-d api \
		-o gen \
		-l go \
		--with-gateway \
		--grpc-gateway_opt generate_unbound_methods=true

extensions:
	make googleapis
	make grpc-gateway

googleapis:
	rm -rf ${PWD}/api/googleapis
	docker run -it --rm -v ${PWD}/api:/git alpine/git clone https://github.com/googleapis/googleapis.git

grpc-gateway:
	rm -rf ${PWD}/api/grpc-gateway
	docker run -it --rm -v ${PWD}/api:/git alpine/git clone https://github.com/grpc-ecosystem/grpc-gateway.git
