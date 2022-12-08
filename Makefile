generate:
	rm -rf gen/*
	docker run --rm -it -v ${PWD}:/defs namely/protoc-all:latest \
		-d api \
		-o gen \
		-l go \
		--with-gateway \
		--grpc-gateway_opt generate_unbound_methods=true

googleapis:
	rm ${PWD}/api/googleapis
	docker run -it --rm -v ${PWD}/api:/git alpine/git clone https://github.com/googleapis/googleapis.git

grpc-gateway:
	rm ${PWD}/api/grpc-gateway
	docker run -it --rm -v ${PWD}/api:/git alpine/git clone https://github.com/grpc-ecosystem/grpc-gateway.git