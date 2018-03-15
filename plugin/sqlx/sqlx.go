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
	sqlPkg   generator.Single
	dbtables map[string][]string
	messages []*generator.Descriptor
}

type validmessage struct {
	Message   *generator.Descriptor
	ccMsgName string
	database  string
	table     string
	tablevar  string
	// columns map?? // Unknown what type this should be
}

type column struct {
	VarName      string
	DBColumnName string
	isKey        bool
	isAutoGenKey bool
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
	p.sqlPkg = p.NewImport("database/sql")
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
		validMessages = append(validMessages, validmessage{Message: message, database: *dbname, table: *tablename, tablevar: *dbname + "_" + *tablename, ccMsgName: generator.CamelCase(message.GetName())})
	}
	// this generates the main constants that the messages build apon for their apis.
	p.genConstants()
	for _, v := range validMessages {
		cols := p.IdentifyColumns(&v)

		p.GenerateFuncs(v, cols)
		p.P()
	}
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
		col := column{VarName: ccTypeName, DBColumnName: *colname}

		if proto.HasExtension(field.Options, dbproto.E_Key) {
			col.isKey = true
		}

		if proto.HasExtension(field.Options, dbproto.E_Autokey) {
			col.isKey = true
			col.isAutoGenKey = true
		}

		columns = append(columns, col)
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

func (p *plugin) GenerateFuncs(v validmessage, cols []column) {
	// generate constants
	p.genMsgConstants(v, cols)
	// generate Insert Func
	p.genComment(v.ccMsgName+"Insert", "Handles Insert")
	p.genInsertFunc(v, cols)
	p.genComment(v.ccMsgName+"Insert", "Handles Insert")
	p.genMultiFunc(v.ccMsgName, "Insert")

	// generate Update Func
	p.genComment(v.ccMsgName+"Update", "Handles Update")
	p.genUpdateFunc(v, cols)
	p.genComment(v.ccMsgName+"Update", "Handles Update")
	p.genMultiFunc(v.ccMsgName, "Update")

	// generate del Func
	p.genComment(v.ccMsgName+"Delete", "Handles deleting")
	p.genDeleteFunc(v, cols)
	p.genComment(v.ccMsgName+"MultiDelete", "Handles deleting multiple")
	p.genMultiFunc(v.ccMsgName, "Delete")

	// generate Get Func
	p.genComment(v.ccMsgName+"Get", "Handles getting")
	p.genGetFunc(v, cols)

}

func (p *plugin) genConstants() {
	// search type
	p.P(`type dbSearchType int`)
	p.P()
	p.P(`const(`)
	p.In()
	p.P(`DBS_EQUAL dbSearchType = iota`)
	p.P(`DBS_GREATERTHAN`)
	p.P(`DBS_LESSTHAN`)
	// following emplictly change it from equal's to Like
	p.P(`DBS_WILDCARD_BOTH`)
	p.P(`DBS_WILDCARD_BACK`)
	p.P(`DBS_WILDCARD_FRONT`)
	// all will ignore the target field and value
	p.P(`DBS_ALL`)
	p.Out()
	p.P(`)`)
	p.P()
}

func (p *plugin) genMsgConstants(v validmessage, cols []column) {
	// Generate functions to get column names for each member.

}

func (p *plugin) genInsertFunc(v validmessage, cols []column) {
	//objNameLower := strings.ToLower(v.Message.GetName())
	p.P(`func `, v.ccMsgName, `Insert(db *`, p.sqlxPkg.Use(), `.DB, in `, v.ccMsgName, `) (`, p.sqlPkg.Use(), `.Result, error) {`)
	p.In()
	keys := ""
	valuesdeclare := ""
	values := ""
	var insertableCols []column
	for _, v := range cols {
		if v.isAutoGenKey == false {
			insertableCols = append(insertableCols, v)
		}
	}
	for k, v := range insertableCols {
		if k != 0 {
			keys += ","
			valuesdeclare += ","
			values += ","
		}
		keys += v.DBColumnName
		valuesdeclare += "?"
		values += "in." + v.VarName
	}
	p.P(`statement := "INSERT INTO `, v.database, `.`, v.table, ` (`, keys, `) VALUES (`, valuesdeclare, `)"`)
	p.P(`return db.Exec(statement,`, values, `)`)
	p.Out()
	p.P(`}`)
}

