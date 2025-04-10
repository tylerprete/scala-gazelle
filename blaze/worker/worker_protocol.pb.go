// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.14.0
// source: blaze/worker/worker_protocol.proto

package worker_protocol

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

type Input struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Path   string `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	Digest []byte `protobuf:"bytes,2,opt,name=digest,proto3" json:"digest,omitempty"`
}

func (x *Input) Reset() {
	*x = Input{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blaze_worker_worker_protocol_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Input) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Input) ProtoMessage() {}

func (x *Input) ProtoReflect() protoreflect.Message {
	mi := &file_blaze_worker_worker_protocol_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Input.ProtoReflect.Descriptor instead.
func (*Input) Descriptor() ([]byte, []int) {
	return file_blaze_worker_worker_protocol_proto_rawDescGZIP(), []int{0}
}

func (x *Input) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *Input) GetDigest() []byte {
	if x != nil {
		return x.Digest
	}
	return nil
}

type WorkRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Arguments  []string `protobuf:"bytes,1,rep,name=arguments,proto3" json:"arguments,omitempty"`
	Inputs     []*Input `protobuf:"bytes,2,rep,name=inputs,proto3" json:"inputs,omitempty"`
	RequestId  int32    `protobuf:"varint,3,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	Cancel     bool     `protobuf:"varint,4,opt,name=cancel,proto3" json:"cancel,omitempty"`
	Verbosity  int32    `protobuf:"varint,5,opt,name=verbosity,proto3" json:"verbosity,omitempty"`
	SandboxDir string   `protobuf:"bytes,6,opt,name=sandbox_dir,json=sandboxDir,proto3" json:"sandbox_dir,omitempty"`
}

func (x *WorkRequest) Reset() {
	*x = WorkRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blaze_worker_worker_protocol_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorkRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorkRequest) ProtoMessage() {}

