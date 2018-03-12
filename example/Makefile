build:
	(protoc -I=/usr/local/include -I=. -I=$(GOPATH)/src -I=$(GOPATH)/src/github.com/gogo/protobuf/protobuf --database_out=plugins=grpc:. test.proto)