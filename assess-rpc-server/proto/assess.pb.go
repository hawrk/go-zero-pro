// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.4
// source: assess.proto

package proto

import (
	context "context"
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

type AssessInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TransactTime          int64   `protobuf:"varint,1,opt,name=transact_time,json=transactTime,proto3" json:"transact_time,omitempty"`                                // 交易时间
	OrderQty              int64   `protobuf:"varint,2,opt,name=order_qty,json=orderQty,proto3" json:"order_qty,omitempty"`                                            // 委托数量
	LastQty               int64   `protobuf:"varint,3,opt,name=last_qty,json=lastQty,proto3" json:"last_qty,omitempty"`                                               //  成交数量
	CancelledQty          int64   `protobuf:"varint,4,opt,name=cancelled_qty,json=cancelledQty,proto3" json:"cancelled_qty,omitempty"`                                // 撤销数量
	RejectedQty           int64   `protobuf:"varint,5,opt,name=rejected_qty,json=rejectedQty,proto3" json:"rejected_qty,omitempty"`                                   // 拒绝数量
	Vwap                  float64 `protobuf:"fixed64,6,opt,name=vwap,proto3" json:"vwap,omitempty"`                                                                   // vwap
	VwapDeviation         float64 `protobuf:"fixed64,7,opt,name=vwap_deviation,json=vwapDeviation,proto3" json:"vwap_deviation,omitempty"`                            // vwap 滑点
	LastPrice             int64   `protobuf:"varint,8,opt,name=last_price,json=lastPrice,proto3" json:"last_price,omitempty"`                                         // 最新价格
	ArrivedPrice          int64   `protobuf:"varint,9,opt,name=arrived_price,json=arrivedPrice,proto3" json:"arrived_price,omitempty"`                                // 到达价格
	ArrivedPriceDeviation float64 `protobuf:"fixed64,10,opt,name=arrived_price_deviation,json=arrivedPriceDeviation,proto3" json:"arrived_price_deviation,omitempty"` // 到达价滑点
	MarketRate            float64 `protobuf:"fixed64,11,opt,name=market_rate,json=marketRate,proto3" json:"market_rate,omitempty"`                                    // 市场参与率
	DealRate              float64 `protobuf:"fixed64,12,opt,name=deal_rate,json=dealRate,proto3" json:"deal_rate,omitempty"`                                          // 成交量比重
	DealProgress          float64 `protobuf:"fixed64,13,opt,name=deal_progress,json=dealProgress,proto3" json:"deal_progress,omitempty"`                              // 成交进度
}

func (x *AssessInfo) Reset() {
	*x = AssessInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_assess_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AssessInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AssessInfo) ProtoMessage() {}

func (x *AssessInfo) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use AssessInfo.ProtoReflect.Descriptor instead.
func (*AssessInfo) Descriptor() ([]byte, []int) {
	return file_assess_proto_rawDescGZIP(), []int{0}
}

func (x *AssessInfo) GetTransactTime() int64 {
	if x != nil {
		return x.TransactTime
	}
	return 0
}

func (x *AssessInfo) GetOrderQty() int64 {
	if x != nil {
		return x.OrderQty
	}
	return 0
}

func (x *AssessInfo) GetLastQty() int64 {
	if x != nil {
		return x.LastQty
	}
	return 0
}

func (x *AssessInfo) GetCancelledQty() int64 {
	if x != nil {
		return x.CancelledQty
	}
	return 0
}

func (x *AssessInfo) GetRejectedQty() int64 {
	if x != nil {
		return x.RejectedQty
	}
	return 0
}

func (x *AssessInfo) GetVwap() float64 {
	if x != nil {
		return x.Vwap
	}
	return 0
}

func (x *AssessInfo) GetVwapDeviation() float64 {
	if x != nil {
		return x.VwapDeviation
	}
	return 0
}

func (x *AssessInfo) GetLastPrice() int64 {
	if x != nil {
		return x.LastPrice
	}
	return 0
}

func (x *AssessInfo) GetArrivedPrice() int64 {
	if x != nil {
		return x.ArrivedPrice
	}
	return 0
}

func (x *AssessInfo) GetArrivedPriceDeviation() float64 {
	if x != nil {
		return x.ArrivedPriceDeviation
	}
	return 0
}

func (x *AssessInfo) GetMarketRate() float64 {
	if x != nil {
		return x.MarketRate
	}
	return 0
}

func (x *AssessInfo) GetDealRate() float64 {
	if x != nil {
		return x.DealRate
	}
	return 0
}

func (x *AssessInfo) GetDealProgress() float64 {
	if x != nil {
		return x.DealProgress
	}
	return 0
}

type GeneralReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AlgoId          int32  `protobuf:"varint,1,opt,name=algo_id,json=algoId,proto3" json:"algo_id,omitempty"`                              // 算法ID
	SecId           string `protobuf:"bytes,2,opt,name=sec_id,json=secId,proto3" json:"sec_id,omitempty"`                                  // 证券ID
	TimeDemension   int32  `protobuf:"varint,3,opt,name=time_demension,json=timeDemension,proto3" json:"time_demension,omitempty"`         // 时间维度
	OrderStatusType int32  `protobuf:"varint,4,opt,name=order_status_type,json=orderStatusType,proto3" json:"order_status_type,omitempty"` // 订单状态
	StartTime       int64  `protobuf:"varint,5,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`                     // 开始时间
	EndTime         int64  `protobuf:"varint,6,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`                           // 结束时间
}

func (x *GeneralReq) Reset() {
	*x = GeneralReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_assess_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GeneralReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GeneralReq) ProtoMessage() {}

func (x *GeneralReq) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use GeneralReq.ProtoReflect.Descriptor instead.
func (*GeneralReq) Descriptor() ([]byte, []int) {
	return file_assess_proto_rawDescGZIP(), []int{1}
}

func (x *GeneralReq) GetAlgoId() int32 {
	if x != nil {
		return x.AlgoId
	}
	return 0
}

func (x *GeneralReq) GetSecId() string {
	if x != nil {
		return x.SecId
	}
	return ""
}

func (x *GeneralReq) GetTimeDemension() int32 {
	if x != nil {
		return x.TimeDemension
	}
	return 0
}

func (x *GeneralReq) GetOrderStatusType() int32 {
	if x != nil {
		return x.OrderStatusType
	}
	return 0
}

func (x *GeneralReq) GetStartTime() int64 {
	if x != nil {
		return x.StartTime
	}
	return 0
}

func (x *GeneralReq) GetEndTime() int64 {
	if x != nil {
		return x.EndTime
	}
	return 0
}

type GeneralRsp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code int32         `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg  string        `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Info []*AssessInfo `protobuf:"bytes,3,rep,name=info,proto3" json:"info,omitempty"`
}

func (x *GeneralRsp) Reset() {
	*x = GeneralRsp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_assess_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GeneralRsp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GeneralRsp) ProtoMessage() {}

func (x *GeneralRsp) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use GeneralRsp.ProtoReflect.Descriptor instead.
func (*GeneralRsp) Descriptor() ([]byte, []int) {
	return file_assess_proto_rawDescGZIP(), []int{2}
}

func (x *GeneralRsp) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *GeneralRsp) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *GeneralRsp) GetInfo() []*AssessInfo {
	if x != nil {
		return x.Info
	}
	return nil
}

var File_assess_proto protoreflect.FileDescriptor

var file_assess_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x22, 0xcb, 0x03, 0x0a, 0x0a, 0x41, 0x73, 0x73, 0x65, 0x73,
	0x73, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x23, 0x0a, 0x0d, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63,
	0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x74, 0x72,
	0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x6f, 0x72,
	0x64, 0x65, 0x72, 0x5f, 0x71, 0x74, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x6f,
	0x72, 0x64, 0x65, 0x72, 0x51, 0x74, 0x79, 0x12, 0x19, 0x0a, 0x08, 0x6c, 0x61, 0x73, 0x74, 0x5f,
	0x71, 0x74, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x6c, 0x61, 0x73, 0x74, 0x51,
	0x74, 0x79, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x6c, 0x65, 0x64, 0x5f,
	0x71, 0x74, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x63, 0x61, 0x6e, 0x63, 0x65,
	0x6c, 0x6c, 0x65, 0x64, 0x51, 0x74, 0x79, 0x12, 0x21, 0x0a, 0x0c, 0x72, 0x65, 0x6a, 0x65, 0x63,
	0x74, 0x65, 0x64, 0x5f, 0x71, 0x74, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x72,
	0x65, 0x6a, 0x65, 0x63, 0x74, 0x65, 0x64, 0x51, 0x74, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x76, 0x77,
	0x61, 0x70, 0x18, 0x06, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04, 0x76, 0x77, 0x61, 0x70, 0x12, 0x25,
	0x0a, 0x0e, 0x76, 0x77, 0x61, 0x70, 0x5f, 0x64, 0x65, 0x76, 0x69, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0d, 0x76, 0x77, 0x61, 0x70, 0x44, 0x65, 0x76, 0x69,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x0a, 0x0a, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x70, 0x72,
	0x69, 0x63, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x6c, 0x61, 0x73, 0x74, 0x50,
	0x72, 0x69, 0x63, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x61, 0x72, 0x72, 0x69, 0x76, 0x65, 0x64, 0x5f,
	0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x61, 0x72, 0x72,
	0x69, 0x76, 0x65, 0x64, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x36, 0x0a, 0x17, 0x61, 0x72, 0x72,
	0x69, 0x76, 0x65, 0x64, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x5f, 0x64, 0x65, 0x76, 0x69, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x01, 0x52, 0x15, 0x61, 0x72, 0x72, 0x69,
	0x76, 0x65, 0x64, 0x50, 0x72, 0x69, 0x63, 0x65, 0x44, 0x65, 0x76, 0x69, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x1f, 0x0a, 0x0b, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x5f, 0x72, 0x61, 0x74, 0x65,
	0x18, 0x0b, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0a, 0x6d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x52, 0x61,
	0x74, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x64, 0x65, 0x61, 0x6c, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x18,
	0x0c, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08, 0x64, 0x65, 0x61, 0x6c, 0x52, 0x61, 0x74, 0x65, 0x12,
	0x23, 0x0a, 0x0d, 0x64, 0x65, 0x61, 0x6c, 0x5f, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73,
	0x18, 0x0d, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0c, 0x64, 0x65, 0x61, 0x6c, 0x50, 0x72, 0x6f, 0x67,
	0x72, 0x65, 0x73, 0x73, 0x22, 0xc9, 0x01, 0x0a, 0x0a, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c,
	0x52, 0x65, 0x71, 0x12, 0x17, 0x0a, 0x07, 0x61, 0x6c, 0x67, 0x6f, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x61, 0x6c, 0x67, 0x6f, 0x49, 0x64, 0x12, 0x15, 0x0a, 0x06,
	0x73, 0x65, 0x63, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x65,
	0x63, 0x49, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x64, 0x65, 0x6d, 0x65,
	0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x74, 0x69, 0x6d,
	0x65, 0x44, 0x65, 0x6d, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x2a, 0x0a, 0x11, 0x6f, 0x72,
	0x64, 0x65, 0x72, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f,
	0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72,
	0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65,
	0x22, 0x5a, 0x0a, 0x0a, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x52, 0x73, 0x70, 0x12, 0x12,
	0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x63, 0x6f,
	0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6d, 0x73, 0x67, 0x12, 0x26, 0x0a, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x03, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x12, 0x2e, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x2e, 0x41, 0x73, 0x73, 0x65,
	0x73, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x32, 0x45, 0x0a, 0x0d,
	0x41, 0x73, 0x73, 0x65, 0x73, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x34, 0x0a,
	0x0a, 0x47, 0x65, 0x74, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x12, 0x12, 0x2e, 0x61, 0x73,
	0x73, 0x65, 0x73, 0x73, 0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x52, 0x65, 0x71, 0x1a,
	0x12, 0x2e, 0x61, 0x73, 0x73, 0x65, 0x73, 0x73, 0x2e, 0x47, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c,
	0x52, 0x73, 0x70, 0x42, 0x07, 0x5a, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
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

