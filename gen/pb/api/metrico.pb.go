// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: proto/metrico.proto

package api

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

type MetricType int32

const (
	MetricType_GAUGE   MetricType = 0
	MetricType_COUNTER MetricType = 1
)

// Enum value maps for MetricType.
var (
	MetricType_name = map[int32]string{
		0: "GAUGE",
		1: "COUNTER",
	}
	MetricType_value = map[string]int32{
		"GAUGE":   0,
		"COUNTER": 1,
	}
)

func (x MetricType) Enum() *MetricType {
	p := new(MetricType)
	*p = x
	return p
}

func (x MetricType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MetricType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_metrico_proto_enumTypes[0].Descriptor()
}

func (MetricType) Type() protoreflect.EnumType {
	return &file_proto_metrico_proto_enumTypes[0]
}

func (x MetricType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MetricType.Descriptor instead.
func (MetricType) EnumDescriptor() ([]byte, []int) {
	return file_proto_metrico_proto_rawDescGZIP(), []int{0}
}

type Status int32

const (
	Status_OK    Status = 0
	Status_ERROR Status = 1
)

// Enum value maps for Status.
var (
	Status_name = map[int32]string{
		0: "OK",
		1: "ERROR",
	}
	Status_value = map[string]int32{
		"OK":    0,
		"ERROR": 1,
	}
)

func (x Status) Enum() *Status {
	p := new(Status)
	*p = x
	return p
}

func (x Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Status) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_metrico_proto_enumTypes[1].Descriptor()
}

func (Status) Type() protoreflect.EnumType {
	return &file_proto_metrico_proto_enumTypes[1]
}

func (x Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Status.Descriptor instead.
func (Status) EnumDescriptor() ([]byte, []int) {
	return file_proto_metrico_proto_rawDescGZIP(), []int{1}
}

type Metric struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string     `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Type  MetricType `protobuf:"varint,2,opt,name=type,proto3,enum=com.github.tony_spark.metrico.MetricType" json:"type,omitempty"`
	Delta *int64     `protobuf:"varint,3,opt,name=delta,proto3,oneof" json:"delta,omitempty"`
	Value *float64   `protobuf:"fixed64,4,opt,name=value,proto3,oneof" json:"value,omitempty"`
	Hash  []byte     `protobuf:"bytes,5,opt,name=hash,proto3,oneof" json:"hash,omitempty"`
}

func (x *Metric) Reset() {
	*x = Metric{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metrico_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metric) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metric) ProtoMessage() {}

func (x *Metric) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metrico_proto_msgTypes[0]
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
	return file_proto_metrico_proto_rawDescGZIP(), []int{0}
}

func (x *Metric) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Metric) GetType() MetricType {
	if x != nil {
		return x.Type
	}
	return MetricType_GAUGE
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

func (x *Metric) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metrico_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metrico_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_proto_metrico_proto_rawDescGZIP(), []int{1}
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status Status  `protobuf:"varint,1,opt,name=status,proto3,enum=com.github.tony_spark.metrico.Status" json:"status,omitempty"`
	Error  *string `protobuf:"bytes,2,opt,name=error,proto3,oneof" json:"error,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_metrico_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_proto_metrico_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_proto_metrico_proto_rawDescGZIP(), []int{2}
}

func (x *Response) GetStatus() Status {
	if x != nil {
		return x.Status
	}
	return Status_OK
}

func (x *Response) GetError() string {
	if x != nil && x.Error != nil {
		return *x.Error
	}
	return ""
}

var File_proto_metrico_proto protoreflect.FileDescriptor

