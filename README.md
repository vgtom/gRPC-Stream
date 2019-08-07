# gRPC-Stream
a gRPC based streaming client and server

# Dependencies
The protoc compiler must be installed on your system
Then download the project from Github

# Build
Juke run
make all
and the project will be built

# Run
go run the client and server. The server first

# Brief Architectere
The client is a streaming gRPC client that generates Random numbers and then signs it with RSA Public key and sends them over the gRPC stream to the server.
The server similarly receives the Signed number and verifies it and if successfully verified then calculates the MAX so far in the stream. If a new Max is found, it sends it to the server.
