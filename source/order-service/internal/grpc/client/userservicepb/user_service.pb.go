// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.31.0
// source: user_service.proto

package userservicepb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GetAllUsersRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetAllUsersRequest) Reset() {
	*x = GetAllUsersRequest{}
	mi := &file_user_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetAllUsersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllUsersRequest) ProtoMessage() {}

func (x *GetAllUsersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllUsersRequest.ProtoReflect.Descriptor instead.
func (*GetAllUsersRequest) Descriptor() ([]byte, []int) {
	return file_user_service_proto_rawDescGZIP(), []int{0}
}

type GetUserByIdRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetUserByIdRequest) Reset() {
	*x = GetUserByIdRequest{}
	mi := &file_user_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserByIdRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserByIdRequest) ProtoMessage() {}

func (x *GetUserByIdRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserByIdRequest.ProtoReflect.Descriptor instead.
func (*GetUserByIdRequest) Descriptor() ([]byte, []int) {
	return file_user_service_proto_rawDescGZIP(), []int{1}
}

func (x *GetUserByIdRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetAllUsersResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Users         []*User                `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetAllUsersResponse) Reset() {
	*x = GetAllUsersResponse{}
	mi := &file_user_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetAllUsersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAllUsersResponse) ProtoMessage() {}

func (x *GetAllUsersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAllUsersResponse.ProtoReflect.Descriptor instead.
func (*GetAllUsersResponse) Descriptor() ([]byte, []int) {
	return file_user_service_proto_rawDescGZIP(), []int{2}
}

func (x *GetAllUsersResponse) GetUsers() []*User {
	if x != nil {
		return x.Users
	}
	return nil
}

type GetUserByIdResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	User          *User                  `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetUserByIdResponse) Reset() {
	*x = GetUserByIdResponse{}
	mi := &file_user_service_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserByIdResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserByIdResponse) ProtoMessage() {}

func (x *GetUserByIdResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_service_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserByIdResponse.ProtoReflect.Descriptor instead.
func (*GetUserByIdResponse) Descriptor() ([]byte, []int) {
	return file_user_service_proto_rawDescGZIP(), []int{3}
}

func (x *GetUserByIdResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

type User struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	FullName      string                 `protobuf:"bytes,2,opt,name=full_name,json=fullName,proto3" json:"full_name,omitempty"`
	Email         string                 `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	Username      string                 `protobuf:"bytes,4,opt,name=username,proto3" json:"username,omitempty"`
	Address       string                 `protobuf:"bytes,5,opt,name=address,proto3" json:"address,omitempty"`
	RoleName      string                 `protobuf:"bytes,6,opt,name=role_name,json=roleName,proto3" json:"role_name,omitempty"`
	CreatedAt     *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *User) Reset() {
	*x = User{}
	mi := &file_user_service_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_user_service_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_user_service_proto_rawDescGZIP(), []int{4}
}

func (x *User) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *User) GetFullName() string {
	if x != nil {
		return x.FullName
	}
	return ""
}

func (x *User) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *User) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *User) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *User) GetRoleName() string {
	if x != nil {
		return x.RoleName
	}
	return ""
}

func (x *User) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *User) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

var File_user_service_proto protoreflect.FileDescriptor

const file_user_service_proto_rawDesc = "" +
	"\n" +
	"\x12user_service.proto\x12\ruserservicepb\x1a\x1fgoogle/protobuf/timestamp.proto\"\x14\n" +
	"\x12GetAllUsersRequest\"$\n" +
	"\x12GetUserByIdRequest\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\"@\n" +
	"\x13GetAllUsersResponse\x12)\n" +
	"\x05users\x18\x01 \x03(\v2\x13.userservicepb.UserR\x05users\">\n" +
	"\x13GetUserByIdResponse\x12'\n" +
	"\x04user\x18\x01 \x01(\v2\x13.userservicepb.UserR\x04user\"\x92\x02\n" +
	"\x04User\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x1b\n" +
	"\tfull_name\x18\x02 \x01(\tR\bfullName\x12\x14\n" +
	"\x05email\x18\x03 \x01(\tR\x05email\x12\x1a\n" +
	"\busername\x18\x04 \x01(\tR\busername\x12\x18\n" +
	"\aaddress\x18\x05 \x01(\tR\aaddress\x12\x1b\n" +
	"\trole_name\x18\x06 \x01(\tR\broleName\x129\n" +
	"\n" +
	"created_at\x18\a \x01(\v2\x1a.google.protobuf.TimestampR\tcreatedAt\x129\n" +
	"\n" +
	"updated_at\x18\b \x01(\v2\x1a.google.protobuf.TimestampR\tupdatedAt2\xbd\x01\n" +
	"\x0fUserServiceGRPC\x12T\n" +
	"\vGetAllUsers\x12!.userservicepb.GetAllUsersRequest\x1a\".userservicepb.GetAllUsersResponse\x12T\n" +
	"\vGetUserById\x12!.userservicepb.GetUserByIdRequest\x1a\".userservicepb.GetUserByIdResponseB\x10Z\x0euserservicepb/b\x06proto3"

var (
	file_user_service_proto_rawDescOnce sync.Once
	file_user_service_proto_rawDescData []byte
)

func file_user_service_proto_rawDescGZIP() []byte {
	file_user_service_proto_rawDescOnce.Do(func() {
		file_user_service_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_user_service_proto_rawDesc), len(file_user_service_proto_rawDesc)))
	})
	return file_user_service_proto_rawDescData
}

var file_user_service_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_user_service_proto_goTypes = []any{
	(*GetAllUsersRequest)(nil),    // 0: userservicepb.GetAllUsersRequest
	(*GetUserByIdRequest)(nil),    // 1: userservicepb.GetUserByIdRequest
	(*GetAllUsersResponse)(nil),   // 2: userservicepb.GetAllUsersResponse
	(*GetUserByIdResponse)(nil),   // 3: userservicepb.GetUserByIdResponse
	(*User)(nil),                  // 4: userservicepb.User
	(*timestamppb.Timestamp)(nil), // 5: google.protobuf.Timestamp
}
var file_user_service_proto_depIdxs = []int32{
	4, // 0: userservicepb.GetAllUsersResponse.users:type_name -> userservicepb.User
	4, // 1: userservicepb.GetUserByIdResponse.user:type_name -> userservicepb.User
	5, // 2: userservicepb.User.created_at:type_name -> google.protobuf.Timestamp
	5, // 3: userservicepb.User.updated_at:type_name -> google.protobuf.Timestamp
	0, // 4: userservicepb.UserServiceGRPC.GetAllUsers:input_type -> userservicepb.GetAllUsersRequest
	1, // 5: userservicepb.UserServiceGRPC.GetUserById:input_type -> userservicepb.GetUserByIdRequest
	2, // 6: userservicepb.UserServiceGRPC.GetAllUsers:output_type -> userservicepb.GetAllUsersResponse
	3, // 7: userservicepb.UserServiceGRPC.GetUserById:output_type -> userservicepb.GetUserByIdResponse
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_user_service_proto_init() }
func file_user_service_proto_init() {
	if File_user_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_user_service_proto_rawDesc), len(file_user_service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_user_service_proto_goTypes,
		DependencyIndexes: file_user_service_proto_depIdxs,
		MessageInfos:      file_user_service_proto_msgTypes,
	}.Build()
	File_user_service_proto = out.File
	file_user_service_proto_goTypes = nil
	file_user_service_proto_depIdxs = nil
}
