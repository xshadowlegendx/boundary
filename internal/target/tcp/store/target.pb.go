// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        (unknown)
// source: controller/storage/target/tcp/store/v1/target.proto

package store

import (
	timestamp "github.com/hashicorp/boundary/internal/db/timestamp"
	_ "github.com/hashicorp/boundary/sdk/pbs/controller/protooptions"
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

type Target struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// public_id is used to access the tcp.Target via an API
	// @inject_tag: gorm:"primary_key"
	PublicId string `protobuf:"bytes,10,opt,name=public_id,json=publicId,proto3" json:"public_id,omitempty" gorm:"primary_key"`
	// project id for the tcp.Target
	// @inject_tag: `gorm:"default:null"`
	ProjectId string `protobuf:"bytes,20,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty" gorm:"default:null"`
	// name is the optional friendly name used to
	// access the tcp.Target via an API
	// @inject_tag: `gorm:"default:null"`
	Name string `protobuf:"bytes,30,opt,name=name,proto3" json:"name,omitempty" gorm:"default:null"`
	// description of the tcp.Target
	// @inject_tag: `gorm:"default:null"`
	Description string `protobuf:"bytes,40,opt,name=description,proto3" json:"description,omitempty" gorm:"default:null"`
	// create_time from the RDBMS
	// @inject_tag: `gorm:"default:current_timestamp"`
	CreateTime *timestamp.Timestamp `protobuf:"bytes,50,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty" gorm:"default:current_timestamp"`
	// update_time from the RDBMS
	// @inject_tag: `gorm:"default:current_timestamp"`
	UpdateTime *timestamp.Timestamp `protobuf:"bytes,60,opt,name=update_time,json=updateTime,proto3" json:"update_time,omitempty" gorm:"default:current_timestamp"`
	// version allows optimistic locking of the tcp.Target when modifying the
	// tcp.Target
	// @inject_tag: `gorm:"default:null"`
	Version uint32 `protobuf:"varint,70,opt,name=version,proto3" json:"version,omitempty" gorm:"default:null"`
	// default port of the tcp.Target
	// @inject_tag: `gorm:"default:null"`
	DefaultPort uint32 `protobuf:"varint,80,opt,name=default_port,json=defaultPort,proto3" json:"default_port,omitempty" gorm:"default:null"`
	// Maximum total lifetime of a created session, in seconds
	// @inject_tag: `gorm:"default:null"`
	SessionMaxSeconds uint32 `protobuf:"varint,100,opt,name=session_max_seconds,json=sessionMaxSeconds,proto3" json:"session_max_seconds,omitempty" gorm:"default:null"`
	// Maximum number of connections in a session
	// @inject_tag: `gorm:"default:null"`
	SessionConnectionLimit int32 `protobuf:"varint,110,opt,name=session_connection_limit,json=sessionConnectionLimit,proto3" json:"session_connection_limit,omitempty" gorm:"default:null"`
	// A boolean expression that allows filtering the workers that can handle a session
	// @inject_tag: `gorm:"default:null"`
	WorkerFilter string `protobuf:"bytes,120,opt,name=worker_filter,json=workerFilter,proto3" json:"worker_filter,omitempty" gorm:"default:null"`
}

func (x *Target) Reset() {
	*x = Target{}
	if protoimpl.UnsafeEnabled {
		mi := &file_controller_storage_target_tcp_store_v1_target_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Target) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Target) ProtoMessage() {}

func (x *Target) ProtoReflect() protoreflect.Message {
	mi := &file_controller_storage_target_tcp_store_v1_target_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Target.ProtoReflect.Descriptor instead.
func (*Target) Descriptor() ([]byte, []int) {
	return file_controller_storage_target_tcp_store_v1_target_proto_rawDescGZIP(), []int{0}
}

func (x *Target) GetPublicId() string {
	if x != nil {
		return x.PublicId
	}
	return ""
}

func (x *Target) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

func (x *Target) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Target) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Target) GetCreateTime() *timestamp.Timestamp {
	if x != nil {
		return x.CreateTime
	}
	return nil
}

func (x *Target) GetUpdateTime() *timestamp.Timestamp {
	if x != nil {
		return x.UpdateTime
	}
	return nil
}

func (x *Target) GetVersion() uint32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *Target) GetDefaultPort() uint32 {
	if x != nil {
		return x.DefaultPort
	}
	return 0
}

func (x *Target) GetSessionMaxSeconds() uint32 {
	if x != nil {
		return x.SessionMaxSeconds
	}
	return 0
}

func (x *Target) GetSessionConnectionLimit() int32 {
	if x != nil {
		return x.SessionConnectionLimit
	}
	return 0
}

func (x *Target) GetWorkerFilter() string {
	if x != nil {
		return x.WorkerFilter
	}
	return ""
}

var File_controller_storage_target_tcp_store_v1_target_proto protoreflect.FileDescriptor