var file_proto_metrico_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x6f, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1d, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x74, 0x6f, 0x6e, 0x79, 0x5f, 0x73, 0x70, 0x61, 0x72, 0x6b, 0x2e, 0x6d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x6f, 0x22, 0xc3, 0x01, 0x0a, 0x06, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x3d, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x29, 0x2e,
	0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x74, 0x6f, 0x6e, 0x79, 0x5f,
	0x73, 0x70, 0x61, 0x72, 0x6b, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x6f, 0x2e, 0x4d, 0x65,
	0x74, 0x72, 0x69, 0x63, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x19,
	0x0a, 0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52,
	0x05, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x48, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x88, 0x01, 0x01, 0x12, 0x17, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0c, 0x48, 0x02, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x88, 0x01, 0x01, 0x42, 0x08, 0x0a,
	0x06, 0x5f, 0x64, 0x65, 0x6c, 0x74, 0x61, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x22, 0x6e, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x3d, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x25, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x74, 0x6f, 0x6e,
	0x79, 0x5f, 0x73, 0x70, 0x61, 0x72, 0x6b, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x6f, 0x2e,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x19,
	0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52,
	0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x88, 0x01, 0x01, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x65, 0x72,
	0x72, 0x6f, 0x72, 0x2a, 0x24, 0x0a, 0x0a, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x09, 0x0a, 0x05, 0x47, 0x41, 0x55, 0x47, 0x45, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07,
	0x43, 0x4f, 0x55, 0x4e, 0x54, 0x45, 0x52, 0x10, 0x01, 0x2a, 0x1b, 0x0a, 0x06, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x06, 0x0a, 0x02, 0x4f, 0x4b, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x45,
	0x52, 0x52, 0x4f, 0x52, 0x10, 0x01, 0x32, 0xca, 0x01, 0x0a, 0x0d, 0x4d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5c, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x12, 0x25, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x74, 0x6f, 0x6e, 0x79, 0x5f, 0x73, 0x70, 0x61, 0x72, 0x6b, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x6f, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x1a, 0x27, 0x2e, 0x63, 0x6f, 0x6d, 0x2e,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x74, 0x6f, 0x6e, 0x79, 0x5f, 0x73, 0x70, 0x61, 0x72,
	0x6b, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x6f, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x12, 0x5b, 0x0a, 0x08, 0x44, 0x42, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x24, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x74, 0x6f, 0x6e, 0x79, 0x5f, 0x73, 0x70, 0x61, 0x72, 0x6b, 0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69,
	0x63, 0x6f, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x27, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x74, 0x6f, 0x6e, 0x79, 0x5f, 0x73, 0x70, 0x61, 0x72, 0x6b,
	0x2e, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x6f, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x42, 0x0c, 0x5a, 0x0a, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x62, 0x2f, 0x61, 0x70,
	0x69, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_metrico_proto_rawDescOnce sync.Once
	file_proto_metrico_proto_rawDescData = file_proto_metrico_proto_rawDesc
)

func file_proto_metrico_proto_rawDescGZIP() []byte {
	file_proto_metrico_proto_rawDescOnce.Do(func() {
		file_proto_metrico_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_metrico_proto_rawDescData)
	})
	return file_proto_metrico_proto_rawDescData
}

var file_proto_metrico_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_proto_metrico_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_proto_metrico_proto_goTypes = []interface{}{
	(MetricType)(0),  // 0: com.github.tony_spark.metrico.MetricType
	(Status)(0),      // 1: com.github.tony_spark.metrico.Status
	(*Metric)(nil),   // 2: com.github.tony_spark.metrico.Metric
	(*Empty)(nil),    // 3: com.github.tony_spark.metrico.Empty
	(*Response)(nil), // 4: com.github.tony_spark.metrico.Response
}
var file_proto_metrico_proto_depIdxs = []int32{
	0, // 0: com.github.tony_spark.metrico.Metric.type:type_name -> com.github.tony_spark.metrico.MetricType
	1, // 1: com.github.tony_spark.metrico.Response.status:type_name -> com.github.tony_spark.metrico.Status
	2, // 2: com.github.tony_spark.metrico.MetricService.Update:input_type -> com.github.tony_spark.metrico.Metric
	3, // 3: com.github.tony_spark.metrico.MetricService.DBStatus:input_type -> com.github.tony_spark.metrico.Empty
	4, // 4: com.github.tony_spark.metrico.MetricService.Update:output_type -> com.github.tony_spark.metrico.Response
	4, // 5: com.github.tony_spark.metrico.MetricService.DBStatus:output_type -> com.github.tony_spark.metrico.Response
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_metrico_proto_init() }
func file_proto_metrico_proto_init() {
	if File_proto_metrico_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_metrico_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_proto_metrico_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
		file_proto_metrico_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
	file_proto_metrico_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_proto_metrico_proto_msgTypes[2].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_metrico_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_metrico_proto_goTypes,
		DependencyIndexes: file_proto_metrico_proto_depIdxs,
		EnumInfos:         file_proto_metrico_proto_enumTypes,
		MessageInfos:      file_proto_metrico_proto_msgTypes,
	}.Build()
	File_proto_metrico_proto = out.File
	file_proto_metrico_proto_rawDesc = nil
	file_proto_metrico_proto_goTypes = nil
	file_proto_metrico_proto_depIdxs = nil
}