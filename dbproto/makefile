setup:
	(cd $(GOPATH)/src/github.com/gogo/protobuf/protobuf && make regenerate)
regenerate:
	protoc -I=/usr/local/include -I=. -I=$(GOPATH)/src -I=$(GOPATH)/src/github.com/gogo/protobuf/protobuf --gogo_out=Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor:. database.proto