var file_controller_storage_target_tcp_store_v1_target_proto_rawDesc = []byte{
	0x0a, 0x33, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2f, 0x73, 0x74, 0x6f,
	0x72, 0x61, 0x67, 0x65, 0x2f, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x2f, 0x74, 0x63, 0x70, 0x2f,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x26, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65,
	0x72, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74,
	0x2e, 0x74, 0x63, 0x70, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x2a, 0x63,
	0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2f, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d,
	0x5f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x6f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2f, 0x63, 0x6f, 0x6e, 0x74, 0x72,
	0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc7, 0x05, 0x0a, 0x06, 0x54,
	0x61, 0x72, 0x67, 0x65, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f,
	0x69, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63,
	0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x14, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x49,
	0x64, 0x12, 0x24, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x1e, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x10, 0xc2, 0xdd, 0x29, 0x0c, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x40, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x28, 0x20, 0x01, 0x28, 0x09, 0x42, 0x1e, 0xc2, 0xdd,
	0x29, 0x1a, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0b, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x4b, 0x0a, 0x0b, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x32, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2a,
	0x2e, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x73, 0x74, 0x6f, 0x72,
	0x61, 0x67, 0x65, 0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x76, 0x31,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x4b, 0x0a, 0x0b, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x3c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x63, 0x6f,
	0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65,
	0x2e, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54,
	0x69, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x46,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x4d, 0x0a,
	0x0c, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x50, 0x20,
	0x01, 0x28, 0x0d, 0x42, 0x2a, 0xc2, 0xdd, 0x29, 0x26, 0x0a, 0x0b, 0x44, 0x65, 0x66, 0x61, 0x75,
	0x6c, 0x74, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x17, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74,
	0x65, 0x73, 0x2e, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x70, 0x6f, 0x72, 0x74, 0x52,
	0x0b, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x50, 0x6f, 0x72, 0x74, 0x12, 0x5c, 0x0a, 0x13,
	0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x6d, 0x61, 0x78, 0x5f, 0x73, 0x65, 0x63, 0x6f,
	0x6e, 0x64, 0x73, 0x18, 0x64, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x2c, 0xc2, 0xdd, 0x29, 0x28, 0x0a,
	0x11, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x4d, 0x61, 0x78, 0x53, 0x65, 0x63, 0x6f, 0x6e,
	0x64, 0x73, 0x12, 0x13, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x6d, 0x61, 0x78, 0x5f,
	0x73, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x52, 0x11, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x4d, 0x61, 0x78, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x12, 0x70, 0x0a, 0x18, 0x73, 0x65,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x6e, 0x20, 0x01, 0x28, 0x05, 0x42, 0x36, 0xc2, 0xdd,
	0x29, 0x32, 0x0a, 0x16, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x18, 0x73, 0x65, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x5f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6c,
	0x69, 0x6d, 0x69, 0x74, 0x52, 0x16, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6e,
	0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x46, 0x0a, 0x0d,
	0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x5f, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x78, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x21, 0xc2, 0xdd, 0x29, 0x1d, 0x0a, 0x0c, 0x57, 0x6f, 0x72, 0x6b, 0x65,
	0x72, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x0d, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x5f,
	0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x52, 0x0c, 0x77, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x46, 0x69,
	0x6c, 0x74, 0x65, 0x72, 0x42, 0x3f, 0x5a, 0x3d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x68, 0x61, 0x73, 0x68, 0x69, 0x63, 0x6f, 0x72, 0x70, 0x2f, 0x62, 0x6f, 0x75,
	0x6e, 0x64, 0x61, 0x72, 0x79, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x74,
	0x61, 0x72, 0x67, 0x65, 0x74, 0x2f, 0x74, 0x63, 0x70, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x3b,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_controller_storage_target_tcp_store_v1_target_proto_rawDescOnce sync.Once
	file_controller_storage_target_tcp_store_v1_target_proto_rawDescData = file_controller_storage_target_tcp_store_v1_target_proto_rawDesc
)

func file_controller_storage_target_tcp_store_v1_target_proto_rawDescGZIP() []byte {
	file_controller_storage_target_tcp_store_v1_target_proto_rawDescOnce.Do(func() {
		file_controller_storage_target_tcp_store_v1_target_proto_rawDescData = protoimpl.X.CompressGZIP(file_controller_storage_target_tcp_store_v1_target_proto_rawDescData)
	})
	return file_controller_storage_target_tcp_store_v1_target_proto_rawDescData
}

var file_controller_storage_target_tcp_store_v1_target_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_controller_storage_target_tcp_store_v1_target_proto_goTypes = []interface{}{
	(*Target)(nil),              // 0: controller.storage.target.tcp.store.v1.Target
	(*timestamp.Timestamp)(nil), // 1: controller.storage.timestamp.v1.Timestamp
}
var file_controller_storage_target_tcp_store_v1_target_proto_depIdxs = []int32{
	1, // 0: controller.storage.target.tcp.store.v1.Target.create_time:type_name -> controller.storage.timestamp.v1.Timestamp
	1, // 1: controller.storage.target.tcp.store.v1.Target.update_time:type_name -> controller.storage.timestamp.v1.Timestamp
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_controller_storage_target_tcp_store_v1_target_proto_init() }
func file_controller_storage_target_tcp_store_v1_target_proto_init() {
	if File_controller_storage_target_tcp_store_v1_target_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_controller_storage_target_tcp_store_v1_target_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Target); i {
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
			RawDescriptor: file_controller_storage_target_tcp_store_v1_target_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_controller_storage_target_tcp_store_v1_target_proto_goTypes,
		DependencyIndexes: file_controller_storage_target_tcp_store_v1_target_proto_depIdxs,
		MessageInfos:      file_controller_storage_target_tcp_store_v1_target_proto_msgTypes,
	}.Build()
	File_controller_storage_target_tcp_store_v1_target_proto = out.File
	file_controller_storage_target_tcp_store_v1_target_proto_rawDesc = nil
	file_controller_storage_target_tcp_store_v1_target_proto_goTypes = nil
	file_controller_storage_target_tcp_store_v1_target_proto_depIdxs = nil
}
