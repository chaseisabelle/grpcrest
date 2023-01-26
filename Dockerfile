FROM golang AS protoc

WORKDIR /workdir

RUN apt-get update -y
RUN apt-get install -y protobuf-compiler

RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

CMD ["protoc", "--help"]


FROM golang AS sqlboiler

RUN go install github.com/volatiletech/sqlboiler/v4@latest
RUN go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest
RUN go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@latest
RUN go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mssql@latest
RUN go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-sqlite3@latest

CMD ["sqlboiler", "--help"]