// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.22.0
// 	protoc        v3.11.4
// source: assess.proto

package proto

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type DemoReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hello string `protobuf:"bytes,1,opt,name=hello,proto3" json:"hello,omitempty"`
}

func (x *DemoReq) Reset() {
	*x = DemoReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_assess_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DemoReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DemoReq) ProtoMessage() {}

func (x *DemoReq) ProtoReflect() protoreflect.Message {
	mi := &file_assess_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DemoReq.ProtoReflect.Descriptor instead.
func (*DemoReq) Descriptor() ([]byte, []int) {
	return file_assess_proto_rawDescGZIP(), []int{0}
}

func (x *DemoReq) GetHello() string {
	if x != nil {
		return x.Hello
	}
	return ""
}

type DemoRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Reply string `protobuf:"bytes,1,opt,name=reply,proto3" json:"reply,omitempty"`
}

func (x *DemoRsp) Reset() {
	*x = DemoRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_assess_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DemoRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DemoRsp) ProtoMessage() {}

func (x *DemoRsp) ProtoReflect() protoreflect.Message {
	mi := &file_assess_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DemoRsp.ProtoReflect.Descriptor instead.
func (*DemoRsp) Descriptor() ([]byte, []int) {
	return file_assess_proto_rawDescGZIP(), []int{1}
}

func (x *DemoRsp) GetReply() string {
	if x != nil {
		return x.Reply
	}
	return ""
}

type OverviewReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StartDate string `protobuf:"bytes,1,opt,name=start_date,json=startDate,proto3" json:"start_date,omitempty"`
	EndDate   string `protobuf:"bytes,2,opt,name=end_date,json=endDate,proto3" json:"end_date,omitempty"`
	AlgoType  int32  `protobuf:"varint,3,opt,name=algo_type,json=algoType,proto3" json:"algo_type,omitempty"`
}

func (x *OverviewReq) Reset() {
	*x = OverviewReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_assess_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OverviewReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OverviewReq) ProtoMessage() {}

