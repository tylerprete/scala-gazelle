// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.14.0
// source: build/stack/gazelle/scala/parse/parser.proto

package parse

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ParseRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Filenames     []string `protobuf:"bytes,1,rep,name=filenames,proto3" json:"filenames,omitempty"`
	WantParseTree bool     `protobuf:"varint,2,opt,name=want_parse_tree,json=wantParseTree,proto3" json:"want_parse_tree,omitempty"`
}

func (x *ParseRequest) Reset() {
	*x = ParseRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_build_stack_gazelle_scala_parse_parser_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ParseRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ParseRequest) ProtoMessage() {}

func (x *ParseRequest) ProtoReflect() protoreflect.Message {
	mi := &file_build_stack_gazelle_scala_parse_parser_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ParseRequest.ProtoReflect.Descriptor instead.
func (*ParseRequest) Descriptor() ([]byte, []int) {
	return file_build_stack_gazelle_scala_parse_parser_proto_rawDescGZIP(), []int{0}
}

func (x *ParseRequest) GetFilenames() []string {
	if x != nil {
		return x.Filenames
	}
	return nil
}

func (x *ParseRequest) GetWantParseTree() bool {
	if x != nil {
		return x.WantParseTree
	}
	return false
}

type ParseResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Files         []*File `protobuf:"bytes,1,rep,name=files,proto3" json:"files,omitempty"`
	Error         string  `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	ElapsedMillis int64   `protobuf:"varint,3,opt,name=elapsed_millis,json=elapsedMillis,proto3" json:"elapsed_millis,omitempty"`
}

func (x *ParseResponse) Reset() {
	*x = ParseResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_build_stack_gazelle_scala_parse_parser_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ParseResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ParseResponse) ProtoMessage() {}

func (x *ParseResponse) ProtoReflect() protoreflect.Message {
	mi := &file_build_stack_gazelle_scala_parse_parser_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ParseResponse.ProtoReflect.Descriptor instead.
func (*ParseResponse) Descriptor() ([]byte, []int) {
	return file_build_stack_gazelle_scala_parse_parser_proto_rawDescGZIP(), []int{1}
}

func (x *ParseResponse) GetFiles() []*File {
	if x != nil {
		return x.Files
	}
	return nil
}

func (x *ParseResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

func (x *ParseResponse) GetElapsedMillis() int64 {
	if x != nil {
		return x.ElapsedMillis
	}
	return 0
}

var File_build_stack_gazelle_scala_parse_parser_proto protoreflect.FileDescriptor

var file_build_stack_gazelle_scala_parse_parser_proto_rawDesc = []byte{
	0x0a, 0x2c, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2f, 0x73, 0x74, 0x61, 0x63, 0x6b, 0x2f, 0x67, 0x61,
	0x7a, 0x65, 0x6c, 0x6c, 0x65, 0x2f, 0x73, 0x63, 0x61, 0x6c, 0x61, 0x2f, 0x70, 0x61, 0x72, 0x73,
	0x65, 0x2f, 0x70, 0x61, 0x72, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1f,
	0x62, 0x75, 0x69, 0x6c, 0x64, 0x2e, 0x73, 0x74, 0x61, 0x63, 0x6b, 0x2e, 0x67, 0x61, 0x7a, 0x65,
	0x6c, 0x6c, 0x65, 0x2e, 0x73, 0x63, 0x61, 0x6c, 0x61, 0x2e, 0x70, 0x61, 0x72, 0x73, 0x65, 0x1a,
	0x2a, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2f, 0x73, 0x74, 0x61, 0x63, 0x6b, 0x2f, 0x67, 0x61, 0x7a,
	0x65, 0x6c, 0x6c, 0x65, 0x2f, 0x73, 0x63, 0x61, 0x6c, 0x61, 0x2f, 0x70, 0x61, 0x72, 0x73, 0x65,
	0x2f, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x54, 0x0a, 0x0c, 0x50,
	0x61, 0x72, 0x73, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x66,
	0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09,
	0x66, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x73, 0x12, 0x26, 0x0a, 0x0f, 0x77, 0x61, 0x6e,
	0x74, 0x5f, 0x70, 0x61, 0x72, 0x73, 0x65, 0x5f, 0x74, 0x72, 0x65, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x0d, 0x77, 0x61, 0x6e, 0x74, 0x50, 0x61, 0x72, 0x73, 0x65, 0x54, 0x72, 0x65,
	0x65, 0x22, 0x89, 0x01, 0x0a, 0x0d, 0x50, 0x61, 0x72, 0x73, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x3b, 0x0a, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x25, 0x2e, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2e, 0x73, 0x74, 0x61, 0x63, 0x6b,
	0x2e, 0x67, 0x61, 0x7a, 0x65, 0x6c, 0x6c, 0x65, 0x2e, 0x73, 0x63, 0x61, 0x6c, 0x61, 0x2e, 0x70,
	0x61, 0x72, 0x73, 0x65, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x05, 0x66, 0x69, 0x6c, 0x65, 0x73,
	0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x25, 0x0a, 0x0e, 0x65, 0x6c, 0x61, 0x70, 0x73, 0x65,
	0x64, 0x5f, 0x6d, 0x69, 0x6c, 0x6c, 0x69, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d,
	0x65, 0x6c, 0x61, 0x70, 0x73, 0x65, 0x64, 0x4d, 0x69, 0x6c, 0x6c, 0x69, 0x73, 0x32, 0x72, 0x0a,
	0x06, 0x50, 0x61, 0x72, 0x73, 0x65, 0x72, 0x12, 0x68, 0x0a, 0x05, 0x50, 0x61, 0x72, 0x73, 0x65,
	0x12, 0x2d, 0x2e, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2e, 0x73, 0x74, 0x61, 0x63, 0x6b, 0x2e, 0x67,
	0x61, 0x7a, 0x65, 0x6c, 0x6c, 0x65, 0x2e, 0x73, 0x63, 0x61, 0x6c, 0x61, 0x2e, 0x70, 0x61, 0x72,
	0x73, 0x65, 0x2e, 0x50, 0x61, 0x72, 0x73, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x2e, 0x2e, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2e, 0x73, 0x74, 0x61, 0x63, 0x6b, 0x2e, 0x67, 0x61,
	0x7a, 0x65, 0x6c, 0x6c, 0x65, 0x2e, 0x73, 0x63, 0x61, 0x6c, 0x61, 0x2e, 0x70, 0x61, 0x72, 0x73,
	0x65, 0x2e, 0x50, 0x61, 0x72, 0x73, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x42, 0x6a, 0x0a, 0x1f, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2e, 0x73, 0x74, 0x61, 0x63, 0x6b,
	0x2e, 0x67, 0x61, 0x7a, 0x65, 0x6c, 0x6c, 0x65, 0x2e, 0x73, 0x63, 0x61, 0x6c, 0x61, 0x2e, 0x70,
	0x61, 0x72, 0x73, 0x65, 0x50, 0x01, 0x5a, 0x45, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x73, 0x74, 0x61, 0x63, 0x6b, 0x62, 0x2f, 0x73, 0x63, 0x61, 0x6c, 0x61, 0x2d,
	0x67, 0x61, 0x7a, 0x65, 0x6c, 0x6c, 0x65, 0x2f, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2f, 0x73, 0x74,
	0x61, 0x63, 0x6b, 0x2f, 0x67, 0x61, 0x7a, 0x65, 0x6c, 0x6c, 0x65, 0x2f, 0x73, 0x63, 0x61, 0x6c,
	0x61, 0x2f, 0x70, 0x61, 0x72, 0x73, 0x65, 0x3b, 0x70, 0x61, 0x72, 0x73, 0x65, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_build_stack_gazelle_scala_parse_parser_proto_rawDescOnce sync.Once
	file_build_stack_gazelle_scala_parse_parser_proto_rawDescData = file_build_stack_gazelle_scala_parse_parser_proto_rawDesc
)

func file_build_stack_gazelle_scala_parse_parser_proto_rawDescGZIP() []byte {
	file_build_stack_gazelle_scala_parse_parser_proto_rawDescOnce.Do(func() {
		file_build_stack_gazelle_scala_parse_parser_proto_rawDescData = protoimpl.X.CompressGZIP(file_build_stack_gazelle_scala_parse_parser_proto_rawDescData)
	})
	return file_build_stack_gazelle_scala_parse_parser_proto_rawDescData
}

var file_build_stack_gazelle_scala_parse_parser_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_build_stack_gazelle_scala_parse_parser_proto_goTypes = []interface{}{
	(*ParseRequest)(nil),  // 0: build.stack.gazelle.scala.parse.ParseRequest
	(*ParseResponse)(nil), // 1: build.stack.gazelle.scala.parse.ParseResponse
	(*File)(nil),          // 2: build.stack.gazelle.scala.parse.File
}
var file_build_stack_gazelle_scala_parse_parser_proto_depIdxs = []int32{
	2, // 0: build.stack.gazelle.scala.parse.ParseResponse.files:type_name -> build.stack.gazelle.scala.parse.File
	0, // 1: build.stack.gazelle.scala.parse.Parser.Parse:input_type -> build.stack.gazelle.scala.parse.ParseRequest
	1, // 2: build.stack.gazelle.scala.parse.Parser.Parse:output_type -> build.stack.gazelle.scala.parse.ParseResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_build_stack_gazelle_scala_parse_parser_proto_init() }
func file_build_stack_gazelle_scala_parse_parser_proto_init() {
	if File_build_stack_gazelle_scala_parse_parser_proto != nil {
		return
	}
	file_build_stack_gazelle_scala_parse_file_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_build_stack_gazelle_scala_parse_parser_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ParseRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_build_stack_gazelle_scala_parse_parser_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ParseResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_build_stack_gazelle_scala_parse_parser_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_build_stack_gazelle_scala_parse_parser_proto_goTypes,
		DependencyIndexes: file_build_stack_gazelle_scala_parse_parser_proto_depIdxs,
		MessageInfos:      file_build_stack_gazelle_scala_parse_parser_proto_msgTypes,
	}.Build()
	File_build_stack_gazelle_scala_parse_parser_proto = out.File
	file_build_stack_gazelle_scala_parse_parser_proto_rawDesc = nil
	file_build_stack_gazelle_scala_parse_parser_proto_goTypes = nil
	file_build_stack_gazelle_scala_parse_parser_proto_depIdxs = nil
}
