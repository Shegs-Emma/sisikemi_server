// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.20.3
// source: rpc_create_collection.proto

package pb

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

type CreateCollectionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CollectionName        string `protobuf:"bytes,1,opt,name=collection_name,json=collectionName,proto3" json:"collection_name,omitempty"`
	CollectionDescription string `protobuf:"bytes,2,opt,name=collection_description,json=collectionDescription,proto3" json:"collection_description,omitempty"`
	ThumbnailImage        string `protobuf:"bytes,3,opt,name=thumbnail_image,json=thumbnailImage,proto3" json:"thumbnail_image,omitempty"`
	HeaderImage           string `protobuf:"bytes,4,opt,name=header_image,json=headerImage,proto3" json:"header_image,omitempty"`
}

func (x *CreateCollectionRequest) Reset() {
	*x = CreateCollectionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_create_collection_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateCollectionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateCollectionRequest) ProtoMessage() {}

func (x *CreateCollectionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_create_collection_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateCollectionRequest.ProtoReflect.Descriptor instead.
func (*CreateCollectionRequest) Descriptor() ([]byte, []int) {
	return file_rpc_create_collection_proto_rawDescGZIP(), []int{0}
}

func (x *CreateCollectionRequest) GetCollectionName() string {
	if x != nil {
		return x.CollectionName
	}
	return ""
}

func (x *CreateCollectionRequest) GetCollectionDescription() string {
	if x != nil {
		return x.CollectionDescription
	}
	return ""
}

func (x *CreateCollectionRequest) GetThumbnailImage() string {
	if x != nil {
		return x.ThumbnailImage
	}
	return ""
}

func (x *CreateCollectionRequest) GetHeaderImage() string {
	if x != nil {
		return x.HeaderImage
	}
	return ""
}

type CreateCollectionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Collection *Collection `protobuf:"bytes,1,opt,name=collection,proto3" json:"collection,omitempty"`
}

func (x *CreateCollectionResponse) Reset() {
	*x = CreateCollectionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_create_collection_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateCollectionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateCollectionResponse) ProtoMessage() {}

func (x *CreateCollectionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_create_collection_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateCollectionResponse.ProtoReflect.Descriptor instead.
func (*CreateCollectionResponse) Descriptor() ([]byte, []int) {
	return file_rpc_create_collection_proto_rawDescGZIP(), []int{1}
}

func (x *CreateCollectionResponse) GetCollection() *Collection {
	if x != nil {
		return x.Collection
	}
	return nil
}

var File_rpc_create_collection_proto protoreflect.FileDescriptor

var file_rpc_create_collection_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x72, 0x70, 0x63, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x63, 0x6f, 0x6c,
	0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70,
	0x62, 0x1a, 0x10, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xc5, 0x01, 0x0a, 0x17, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6f,
	0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x27, 0x0a, 0x0f, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x35, 0x0a, 0x16, 0x63, 0x6f, 0x6c, 0x6c,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x15, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x27, 0x0a, 0x0f, 0x74, 0x68, 0x75, 0x6d, 0x62, 0x6e, 0x61, 0x69, 0x6c, 0x5f, 0x69, 0x6d, 0x61,
	0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x74, 0x68, 0x75, 0x6d, 0x62, 0x6e,
	0x61, 0x69, 0x6c, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x68, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x5f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x22, 0x4a, 0x0a, 0x18, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x0a, 0x63, 0x6f, 0x6c, 0x6c, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x62,
	0x2e, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x63, 0x6f, 0x6c,
	0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x25, 0x5a, 0x23, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x65, 0x63, 0x68, 0x73, 0x63, 0x68, 0x6f, 0x6f, 0x6c,
	0x2f, 0x73, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x62, 0x61, 0x6e, 0x6b, 0x2f, 0x70, 0x62, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_create_collection_proto_rawDescOnce sync.Once
	file_rpc_create_collection_proto_rawDescData = file_rpc_create_collection_proto_rawDesc
)

func file_rpc_create_collection_proto_rawDescGZIP() []byte {
	file_rpc_create_collection_proto_rawDescOnce.Do(func() {
		file_rpc_create_collection_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_create_collection_proto_rawDescData)
	})
	return file_rpc_create_collection_proto_rawDescData
}

var file_rpc_create_collection_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_rpc_create_collection_proto_goTypes = []any{
	(*CreateCollectionRequest)(nil),  // 0: pb.CreateCollectionRequest
	(*CreateCollectionResponse)(nil), // 1: pb.CreateCollectionResponse
	(*Collection)(nil),               // 2: pb.Collection
}
var file_rpc_create_collection_proto_depIdxs = []int32{
	2, // 0: pb.CreateCollectionResponse.collection:type_name -> pb.Collection
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_rpc_create_collection_proto_init() }
func file_rpc_create_collection_proto_init() {
	if File_rpc_create_collection_proto != nil {
		return
	}
	file_collection_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_rpc_create_collection_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*CreateCollectionRequest); i {
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
		file_rpc_create_collection_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*CreateCollectionResponse); i {
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
			RawDescriptor: file_rpc_create_collection_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_rpc_create_collection_proto_goTypes,
		DependencyIndexes: file_rpc_create_collection_proto_depIdxs,
		MessageInfos:      file_rpc_create_collection_proto_msgTypes,
	}.Build()
	File_rpc_create_collection_proto = out.File
	file_rpc_create_collection_proto_rawDesc = nil
	file_rpc_create_collection_proto_goTypes = nil
	file_rpc_create_collection_proto_depIdxs = nil
}