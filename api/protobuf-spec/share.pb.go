// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.19.1
// source: share.proto

package protobuf_spec

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AddShareReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PostContent *PostItem `protobuf:"bytes,1,opt,name=post_content,json=postContent,proto3" json:"post_content,omitempty"`
}

func (x *AddShareReq) Reset() {
	*x = AddShareReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_share_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddShareReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddShareReq) ProtoMessage() {}

func (x *AddShareReq) ProtoReflect() protoreflect.Message {
	mi := &file_share_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddShareReq.ProtoReflect.Descriptor instead.
func (*AddShareReq) Descriptor() ([]byte, []int) {
	return file_share_proto_rawDescGZIP(), []int{0}
}

func (x *AddShareReq) GetPostContent() *PostItem {
	if x != nil {
		return x.PostContent
	}
	return nil
}

type AddShareResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *AddShareResp) Reset() {
	*x = AddShareResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_share_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddShareResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddShareResp) ProtoMessage() {}

func (x *AddShareResp) ProtoReflect() protoreflect.Message {
	mi := &file_share_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddShareResp.ProtoReflect.Descriptor instead.
func (*AddShareResp) Descriptor() ([]byte, []int) {
	return file_share_proto_rawDescGZIP(), []int{1}
}

func (x *AddShareResp) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type GetShareByPageReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PageIndex int64 `protobuf:"varint,1,opt,name=page_index,json=pageIndex,proto3" json:"page_index,omitempty"`
	PageSize  int64 `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
}

func (x *GetShareByPageReq) Reset() {
	*x = GetShareByPageReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_share_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetShareByPageReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetShareByPageReq) ProtoMessage() {}

func (x *GetShareByPageReq) ProtoReflect() protoreflect.Message {
	mi := &file_share_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetShareByPageReq.ProtoReflect.Descriptor instead.
func (*GetShareByPageReq) Descriptor() ([]byte, []int) {
	return file_share_proto_rawDescGZIP(), []int{2}
}

func (x *GetShareByPageReq) GetPageIndex() int64 {
	if x != nil {
		return x.PageIndex
	}
	return 0
}

func (x *GetShareByPageReq) GetPageSize() int64 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

type GetShareByPageResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Total int64       `protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`
	Data  []*PostItem `protobuf:"bytes,2,rep,name=data,proto3" json:"data,omitempty"`
}

func (x *GetShareByPageResp) Reset() {
	*x = GetShareByPageResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_share_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetShareByPageResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetShareByPageResp) ProtoMessage() {}

func (x *GetShareByPageResp) ProtoReflect() protoreflect.Message {
	mi := &file_share_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetShareByPageResp.ProtoReflect.Descriptor instead.
func (*GetShareByPageResp) Descriptor() ([]byte, []int) {
	return file_share_proto_rawDescGZIP(), []int{3}
}

func (x *GetShareByPageResp) GetTotal() int64 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *GetShareByPageResp) GetData() []*PostItem {
	if x != nil {
		return x.Data
	}
	return nil
}

type PostItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Text      string                 `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
	Img       []string               `protobuf:"bytes,2,rep,name=img,proto3" json:"img,omitempty"`
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=CreatedAt,proto3" json:"CreatedAt,omitempty"`
}

func (x *PostItem) Reset() {
	*x = PostItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_share_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PostItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PostItem) ProtoMessage() {}

func (x *PostItem) ProtoReflect() protoreflect.Message {
	mi := &file_share_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PostItem.ProtoReflect.Descriptor instead.
func (*PostItem) Descriptor() ([]byte, []int) {
	return file_share_proto_rawDescGZIP(), []int{4}
}

func (x *PostItem) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *PostItem) GetImg() []string {
	if x != nil {
		return x.Img
	}
	return nil
}

func (x *PostItem) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

var File_share_proto protoreflect.FileDescriptor

