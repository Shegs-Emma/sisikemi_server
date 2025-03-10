// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v5.27.3
// source: rpc_update_product.proto

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

type UpdateProductRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProductId          int64    `protobuf:"varint,1,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	ProductName        *string  `protobuf:"bytes,2,opt,name=product_name,json=productName,proto3,oneof" json:"product_name,omitempty"`
	ProductDescription *string  `protobuf:"bytes,3,opt,name=product_description,json=productDescription,proto3,oneof" json:"product_description,omitempty"`
	ProductCode        *string  `protobuf:"bytes,4,opt,name=product_code,json=productCode,proto3,oneof" json:"product_code,omitempty"`
	Price              *int64   `protobuf:"varint,5,opt,name=price,proto3,oneof" json:"price,omitempty"`
	SalePrice          *string  `protobuf:"bytes,6,opt,name=sale_price,json=salePrice,proto3,oneof" json:"sale_price,omitempty"`
	Collection         *int64   `protobuf:"varint,7,opt,name=collection,proto3,oneof" json:"collection,omitempty"`
	Quantity           *int32   `protobuf:"varint,8,opt,name=quantity,proto3,oneof" json:"quantity,omitempty"`
	Color              []string `protobuf:"bytes,9,rep,name=color,proto3" json:"color,omitempty"`
	Size               []string `protobuf:"bytes,10,rep,name=size,proto3" json:"size,omitempty"`
	Status             *string  `protobuf:"bytes,11,opt,name=status,proto3,oneof" json:"status,omitempty"`
	MainImage          *string  `protobuf:"bytes,12,opt,name=main_image,json=mainImage,proto3,oneof" json:"main_image,omitempty"`
	OtherImage_1       *string  `protobuf:"bytes,13,opt,name=other_image_1,json=otherImage1,proto3,oneof" json:"other_image_1,omitempty"`
	OtherImage_2       *string  `protobuf:"bytes,14,opt,name=other_image_2,json=otherImage2,proto3,oneof" json:"other_image_2,omitempty"`
	OtherImage_3       *string  `protobuf:"bytes,15,opt,name=other_image_3,json=otherImage3,proto3,oneof" json:"other_image_3,omitempty"`
}

func (x *UpdateProductRequest) Reset() {
	*x = UpdateProductRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_update_product_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateProductRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateProductRequest) ProtoMessage() {}

func (x *UpdateProductRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_update_product_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateProductRequest.ProtoReflect.Descriptor instead.
func (*UpdateProductRequest) Descriptor() ([]byte, []int) {
	return file_rpc_update_product_proto_rawDescGZIP(), []int{0}
}

func (x *UpdateProductRequest) GetProductId() int64 {
	if x != nil {
		return x.ProductId
	}
	return 0
}

func (x *UpdateProductRequest) GetProductName() string {
	if x != nil && x.ProductName != nil {
		return *x.ProductName
	}
	return ""
}

func (x *UpdateProductRequest) GetProductDescription() string {
	if x != nil && x.ProductDescription != nil {
		return *x.ProductDescription
	}
	return ""
}

func (x *UpdateProductRequest) GetProductCode() string {
	if x != nil && x.ProductCode != nil {
		return *x.ProductCode
	}
	return ""
}

func (x *UpdateProductRequest) GetPrice() int64 {
	if x != nil && x.Price != nil {
		return *x.Price
	}
	return 0
}

func (x *UpdateProductRequest) GetSalePrice() string {
	if x != nil && x.SalePrice != nil {
		return *x.SalePrice
	}
	return ""
}

func (x *UpdateProductRequest) GetCollection() int64 {
	if x != nil && x.Collection != nil {
		return *x.Collection
	}
	return 0
}

func (x *UpdateProductRequest) GetQuantity() int32 {
	if x != nil && x.Quantity != nil {
		return *x.Quantity
	}
	return 0
}

func (x *UpdateProductRequest) GetColor() []string {
	if x != nil {
		return x.Color
	}
	return nil
}

func (x *UpdateProductRequest) GetSize() []string {
	if x != nil {
		return x.Size
	}
	return nil
}

func (x *UpdateProductRequest) GetStatus() string {
	if x != nil && x.Status != nil {
		return *x.Status
	}
	return ""
}

func (x *UpdateProductRequest) GetMainImage() string {
	if x != nil && x.MainImage != nil {
		return *x.MainImage
	}
	return ""
}

func (x *UpdateProductRequest) GetOtherImage_1() string {
	if x != nil && x.OtherImage_1 != nil {
		return *x.OtherImage_1
	}
	return ""
}

func (x *UpdateProductRequest) GetOtherImage_2() string {
	if x != nil && x.OtherImage_2 != nil {
		return *x.OtherImage_2
	}
	return ""
}

func (x *UpdateProductRequest) GetOtherImage_3() string {
	if x != nil && x.OtherImage_3 != nil {
		return *x.OtherImage_3
	}
	return ""
}

type UpdateProductResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Product *Product `protobuf:"bytes,1,opt,name=product,proto3" json:"product,omitempty"`
}

func (x *UpdateProductResponse) Reset() {
	*x = UpdateProductResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_update_product_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateProductResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateProductResponse) ProtoMessage() {}

func (x *UpdateProductResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_update_product_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateProductResponse.ProtoReflect.Descriptor instead.
func (*UpdateProductResponse) Descriptor() ([]byte, []int) {
	return file_rpc_update_product_proto_rawDescGZIP(), []int{1}
}

func (x *UpdateProductResponse) GetProduct() *Product {
	if x != nil {
		return x.Product
	}
	return nil
}

var File_rpc_update_product_proto protoreflect.FileDescriptor

var file_rpc_update_product_proto_rawDesc = []byte{
	0x0a, 0x18, 0x72, 0x70, 0x63, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x70, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x1a, 0x0d,
	0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xe5, 0x05,
	0x0a, 0x14, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63,
	0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x64,
	0x75, 0x63, 0x74, 0x49, 0x64, 0x12, 0x26, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0b, 0x70,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x34, 0x0a,
	0x13, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x01, 0x52, 0x12, 0x70, 0x72,
	0x6f, 0x64, 0x75, 0x63, 0x74, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x88, 0x01, 0x01, 0x12, 0x26, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x63,
	0x6f, 0x64, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x48, 0x02, 0x52, 0x0b, 0x70, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x88, 0x01, 0x01, 0x12, 0x19, 0x0a, 0x05, 0x70,
	0x72, 0x69, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x48, 0x03, 0x52, 0x05, 0x70, 0x72,
	0x69, 0x63, 0x65, 0x88, 0x01, 0x01, 0x12, 0x22, 0x0a, 0x0a, 0x73, 0x61, 0x6c, 0x65, 0x5f, 0x70,
	0x72, 0x69, 0x63, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x48, 0x04, 0x52, 0x09, 0x73, 0x61,
	0x6c, 0x65, 0x50, 0x72, 0x69, 0x63, 0x65, 0x88, 0x01, 0x01, 0x12, 0x23, 0x0a, 0x0a, 0x63, 0x6f,
	0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x48, 0x05,
	0x52, 0x0a, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x12,
	0x1f, 0x0a, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x05, 0x48, 0x06, 0x52, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x88, 0x01, 0x01,
	0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x18, 0x09, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x0a,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x1b, 0x0a, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x48, 0x07, 0x52, 0x06, 0x73, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x88, 0x01, 0x01, 0x12, 0x22, 0x0a, 0x0a, 0x6d, 0x61, 0x69, 0x6e, 0x5f,
	0x69, 0x6d, 0x61, 0x67, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x48, 0x08, 0x52, 0x09, 0x6d,
	0x61, 0x69, 0x6e, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x88, 0x01, 0x01, 0x12, 0x27, 0x0a, 0x0d, 0x6f,
	0x74, 0x68, 0x65, 0x72, 0x5f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x31, 0x18, 0x0d, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x09, 0x52, 0x0b, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x49, 0x6d, 0x61, 0x67, 0x65,
	0x31, 0x88, 0x01, 0x01, 0x12, 0x27, 0x0a, 0x0d, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x5f, 0x69, 0x6d,
	0x61, 0x67, 0x65, 0x5f, 0x32, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x09, 0x48, 0x0a, 0x52, 0x0b, 0x6f,
	0x74, 0x68, 0x65, 0x72, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x32, 0x88, 0x01, 0x01, 0x12, 0x27, 0x0a,
	0x0d, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x5f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x33, 0x18, 0x0f,
	0x20, 0x01, 0x28, 0x09, 0x48, 0x0b, 0x52, 0x0b, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x49, 0x6d, 0x61,
	0x67, 0x65, 0x33, 0x88, 0x01, 0x01, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x70, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x16, 0x0a, 0x14, 0x5f, 0x70, 0x72, 0x6f, 0x64,
	0x75, 0x63, 0x74, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x42,
	0x0f, 0x0a, 0x0d, 0x5f, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x63, 0x6f, 0x64, 0x65,
	0x42, 0x08, 0x0a, 0x06, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x73,
	0x61, 0x6c, 0x65, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x63, 0x6f,
	0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x0b, 0x0a, 0x09, 0x5f, 0x71, 0x75, 0x61,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x42, 0x09, 0x0a, 0x07, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x6d, 0x61, 0x69, 0x6e, 0x5f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x42,
	0x10, 0x0a, 0x0e, 0x5f, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x5f, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f,
	0x31, 0x42, 0x10, 0x0a, 0x0e, 0x5f, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x5f, 0x69, 0x6d, 0x61, 0x67,
	0x65, 0x5f, 0x32, 0x42, 0x10, 0x0a, 0x0e, 0x5f, 0x6f, 0x74, 0x68, 0x65, 0x72, 0x5f, 0x69, 0x6d,
	0x61, 0x67, 0x65, 0x5f, 0x33, 0x22, 0x3e, 0x0a, 0x15, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x50,
	0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x25,
	0x0a, 0x07, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0b, 0x2e, 0x70, 0x62, 0x2e, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x52, 0x07, 0x70, 0x72,
	0x6f, 0x64, 0x75, 0x63, 0x74, 0x42, 0x25, 0x5a, 0x23, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x65, 0x63, 0x68, 0x73, 0x63, 0x68, 0x6f, 0x6f, 0x6c, 0x2f, 0x73,
	0x69, 0x6d, 0x70, 0x6c, 0x65, 0x62, 0x61, 0x6e, 0x6b, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_update_product_proto_rawDescOnce sync.Once
	file_rpc_update_product_proto_rawDescData = file_rpc_update_product_proto_rawDesc
)

func file_rpc_update_product_proto_rawDescGZIP() []byte {
	file_rpc_update_product_proto_rawDescOnce.Do(func() {
		file_rpc_update_product_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_update_product_proto_rawDescData)
	})
	return file_rpc_update_product_proto_rawDescData
}

var file_rpc_update_product_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_rpc_update_product_proto_goTypes = []interface{}{
	(*UpdateProductRequest)(nil),  // 0: pb.UpdateProductRequest
	(*UpdateProductResponse)(nil), // 1: pb.UpdateProductResponse
	(*Product)(nil),               // 2: pb.Product
}
var file_rpc_update_product_proto_depIdxs = []int32{
	2, // 0: pb.UpdateProductResponse.product:type_name -> pb.Product
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_rpc_update_product_proto_init() }
func file_rpc_update_product_proto_init() {
	if File_rpc_update_product_proto != nil {
		return
	}
	file_product_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_rpc_update_product_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateProductRequest); i {
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
		file_rpc_update_product_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateProductResponse); i {
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
	file_rpc_update_product_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_rpc_update_product_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_rpc_update_product_proto_goTypes,
		DependencyIndexes: file_rpc_update_product_proto_depIdxs,
		MessageInfos:      file_rpc_update_product_proto_msgTypes,
	}.Build()
	File_rpc_update_product_proto = out.File
	file_rpc_update_product_proto_rawDesc = nil
	file_rpc_update_product_proto_goTypes = nil
	file_rpc_update_product_proto_depIdxs = nil
}
