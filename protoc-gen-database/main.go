package main

import (
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/vanity"
	"github.com/gogo/protobuf/vanity/command"
	"github.com/zanven42/proto-database-go/plugin/sqlx"
)

func main() {
	req := command.Read()
	files := req.GetProtoFile()
	files = vanity.FilterFiles(files, vanity.NotGoogleProtobufDescriptorProto)
	for _, opt := range []func(*descriptor.FileDescriptorProto){
		vanity.TurnOffGoGettersAll,
		vanity.TurnOffGoStringerAll,
		//vanity.TurnOnMarshalerAll,
		//vanity.TurnOnStringerAll,
		//vanity.TurnOnUnmarshalerAll,
		//vanity.TurnOnSizerAll,
	} {
		vanity.ForEachFile(files, opt)
	}

	resp := command.GeneratePlugin(req, sqlx.NewPlugin(), ".database.go")
	command.Write(resp)
}