func (p *plugin) genUpdateFunc(v validmessage, cols []column) {
	// setup where field based off key's.
	// and only asign to things that are not keys
	var keys []column
	var variables []column
	for _, v := range cols {
		if v.isKey || v.isAutoGenKey {
			keys = append(keys, v)
			continue
		}
		variables = append(variables, v)
	}
	var whereSetup string
	var whereAssign string
	for k, v := range keys {
		if k != 0 {
			whereSetup += " AND "
			whereAssign += ", "
		}
		whereSetup += v.DBColumnName + "=?"
		whereAssign += "in." + v.VarName
	}

	var variableSetup string
	var variableAssign string
	for k, v := range variables {
		if k != 0 {
			variableSetup += ", "
			variableAssign += ", "
		}
		variableSetup += v.DBColumnName + "=?"
		variableAssign += "in." + v.VarName
	}

	p.P(`func `, v.ccMsgName, `Update(db *`, p.sqlxPkg.Use(), `.DB, in `, v.ccMsgName, `) (`, p.sqlPkg.Use(), `.Result, error) {`)
	p.In()
	p.P(`statement := "UPDATE `, v.database, `.`, v.table, ` SET `, variableSetup, ` WHERE `, whereSetup, `"`)
	p.P(`return db.Exec(statement,`, variableAssign, `,`, whereAssign, `)`)
	p.Out()
	p.P(`}`)
	return
}

func (p *plugin) genDeleteFunc(v validmessage, cols []column) {
	var whereSetup string
	var whereAssign string
	var keys []column
	for _, v := range cols {
		if v.isAutoGenKey || v.isKey {
			keys = append(keys, v)
		}
	}
	for k, v := range keys {
		if k != 0 {
			whereSetup += " AND "
			whereAssign += ", "
		}
		whereSetup += v.DBColumnName + "=?"
		whereAssign += "in." + v.VarName
	}

	p.P(`func `, v.ccMsgName, `Delete(db *`, p.sqlxPkg.Use(), `.DB, in `, v.ccMsgName, `) (`, p.sqlPkg.Use(), `.Result, error) {`)
	p.In()
	p.P(`statement := "DELETE FROM `, v.database, `.`, v.table, ` WHERE `, whereSetup, `"`)
	p.P(`return db.Exec(statement,`, whereAssign, `)`)
	p.Out()
	p.P(`}`)
}

func (p *plugin) genGetFunc(v validmessage, cols []column) {
	p.P(`func `, v.ccMsgName, `Get(db *`, p.sqlxPkg.Use(), `.DB,column string, searchType dbSearchType,value string) (`, v.ccMsgName, `, error) {`)
	p.In()

	p.P(`return `, v.ccMsgName, `{},nil`)
	p.Out()
	p.P(`}`)
}

func (p *plugin) genMultiFunc(msgName string, mainFuncName string) {
	p.P(`func `, msgName, mainFuncName, `Multi(db *`, p.sqlxPkg.Use(), `.DB, in []`, msgName, `) ( results []`, p.sqlPkg.Use(), `.Result, errors []error) {`)
	p.In()
	p.P(`for _ , v := range in {`)
	p.In()
	p.P(`res , err := `, msgName, mainFuncName, `(db,v)`)
	p.P(`results = append(results,res)`)
	p.P(`errors = append(errors,err)`)
	p.Out()
	p.P(`}`)
	p.P(`return`)
	p.Out()
	p.P(`}`)
}

func (p *plugin) genComment(funcName string, comment string) {
	if comment != "" && funcName != "" {
		p.P(`// `, funcName, " ", comment)
	}
}