var file_assess_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_assess_proto_goTypes = []interface{}{
	(*AssessInfo)(nil), // 0: assess.AssessInfo
	(*GeneralReq)(nil), // 1: assess.GeneralReq
	(*GeneralRsp)(nil), // 2: assess.GeneralRsp
}
var file_assess_proto_depIdxs = []int32{
	0, // 0: assess.GeneralRsp.info:type_name -> assess.AssessInfo
	1, // 1: assess.AssessService.GetGeneral:input_type -> assess.GeneralReq
	2, // 2: assess.AssessService.GetGeneral:output_type -> assess.GeneralRsp
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_assess_proto_init() }
func file_assess_proto_init() {
	if File_assess_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_assess_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AssessInfo); i {
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
			switch v := v.(*GeneralReq); i {
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
			switch v := v.(*GeneralRsp); i {
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
			NumMessages:   3,
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

// AssessServiceClient is the client API for AssessService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AssessServiceClient interface {
	// 获取绩效概况
	GetGeneral(ctx context.Context, in *GeneralReq, opts ...grpc.CallOption) (*GeneralRsp, error)
}

type assessServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAssessServiceClient(cc grpc.ClientConnInterface) AssessServiceClient {
	return &assessServiceClient{cc}
}

func (c *assessServiceClient) GetGeneral(ctx context.Context, in *GeneralReq, opts ...grpc.CallOption) (*GeneralRsp, error) {
	out := new(GeneralRsp)
	err := c.cc.Invoke(ctx, "/assess.AssessService/GetGeneral", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AssessServiceServer is the server API for AssessService service.
type AssessServiceServer interface {
	// 获取绩效概况
	GetGeneral(context.Context, *GeneralReq) (*GeneralRsp, error)
}

// UnimplementedAssessServiceServer can be embedded to have forward compatible implementations.
type UnimplementedAssessServiceServer struct {
}

func (*UnimplementedAssessServiceServer) GetGeneral(context.Context, *GeneralReq) (*GeneralRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGeneral not implemented")
}

func RegisterAssessServiceServer(s *grpc.Server, srv AssessServiceServer) {
	s.RegisterService(&_AssessService_serviceDesc, srv)
}

func _AssessService_GetGeneral_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GeneralReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AssessServiceServer).GetGeneral(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/assess.AssessService/GetGeneral",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AssessServiceServer).GetGeneral(ctx, req.(*GeneralReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _AssessService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "assess.AssessService",
	HandlerType: (*AssessServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetGeneral",
			Handler:    _AssessService_GetGeneral_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "assess.proto",
}
