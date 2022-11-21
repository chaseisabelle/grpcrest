generate:
	docker run --rm -it -v ${PWD}:/defs namely/protoc-all:latest -d api -o gen -l go --with-gateway

ungenerate:
	rm -rf gen/*

regenerate:
	make ungenerate
	make generate

googleapis:
	cd api && git clone https://github.com/googleapis/googleapis.git