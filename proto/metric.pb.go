// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.8
// source: proto/metric.proto

package proto

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Metric_Type int32

const (
	Metric_UNKNOWN Metric_Type = 0
	Metric_GAUGE   Metric_Type = 1
	Metric_COUNTER Metric_Type = 2
)

// Enum value maps for Metric_Type.
var (
	Metric_Type_name = map[int32]string{
		0: "UNKNOWN",
		1: "GAUGE",
		2: "COUNTER",
	}
	Metric_Type_value = map[string]int32{
		"UNKNOWN": 0,
		"GAUGE":   1,
		"COUNTER": 2,
	}
)

func (x Metric_Type) Enum() *Metric_Type {
	p := new(Metric_Type)
	*p = x
	return p
}

func (x Metric_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Metric_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_metric_proto_enumTypes[0].Descriptor()
}

func (Metric_Type) Type() protoreflect.EnumType {
	return &file_proto_metric_proto_enumTypes[0]
}

func (x Metric_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Metric_Type.Descriptor instead.
func (Metric_Type) EnumDescriptor() ([]byte, []int) {
	return file_proto_metric_proto_rawDescGZIP(), []int{0, 0}
}

type Metric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string      `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Type  Metric_Type `protobuf:"varint,2,opt,name=type,proto3,enum=proto.Metric_Type" json:"type,omitempty"`
	Delta *int64      `protobuf:"varint,3,opt,name=delta,proto3,oneof" json:"delta,omitempty"`
	Value *float64    `protobuf:"fixed64,4,opt,name=value,proto3,oneof" json:"value,omitempty"`
	Hash  *string     `protobuf:"bytes,5,opt,name=hash,proto3,oneof" json:"hash,omitempty"`
}

func (x *Metric) Reset() {
	*x = Metric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metric_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metric) ProtoMessage() {}

func (x *Metric) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metric_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metric.ProtoReflect.Descriptor instead.
func (*Metric) Descriptor() ([]byte, []int) {
	return file_proto_metric_proto_rawDescGZIP(), []int{0}
}

func (x *Metric) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Metric) GetType() Metric_Type {
	if x != nil {
		return x.Type
	}
	return Metric_UNKNOWN
}

func (x *Metric) GetDelta() int64 {
	if x != nil && x.Delta != nil {
		return *x.Delta
	}
	return 0
}

func (x *Metric) GetValue() float64 {
	if x != nil && x.Value != nil {
		return *x.Value
	}
	return 0
}

func (x *Metric) GetHash() string {
	if x != nil && x.Hash != nil {
		return *x.Hash
	}
	return ""
}

type BulkMetrics struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metrics []*Metric `protobuf:"bytes,1,rep,name=metrics,proto3" json:"metrics,omitempty"`
}

func (x *BulkMetrics) Reset() {
	*x = BulkMetrics{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metric_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BulkMetrics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BulkMetrics) ProtoMessage() {}

func (x *BulkMetrics) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metric_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BulkMetrics.ProtoReflect.Descriptor instead.
func (*BulkMetrics) Descriptor() ([]byte, []int) {
	return file_proto_metric_proto_rawDescGZIP(), []int{1}
}

func (x *BulkMetrics) GetMetrics() []*Metric {
	if x != nil {
		return x.Metrics
	}
	return nil
}

type AddBulkMetricsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metrics *BulkMetrics `protobuf:"bytes,1,opt,name=metrics,proto3" json:"metrics,omitempty"`
}

func (x *AddBulkMetricsRequest) Reset() {
	*x = AddBulkMetricsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metric_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddBulkMetricsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddBulkMetricsRequest) ProtoMessage() {}

func (x *AddBulkMetricsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metric_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddBulkMetricsRequest.ProtoReflect.Descriptor instead.
func (*AddBulkMetricsRequest) Descriptor() ([]byte, []int) {
	return file_proto_metric_proto_rawDescGZIP(), []int{2}
}

func (x *AddBulkMetricsRequest) GetMetrics() *BulkMetrics {
	if x != nil {
		return x.Metrics
	}
	return nil
}

var File_proto_metric_proto protoreflect.FileDescriptor

var file_proto_metric_proto_rawDesc = []byte{
	0x0a, 0x12, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70,
	0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd9, 0x01, 0x0a, 0x06, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x26, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x19, 0x0a, 0x05, 0x64,
	0x65, 0x6c, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x05, 0x64, 0x65,
	0x6c, 0x74, 0x61, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x01, 0x48, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x88, 0x01,
	0x01, 0x12, 0x17, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x48,
	0x02, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x88, 0x01, 0x01, 0x22, 0x2b, 0x0a, 0x04, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12,
	0x09, 0x0a, 0x05, 0x47, 0x41, 0x55, 0x47, 0x45, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x4f,
	0x55, 0x4e, 0x54, 0x45, 0x52, 0x10, 0x02, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x64, 0x65, 0x6c, 0x74,
	0x61, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x07, 0x0a, 0x05, 0x5f,
	0x68, 0x61, 0x73, 0x68, 0x22, 0x36, 0x0a, 0x0b, 0x42, 0x75, 0x6c, 0x6b, 0x4d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x73, 0x12, 0x27, 0x0a, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x22, 0x45, 0x0a, 0x15,
	0x41, 0x64, 0x64, 0x42, 0x75, 0x6c, 0x6b, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2c, 0x0a, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42,
	0x75, 0x6c, 0x6b, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x73, 0x32, 0x52, 0x0a, 0x07, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x47,
	0x0a, 0x0f, 0x42, 0x75, 0x6c, 0x6b, 0x53, 0x61, 0x76, 0x65, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x73, 0x12, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x64, 0x64, 0x42, 0x75, 0x6c,
	0x6b, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x20, 0x5a, 0x1e, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x76, 0x6c, 0x6c, 0x76, 0x6c, 0x6c, 0x2f, 0x64, 0x65, 0x76,
	0x6f, 0x70, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_proto_metric_proto_rawDescOnce sync.Once
	file_proto_metric_proto_rawDescData = file_proto_metric_proto_rawDesc
)

func file_proto_metric_proto_rawDescGZIP() []byte {
	file_proto_metric_proto_rawDescOnce.Do(func() {
		file_proto_metric_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_metric_proto_rawDescData)
	})
	return file_proto_metric_proto_rawDescData
}

var file_proto_metric_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_metric_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_metric_proto_goTypes = []interface{}{
	(Metric_Type)(0),              // 0: proto.Metric.Type
	(*Metric)(nil),                // 1: proto.Metric
	(*BulkMetrics)(nil),           // 2: proto.BulkMetrics
	(*AddBulkMetricsRequest)(nil), // 3: proto.AddBulkMetricsRequest
	(*emptypb.Empty)(nil),         // 4: google.protobuf.Empty
}
var file_proto_metric_proto_depIdxs = []int32{
	0, // 0: proto.Metric.type:type_name -> proto.Metric.Type
	1, // 1: proto.BulkMetrics.metrics:type_name -> proto.Metric
	2, // 2: proto.AddBulkMetricsRequest.metrics:type_name -> proto.BulkMetrics
	3, // 3: proto.Metrics.BulkSaveMetrics:input_type -> proto.AddBulkMetricsRequest
	4, // 4: proto.Metrics.BulkSaveMetrics:output_type -> google.protobuf.Empty
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_proto_metric_proto_init() }
func file_proto_metric_proto_init() {
	if File_proto_metric_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_metric_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metric); i {
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
		file_proto_metric_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BulkMetrics); i {
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
		file_proto_metric_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddBulkMetricsRequest); i {
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
	file_proto_metric_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_metric_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_metric_proto_goTypes,
		DependencyIndexes: file_proto_metric_proto_depIdxs,
		EnumInfos:         file_proto_metric_proto_enumTypes,
		MessageInfos:      file_proto_metric_proto_msgTypes,
	}.Build()
	File_proto_metric_proto = out.File
	file_proto_metric_proto_rawDesc = nil
	file_proto_metric_proto_goTypes = nil
	file_proto_metric_proto_depIdxs = nil
}