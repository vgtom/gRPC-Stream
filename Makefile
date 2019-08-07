
all: client server

protoc:
	@echo "Generating Go files"
	cd src/proto && protoc --go_out=plugins=grpc:. *.proto

server: protoc
	@echo "Building server"
	go build -o server \
		soln/src/server

client: protoc
	@echo "Building client"
	go build -o client \
		soln/src/client

clean:
	go clean
	rm -f server client

.PHONY: client server protoc