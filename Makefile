generate:
	docker run --rm -it -v ${PWD}:/defs namely/protoc-all:latest -d api -o gen -l go --with-gateway

regenerate:
	rm -rf gen/*
	make generate

googleapis:
	docker run -it --rm -v ${PWD}/api:/git alpine/git clone https://github.com/googleapis/googleapis.git