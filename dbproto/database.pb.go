// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: database.proto

/*
Package dbproto is a generated protocol buffer package.

It is generated from these files:
	database.proto

It has these top-level messages:
*/
package dbproto

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

var E_Sqlxdb = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.MessageOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         64100,
	Name:          "dbproto.sqlxdb",
	Tag:           "varint,64100,opt,name=sqlxdb",
	Filename:      "database.proto",
}

var E_Dbname = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.MessageOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         64101,
	Name:          "dbproto.dbname",
	Tag:           "bytes,64101,opt,name=dbname",
	Filename:      "database.proto",
}

var E_Tablename = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.MessageOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         64102,
	Name:          "dbproto.tablename",
	Tag:           "bytes,64102,opt,name=tablename",
	Filename:      "database.proto",
}

var E_Colname = &proto.ExtensionDesc{
	ExtendedType:  (*google_protobuf.FieldOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         65100,
	Name:          "dbproto.colname",
	Tag:           "bytes,65100,opt,name=colname",
	Filename:      "database.proto",
}

func init() {
	proto.RegisterExtension(E_Sqlxdb)
	proto.RegisterExtension(E_Dbname)
	proto.RegisterExtension(E_Tablename)
	proto.RegisterExtension(E_Colname)
}

func init() { proto.RegisterFile("database.proto", fileDescriptorDatabase) }

var fileDescriptorDatabase = []byte{
	// 181 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4b, 0x49, 0x2c, 0x49,
	0x4c, 0x4a, 0x2c, 0x4e, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x4f, 0x49, 0x02, 0x33,
	0xa4, 0x14, 0xd2, 0xf3, 0xf3, 0xd3, 0x73, 0x52, 0xf5, 0xc1, 0xbc, 0xa4, 0xd2, 0x34, 0xfd, 0x94,
	0xd4, 0xe2, 0xe4, 0xa2, 0xcc, 0x82, 0x92, 0xfc, 0x22, 0x88, 0x52, 0x2b, 0x4b, 0x2e, 0xb6, 0xe2,
	0xc2, 0x9c, 0x8a, 0x94, 0x24, 0x21, 0x79, 0x3d, 0x88, 0x62, 0x3d, 0x98, 0x62, 0x3d, 0xdf, 0xd4,
	0xe2, 0xe2, 0xc4, 0xf4, 0x54, 0xff, 0x82, 0x92, 0xcc, 0xfc, 0xbc, 0x62, 0x89, 0x27, 0x5f, 0x98,
	0x15, 0x18, 0x35, 0x38, 0x82, 0xa0, 0x1a, 0x40, 0x5a, 0x53, 0x92, 0xf2, 0x12, 0x73, 0x53, 0x09,
	0x6b, 0x7d, 0x0a, 0xd6, 0xca, 0x19, 0x04, 0xd5, 0x60, 0x65, 0xcf, 0xc5, 0x59, 0x92, 0x98, 0x94,
	0x93, 0x4a, 0x9c, 0xee, 0x67, 0x50, 0xdd, 0x08, 0x3d, 0x56, 0x96, 0x5c, 0xec, 0xc9, 0xf9, 0x39,
	0x60, 0xed, 0xb2, 0x18, 0xda, 0xdd, 0x32, 0x53, 0x73, 0x52, 0x60, 0x9a, 0xcf, 0xfc, 0x81, 0x68,
	0x86, 0xa9, 0x77, 0xe2, 0x8c, 0x82, 0x05, 0x0f, 0x20, 0x00, 0x00, 0xff, 0xff, 0x6d, 0x18, 0xfe,
	0xb9, 0x38, 0x01, 0x00, 0x00,
}
