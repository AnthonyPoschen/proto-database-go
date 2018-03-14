package sqlx

import (
	"fmt"
	"strings"

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
	messages []*generator.Descriptor
}

type validmessage struct {
	Message  *generator.Descriptor
	database string
	table    string
	tablevar string
	// columns map?? // Unknown what type this should be
}

type column struct {
	VarName      string
	DBColumnName string
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
	p.dbtables = map[string][]string{}
	p.messages = file.Messages()
	// maps lists of tables to databases for later processing
	// do a pass purely to grab all databases and tables
	var validMessages []validmessage
	for _, message := range p.messages {
		if message.Options == nil {
			continue
		}
		if proto.GetBoolExtension(message.Options, dbproto.E_Sqlxdb, false) == false {
			continue
		}
		dbnameraw, err := proto.GetExtension(message.Options, dbproto.E_Dbname)
		if err != nil {
			continue
		}
		dbname, ok := dbnameraw.(*string)
		if !ok {
			continue
		}
		tablenameraw, err := proto.GetExtension(message.Options, dbproto.E_Tablename)
		if err != nil {
			continue
		}
		tablename, ok := tablenameraw.(*string)
		if !ok {
			continue
		}
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
		cols := p.IdentifyColumns(&v)
		for _, col := range cols {
			p.P("// var: ", col.VarName, " - DBName: ", col.DBColumnName)
		}
		p.GenerateFuncs(v)
		p.P()
	}
}

func (p *plugin) GenerateFuncs(v validmessage) {
	// generate New Func

}

func (p *plugin) IdentifyColumns(v *validmessage) []column {
	msg := v.Message
	var columns []column
	for _, field := range msg.Field {
		// if this is not a basic type lets recursion in
		ccTypeName := generator.CamelCase(field.GetName())

		if field.IsMessage() {
			p.P()
			typeName := field.GetTypeName()
			typeNamesplit := strings.Split(typeName, ".")
			if len(typeNamesplit) > 0 {
				typeName = typeNamesplit[len(typeNamesplit)-1]
			}
			msg, err := p.GetMessageByName(typeName)
			if err != nil {
				continue
			}
			cols := p.IdentifyColumns(&validmessage{Message: msg})
			// append the object name to the result
			for k, col := range cols {

				col.VarName = ccTypeName + "." + col.VarName
				cols[k] = col
			}
			columns = append(columns, cols...)
			continue
		}
		if proto.HasExtension(field.Options, dbproto.E_Colname) == false {
			continue
		}
		colnameFace, err := proto.GetExtension(field.Options, dbproto.E_Colname)
		if err != nil {
			continue
		}
		colname, ok := colnameFace.(*string)
		if !ok {
			continue
		}
		if colname == nil {
			continue
		}
		columns = append(columns, column{VarName: ccTypeName, DBColumnName: *colname})
	}
	return columns
}

func (p *plugin) GetMessageByName(name string) (*generator.Descriptor, error) {
	for _, msg := range p.messages {
		if msg.GetName() == name {
			return msg, nil
		}
	}
	return nil, fmt.Errorf("Message not found")
}

//func genNewFunc(dbvar string, tablevar string, objType string)
