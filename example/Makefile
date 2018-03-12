gen:
	(protoc -I=/usr/local/include -I=. -I=$(GOPATH)/src -I=$(GOPATH)/src/github.com/gogo/protobuf/protobuf --database_out=plugins=grpc:. test.proto && \
	 protoc -I=/usr/local/include -I=. -I=$(GOPATH)/src -I=$(GOPATH)/src/github.com/gogo/protobuf/protobuf --go_out=plugins=grpc:. test.proto)
install:
	(cd ../protoc-gen-database && go install && cd ../example)
test:
	(make install && make gen)