package sqlx

import (
	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

type plugin struct {
	*generator.Generator
	generator.PluginImports
	sqlxPkg generator.Single
}

func init() {
	//generator.RegisterPlugin(NewPlugin())
}

func NewPlugin() *plugin {
	return &plugin{}
}

func (p *plugin) Name() string {
	return "database"
}

func (p *plugin) Init(g *generator.Generator) {
	p.Generator = g
}

func (p *plugin) Generate(file *generator.FileDescriptor) {
	p.PluginImports = generator.NewPluginImports(p.Generator)
	//protopkg := p.NewImport("github.com/gogo/protobuf/proto")
	// support for if this plugin is used externally with standard proto and not gogo
	if !gogoproto.ImportsGoGoProto(file.FileDescriptorProto) {
		//protopkg = p.NewImport("github.com/golang/protobuf/proto")
	}
	p.sqlxPkg = p.NewImport("github.com/jmoiron/sqlx")
	for _, message := range file.Messages() {
		p.P(`func `, message.Name, `() {fmt.Println(`, p.sqlxPkg.Use(), `.NAMED)}`)
		//message.DescriptorProto

	}
}
