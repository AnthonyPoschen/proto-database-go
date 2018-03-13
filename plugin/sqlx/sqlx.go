package sqlx

import (
	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/zanven42/proto-database-go/dbproto"
)

type plugin struct {
	*generator.Generator
	generator.PluginImports
	sqlxPkg  generator.Single
	dbtables map[string][]string
}

type validmessage struct {
	Message  *generator.Descriptor
	database string
	table    string
	tablevar string
	// columns map?? // Unknown what type this should be
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

	// maps lists of tables to databases for later processing
	// do a pass purely to grab all databases and tables
	var validMessages []validmessage
	for _, message := range file.Messages() {
		if message.Options == nil {
			continue
		}
		if proto.GetBoolExtension(message.Options, dbproto.E_Sqlxdb, false) == false {
			continue
		}
		p.P("// sqlx")
		dbnameraw, err := proto.GetExtension(message.Options, dbproto.E_Dbname)
		if err != nil {
			continue
		}
		p.P("// dbName1")
		dbname, ok := dbnameraw.(*string)
		if !ok {
			continue
		}
		p.P("// dbName2")
		tablenameraw, err := proto.GetExtension(message.Options, dbproto.E_Tablename)
		if err != nil {
			continue
		}
		tablename, ok := tablenameraw.(*string)
		if !ok {
			continue
		}
		p.P("// ", *dbname, " ", *tablename)
		p.dbtables[*dbname] = append(p.dbtables[*dbname], *tablename)
		validMessages = append(validMessages, validmessage{Message: message, database: *dbname, table: *tablename, tablevar: *dbname + "_" + *tablename})
	}

	p.P("// Database and Table definitions")
	// for every db and table lets create the variables
	for db, tables := range p.dbtables {
		p.P(`var `, db, ` string = "`, db, `"`)
		for _, table := range tables {
			tableName := db + "_" + table
			p.P(`var `, tableName, ` string = "`, table, `"`)
		}
	}

	for _, v := range validMessages {
		p.GenerateFuncs(v)
		p.P()
	}
}

func (p *plugin) GenerateFuncs(v validmessage) {
	// generate New Func
	p.P()
}

//func genNewFunc(dbvar string, tablevar string, objType string)