func (x *WorkRequest) ProtoReflect() protoreflect.Message {
	mi := &file_blaze_worker_worker_protocol_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorkRequest.ProtoReflect.Descriptor instead.
func (*WorkRequest) Descriptor() ([]byte, []int) {
	return file_blaze_worker_worker_protocol_proto_rawDescGZIP(), []int{1}
}

func (x *WorkRequest) GetArguments() []string {
	if x != nil {
		return x.Arguments
	}
	return nil
}

func (x *WorkRequest) GetInputs() []*Input {
	if x != nil {
		return x.Inputs
	}
	return nil
}

func (x *WorkRequest) GetRequestId() int32 {
	if x != nil {
		return x.RequestId
	}
	return 0
}

func (x *WorkRequest) GetCancel() bool {
	if x != nil {
		return x.Cancel
	}
	return false
}

func (x *WorkRequest) GetVerbosity() int32 {
	if x != nil {
		return x.Verbosity
	}
	return 0
}

func (x *WorkRequest) GetSandboxDir() string {
	if x != nil {
		return x.SandboxDir
	}
	return ""
}

type WorkResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ExitCode     int32  `protobuf:"varint,1,opt,name=exit_code,json=exitCode,proto3" json:"exit_code,omitempty"`
	Output       string `protobuf:"bytes,2,opt,name=output,proto3" json:"output,omitempty"`
	RequestId    int32  `protobuf:"varint,3,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	WasCancelled bool   `protobuf:"varint,4,opt,name=was_cancelled,json=wasCancelled,proto3" json:"was_cancelled,omitempty"`
}

func (x *WorkResponse) Reset() {
	*x = WorkResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_blaze_worker_worker_protocol_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorkResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorkResponse) ProtoMessage() {}

func (x *WorkResponse) ProtoReflect() protoreflect.Message {
	mi := &file_blaze_worker_worker_protocol_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorkResponse.ProtoReflect.Descriptor instead.
func (*WorkResponse) Descriptor() ([]byte, []int) {
	return file_blaze_worker_worker_protocol_proto_rawDescGZIP(), []int{2}
}

func (x *WorkResponse) GetExitCode() int32 {
	if x != nil {
		return x.ExitCode
	}
	return 0
}

func (x *WorkResponse) GetOutput() string {
	if x != nil {
		return x.Output
	}
	return ""
}

func (x *WorkResponse) GetRequestId() int32 {
	if x != nil {
		return x.RequestId
	}
	return 0
}

func (x *WorkResponse) GetWasCancelled() bool {
	if x != nil {
		return x.WasCancelled
	}
	return false
}

var File_blaze_worker_worker_protocol_proto protoreflect.FileDescriptor

var file_blaze_worker_worker_protocol_proto_rawDesc = []byte{
	0x0a, 0x22, 0x62, 0x6c, 0x61, 0x7a, 0x65, 0x2f, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2f, 0x77,
	0x6f, 0x72, 0x6b, 0x65, 0x72, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x62, 0x6c, 0x61, 0x7a, 0x65, 0x2e, 0x77, 0x6f, 0x72, 0x6b,
	0x65, 0x72, 0x22, 0x33, 0x0a, 0x05, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70,
	0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12,
	0x16, 0x0a, 0x06, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x06, 0x64, 0x69, 0x67, 0x65, 0x73, 0x74, 0x22, 0xce, 0x01, 0x0a, 0x0b, 0x57, 0x6f, 0x72, 0x6b,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x72, 0x67, 0x75, 0x6d,
	0x65, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x61, 0x72, 0x67, 0x75,
	0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x2b, 0x0a, 0x06, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x62, 0x6c, 0x61, 0x7a, 0x65, 0x2e, 0x77, 0x6f,
	0x72, 0x6b, 0x65, 0x72, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x52, 0x06, 0x69, 0x6e, 0x70, 0x75,
	0x74, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49,
	0x64, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x06, 0x63, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x12, 0x1c, 0x0a, 0x09, 0x76, 0x65, 0x72,
	0x62, 0x6f, 0x73, 0x69, 0x74, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x76, 0x65,
	0x72, 0x62, 0x6f, 0x73, 0x69, 0x74, 0x79, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x61, 0x6e, 0x64, 0x62,
	0x6f, 0x78, 0x5f, 0x64, 0x69, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x61,
	0x6e, 0x64, 0x62, 0x6f, 0x78, 0x44, 0x69, 0x72, 0x22, 0x87, 0x01, 0x0a, 0x0c, 0x57, 0x6f, 0x72,
	0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x65, 0x78, 0x69,
	0x74, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x65, 0x78,
	0x69, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x12, 0x1d,
	0x0a, 0x0a, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x64, 0x12, 0x23, 0x0a,
	0x0d, 0x77, 0x61, 0x73, 0x5f, 0x63, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x6c, 0x65, 0x64, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x77, 0x61, 0x73, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x6c,
	0x65, 0x64, 0x42, 0x74, 0x0a, 0x24, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x64, 0x65, 0x76, 0x74, 0x6f, 0x6f, 0x6c, 0x73, 0x2e, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2e,
	0x6c, 0x69, 0x62, 0x2e, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x5a, 0x4c, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x74, 0x61, 0x63, 0x6b, 0x62, 0x2f, 0x73, 0x63,
	0x61, 0x6c, 0x61, 0x2d, 0x67, 0x61, 0x7a, 0x65, 0x6c, 0x6c, 0x65, 0x2f, 0x62, 0x6c, 0x61, 0x7a,
	0x65, 0x2f, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x2f, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x5f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x3b, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x5f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_blaze_worker_worker_protocol_proto_rawDescOnce sync.Once
	file_blaze_worker_worker_protocol_proto_rawDescData = file_blaze_worker_worker_protocol_proto_rawDesc
)

func file_blaze_worker_worker_protocol_proto_rawDescGZIP() []byte {
	file_blaze_worker_worker_protocol_proto_rawDescOnce.Do(func() {
		file_blaze_worker_worker_protocol_proto_rawDescData = protoimpl.X.CompressGZIP(file_blaze_worker_worker_protocol_proto_rawDescData)
	})
	return file_blaze_worker_worker_protocol_proto_rawDescData
}

var file_blaze_worker_worker_protocol_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_blaze_worker_worker_protocol_proto_goTypes = []interface{}{
	(*Input)(nil),        // 0: blaze.worker.Input
	(*WorkRequest)(nil),  // 1: blaze.worker.WorkRequest
	(*WorkResponse)(nil), // 2: blaze.worker.WorkResponse
}
var file_blaze_worker_worker_protocol_proto_depIdxs = []int32{
	0, // 0: blaze.worker.WorkRequest.inputs:type_name -> blaze.worker.Input
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_blaze_worker_worker_protocol_proto_init() }
func file_blaze_worker_worker_protocol_proto_init() {
	if File_blaze_worker_worker_protocol_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_blaze_worker_worker_protocol_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Input); i {
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
		file_blaze_worker_worker_protocol_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WorkRequest); i {
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
		file_blaze_worker_worker_protocol_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WorkResponse); i {
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
			RawDescriptor: file_blaze_worker_worker_protocol_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_blaze_worker_worker_protocol_proto_goTypes,
		DependencyIndexes: file_blaze_worker_worker_protocol_proto_depIdxs,
		MessageInfos:      file_blaze_worker_worker_protocol_proto_msgTypes,
	}.Build()
	File_blaze_worker_worker_protocol_proto = out.File
	file_blaze_worker_worker_protocol_proto_rawDesc = nil
	file_blaze_worker_worker_protocol_proto_goTypes = nil
	file_blaze_worker_worker_protocol_proto_depIdxs = nil
}