var file_share_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x73, 0x68, 0x61, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3b,
	0x0a, 0x0b, 0x41, 0x64, 0x64, 0x53, 0x68, 0x61, 0x72, 0x65, 0x52, 0x65, 0x71, 0x12, 0x2c, 0x0a,
	0x0c, 0x70, 0x6f, 0x73, 0x74, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x0b,
	0x70, 0x6f, 0x73, 0x74, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x28, 0x0a, 0x0c, 0x41,
	0x64, 0x64, 0x53, 0x68, 0x61, 0x72, 0x65, 0x52, 0x65, 0x73, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x4f, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x53, 0x68, 0x61, 0x72,
	0x65, 0x42, 0x79, 0x50, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x61,
	0x67, 0x65, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09,
	0x70, 0x61, 0x67, 0x65, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67,
	0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x70, 0x61,
	0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x22, 0x49, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x53, 0x68, 0x61,
	0x72, 0x65, 0x42, 0x79, 0x50, 0x61, 0x67, 0x65, 0x52, 0x65, 0x73, 0x70, 0x12, 0x14, 0x0a, 0x05,
	0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x74, 0x6f, 0x74,
	0x61, 0x6c, 0x12, 0x1d, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x09, 0x2e, 0x50, 0x6f, 0x73, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x22, 0x6a, 0x0a, 0x08, 0x50, 0x6f, 0x73, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x12, 0x0a,
	0x04, 0x74, 0x65, 0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78,
	0x74, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x6d, 0x67, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03,
	0x69, 0x6d, 0x67, 0x12, 0x38, 0x0a, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x32, 0x76, 0x0a,
	0x0c, 0x53, 0x68, 0x61, 0x72, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x29, 0x0a,
	0x08, 0x41, 0x64, 0x64, 0x53, 0x68, 0x61, 0x72, 0x65, 0x12, 0x0c, 0x2e, 0x41, 0x64, 0x64, 0x53,
	0x68, 0x61, 0x72, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x0d, 0x2e, 0x41, 0x64, 0x64, 0x53, 0x68, 0x61,
	0x72, 0x65, 0x52, 0x65, 0x73, 0x70, 0x22, 0x00, 0x12, 0x3b, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x53,
	0x68, 0x61, 0x72, 0x65, 0x42, 0x79, 0x50, 0x61, 0x67, 0x65, 0x12, 0x12, 0x2e, 0x47, 0x65, 0x74,
	0x53, 0x68, 0x61, 0x72, 0x65, 0x42, 0x79, 0x50, 0x61, 0x67, 0x65, 0x52, 0x65, 0x71, 0x1a, 0x13,
	0x2e, 0x47, 0x65, 0x74, 0x53, 0x68, 0x61, 0x72, 0x65, 0x42, 0x79, 0x50, 0x61, 0x67, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x22, 0x00, 0x42, 0x13, 0x5a, 0x11, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2d, 0x73, 0x70, 0x65, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_share_proto_rawDescOnce sync.Once
	file_share_proto_rawDescData = file_share_proto_rawDesc
)

func file_share_proto_rawDescGZIP() []byte {
	file_share_proto_rawDescOnce.Do(func() {
		file_share_proto_rawDescData = protoimpl.X.CompressGZIP(file_share_proto_rawDescData)
	})
	return file_share_proto_rawDescData
}

var file_share_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_share_proto_goTypes = []interface{}{
	(*AddShareReq)(nil),           // 0: AddShareReq
	(*AddShareResp)(nil),          // 1: AddShareResp
	(*GetShareByPageReq)(nil),     // 2: GetShareByPageReq
	(*GetShareByPageResp)(nil),    // 3: GetShareByPageResp
	(*PostItem)(nil),              // 4: PostItem
	(*timestamppb.Timestamp)(nil), // 5: google.protobuf.Timestamp
}
var file_share_proto_depIdxs = []int32{
	4, // 0: AddShareReq.post_content:type_name -> PostItem
	4, // 1: GetShareByPageResp.data:type_name -> PostItem
	5, // 2: PostItem.CreatedAt:type_name -> google.protobuf.Timestamp
	0, // 3: ShareService.AddShare:input_type -> AddShareReq
	2, // 4: ShareService.GetShareByPage:input_type -> GetShareByPageReq
	1, // 5: ShareService.AddShare:output_type -> AddShareResp
	3, // 6: ShareService.GetShareByPage:output_type -> GetShareByPageResp
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_share_proto_init() }
func file_share_proto_init() {
	if File_share_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_share_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddShareReq); i {
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
		file_share_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddShareResp); i {
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
		file_share_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetShareByPageReq); i {
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
		file_share_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetShareByPageResp); i {
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
		file_share_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PostItem); i {
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
			RawDescriptor: file_share_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_share_proto_goTypes,
		DependencyIndexes: file_share_proto_depIdxs,
		MessageInfos:      file_share_proto_msgTypes,
	}.Build()
	File_share_proto = out.File
	file_share_proto_rawDesc = nil
	file_share_proto_goTypes = nil
	file_share_proto_depIdxs = nil
}