func (x *OverviewReq) ProtoReflect() protoreflect.Message {
	mi := &file_assess_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OverviewReq.ProtoReflect.Descriptor instead.
func (*OverviewReq) Descriptor() ([]byte, []int) {
	return file_assess_proto_rawDescGZIP(), []int{2}
}

func (x *OverviewReq) GetStartDate() string {
	if x != nil {
		return x.StartDate
	}
	return ""
}

func (x *OverviewReq) GetEndDate() string {
	if x != nil {
		return x.EndDate
	}
	return ""
}

func (x *OverviewReq) GetAlgoType() int32 {
	if x != nil {
		return x.AlgoType
	}
	return 0
}

type OverVewRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code int32              `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg  string             `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Data []*OverVewRsp_Data `protobuf:"bytes,3,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *OverVewRsp) Reset() {
	*x = OverVewRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_assess_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OverVewRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OverVewRsp) ProtoMessage() {}

func (x *OverVewRsp) ProtoReflect() protoreflect.Message {
	mi := &file_assess_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OverVewRsp.ProtoReflect.Descriptor instead.
func (*OverVewRsp) Descriptor() ([]byte, []int) {
	return file_assess_proto_rawDescGZIP(), []int{3}
}

func (x *OverVewRsp) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *OverVewRsp) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *OverVewRsp) GetData() []*OverVewRsp_Data {
	if x != nil {
		return x.Data
	}
	return nil
}

type OverVewRsp_Data struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AlgoName   string                  `protobuf:"bytes,1,opt,name=algo_name,json=algoName,proto3" json:"algo_name,omitempty"`
	AssessType int32                   `protobuf:"varint,2,opt,name=assess_type,json=assessType,proto3" json:"assess_type,omitempty"`
	Info       []*OverVewRsp_Data_Info `protobuf:"bytes,3,rep,name=info,proto3" json:"info,omitempty"`
}

func (x *OverVewRsp_Data) Reset() {
	*x = OverVewRsp_Data{}
	if protoimpl.UnsafeEnabled {
		mi := &file_assess_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OverVewRsp_Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OverVewRsp_Data) ProtoMessage() {}

func (x *OverVewRsp_Data) ProtoReflect() protoreflect.Message {
	mi := &file_assess_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OverVewRsp_Data.ProtoReflect.Descriptor instead.
func (*OverVewRsp_Data) Descriptor() ([]byte, []int) {
	return file_assess_proto_rawDescGZIP(), []int{3, 0}
}

func (x *OverVewRsp_Data) GetAlgoName() string {
	if x != nil {
		return x.AlgoName
	}
	return ""
}

func (x *OverVewRsp_Data) GetAssessType() int32 {
	if x != nil {
		return x.AssessType
	}
	return 0
}

func (x *OverVewRsp_Data) GetInfo() []*OverVewRsp_Data_Info {
	if x != nil {
		return x.Info
	}
	return nil
}

type OverVewRsp_Data_Info struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Date  string  `protobuf:"bytes,1,opt,name=date,proto3" json:"date,omitempty"`
	Score float32 `protobuf:"fixed32,2,opt,name=score,proto3" json:"score,omitempty"`
}

func (x *OverVewRsp_Data_Info) Reset() {
	*x = OverVewRsp_Data_Info{}
	if protoimpl.UnsafeEnabled {
		mi := &file_assess_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OverVewRsp_Data_Info) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OverVewRsp_Data_Info) ProtoMessage() {}

func (x *OverVewRsp_Data_Info) ProtoReflect() protoreflect.Message {
	mi := &file_assess_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OverVewRsp_Data_Info.ProtoReflect.Descriptor instead.
func (*OverVewRsp_Data_Info) Descriptor() ([]byte, []int) {
	return file_assess_proto_rawDescGZIP(), []int{3, 0, 0}
}

func (x *OverVewRsp_Data_Info) GetDate() string {
	if x != nil {
		return x.Date
	}
	return ""
}

func (x *OverVewRsp_Data_Info) GetScore() float32 {
	if x != nil {
		return x.Score
	}
	return 0
}

var File_assess_proto protoreflect.FileDescriptor

var file_assess_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x22, 0x1f, 0x0a, 0x07, 0x44, 0x65, 0x6d, 0x6f, 0x52, 0x65,
	0x71, 0x12, 0x14, 0x0a, 0x05, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x22, 0x1f, 0x0a, 0x07, 0x44, 0x65, 0x6d, 0x6f, 0x52,
	0x73, 0x70, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x65, 0x70, 0x6c, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x72, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x64, 0x0a, 0x0b, 0x4f, 0x76, 0x65, 0x72,
	0x76, 0x69, 0x65, 0x77, 0x52, 0x65, 0x71, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x5f, 0x64, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x44, 0x61, 0x74, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x64, 0x61,
	0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x44, 0x61, 0x74,
	0x65, 0x12, 0x1b, 0x0a, 0x09, 0x61, 0x6c, 0x67, 0x6f, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x61, 0x6c, 0x67, 0x6f, 0x54, 0x79, 0x70, 0x65, 0x22, 0x8a,
	0x02, 0x0a, 0x0a, 0x4f, 0x76, 0x65, 0x72, 0x56, 0x65, 0x77, 0x52, 0x73, 0x70, 0x12, 0x12, 0x0a,
	0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f, 0x64,
	0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6d, 0x73, 0x67, 0x12, 0x2b, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x17, 0x2e, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x2e, 0x4f, 0x76, 0x65, 0x72, 0x56,
	0x65, 0x77, 0x52, 0x73, 0x70, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x1a, 0xa8, 0x01, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1b, 0x0a, 0x09, 0x61, 0x6c, 0x67,
	0x6f, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x61, 0x6c,
	0x67, 0x6f, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73,
	0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x61, 0x73, 0x73,
	0x65, 0x73, 0x73, 0x54, 0x79, 0x70, 0x65, 0x12, 0x30, 0x0a, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x2e, 0x4f,
	0x76, 0x65, 0x72, 0x56, 0x65, 0x77, 0x52, 0x73, 0x70, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x49,
	0x6e, 0x66, 0x6f, 0x52, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x1a, 0x30, 0x0a, 0x04, 0x49, 0x6e, 0x66,
	0x6f, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x64, 0x61, 0x74, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x02, 0x52, 0x05, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x32, 0x6d, 0x0a, 0x06, 0x41,
	0x73, 0x73, 0x65, 0x73, 0x73, 0x12, 0x2b, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x44, 0x65, 0x6d, 0x6f,
	0x12, 0x0f, 0x2e, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x2e, 0x44, 0x65, 0x6d, 0x6f, 0x52, 0x65,
	0x71, 0x1a, 0x0f, 0x2e, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x2e, 0x44, 0x65, 0x6d, 0x6f, 0x52,
	0x73, 0x70, 0x12, 0x36, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x4f, 0x76, 0x65, 0x72, 0x76, 0x69, 0x65,
	0x77, 0x12, 0x13, 0x2e, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x2e, 0x4f, 0x76, 0x65, 0x72, 0x76,
	0x69, 0x65, 0x77, 0x52, 0x65, 0x71, 0x1a, 0x12, 0x2e, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x2e,
	0x4f, 0x76, 0x65, 0x72, 0x56, 0x65, 0x77, 0x52, 0x73, 0x70, 0x42, 0x07, 0x5a, 0x05, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_assess_proto_rawDescOnce sync.Once
	file_assess_proto_rawDescData = file_assess_proto_rawDesc
)

func file_assess_proto_rawDescGZIP() []byte {
	file_assess_proto_rawDescOnce.Do(func() {
		file_assess_proto_rawDescData = protoimpl.X.CompressGZIP(file_assess_proto_rawDescData)
	})
	return file_assess_proto_rawDescData
}

var file_assess_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_assess_proto_goTypes = []interface{}{
	(*DemoReq)(nil),              // 0: assess.DemoReq
	(*DemoRsp)(nil),              // 1: assess.DemoRsp
	(*OverviewReq)(nil),          // 2: assess.OverviewReq
	(*OverVewRsp)(nil),           // 3: assess.OverVewRsp
	(*OverVewRsp_Data)(nil),      // 4: assess.OverVewRsp.Data
	(*OverVewRsp_Data_Info)(nil), // 5: assess.OverVewRsp.Data.Info
}
var file_assess_proto_depIdxs = []int32{
	4, // 0: assess.OverVewRsp.data:type_name -> assess.OverVewRsp.Data
	5, // 1: assess.OverVewRsp.Data.info:type_name -> assess.OverVewRsp.Data.Info
	0, // 2: assess.Assess.GetDemo:input_type -> assess.DemoReq
	2, // 3: assess.Assess.GetOverview:input_type -> assess.OverviewReq
	1, // 4: assess.Assess.GetDemo:output_type -> assess.DemoRsp
	3, // 5: assess.Assess.GetOverview:output_type -> assess.OverVewRsp
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_assess_proto_init() }
func file_assess_proto_init() {
	if File_assess_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_assess_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DemoReq); i {
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
		file_assess_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DemoRsp); i {
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
		file_assess_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OverviewReq); i {
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
		file_assess_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OverVewRsp); i {
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
		file_assess_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OverVewRsp_Data); i {
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
		file_assess_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OverVewRsp_Data_Info); i {
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
			RawDescriptor: file_assess_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_assess_proto_goTypes,
		DependencyIndexes: file_assess_proto_depIdxs,
		MessageInfos:      file_assess_proto_msgTypes,
	}.Build()
	File_assess_proto = out.File
	file_assess_proto_rawDesc = nil
	file_assess_proto_goTypes = nil
	file_assess_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AssessClient is the client API for Assess service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AssessClient interface {
	GetDemo(ctx context.Context, in *DemoReq, opts ...grpc.CallOption) (*DemoRsp, error)
	GetOverview(ctx context.Context, in *OverviewReq, opts ...grpc.CallOption) (*OverVewRsp, error)
}

type assessClient struct {
	cc grpc.ClientConnInterface
}

func NewAssessClient(cc grpc.ClientConnInterface) AssessClient {
	return &assessClient{cc}
}

func (c *assessClient) GetDemo(ctx context.Context, in *DemoReq, opts ...grpc.CallOption) (*DemoRsp, error) {
	out := new(DemoRsp)
	err := c.cc.Invoke(ctx, "/assess.Assess/GetDemo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *assessClient) GetOverview(ctx context.Context, in *OverviewReq, opts ...grpc.CallOption) (*OverVewRsp, error) {
	out := new(OverVewRsp)
	err := c.cc.Invoke(ctx, "/assess.Assess/GetOverview", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AssessServer is the server API for Assess service.
type AssessServer interface {
	GetDemo(context.Context, *DemoReq) (*DemoRsp, error)
	GetOverview(context.Context, *OverviewReq) (*OverVewRsp, error)
}

// UnimplementedAssessServer can be embedded to have forward compatible implementations.
type UnimplementedAssessServer struct {
}

func (*UnimplementedAssessServer) GetDemo(context.Context, *DemoReq) (*DemoRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDemo not implemented")
}
func (*UnimplementedAssessServer) GetOverview(context.Context, *OverviewReq) (*OverVewRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOverview not implemented")
}

func RegisterAssessServer(s *grpc.Server, srv AssessServer) {
	s.RegisterService(&_Assess_serviceDesc, srv)
}

func _Assess_GetDemo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DemoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AssessServer).GetDemo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/assess.Assess/GetDemo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AssessServer).GetDemo(ctx, req.(*DemoReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Assess_GetOverview_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OverviewReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AssessServer).GetOverview(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/assess.Assess/GetOverview",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AssessServer).GetOverview(ctx, req.(*OverviewReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _Assess_serviceDesc = grpc.ServiceDesc{
	ServiceName: "assess.Assess",
	HandlerType: (*AssessServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDemo",
			Handler:    _Assess_GetDemo_Handler,
		},
		{
			MethodName: "GetOverview",
			Handler:    _Assess_GetOverview_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "assess.proto",
}
