// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v3.6.1
// source: events/event_model.proto

package pb

import (
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type EventNew struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Title             string               `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	StartTime         *timestamp.Timestamp `protobuf:"bytes,2,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	EndTime           *timestamp.Timestamp `protobuf:"bytes,3,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
	Description       string               `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	UserId            string               `protobuf:"bytes,5,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	NotifyBeforeEvent int64                `protobuf:"varint,6,opt,name=notify_before_event,json=notifyBeforeEvent,proto3" json:"notify_before_event,omitempty"`
}

func (x *EventNew) Reset() {
	*x = EventNew{}
	if protoimpl.UnsafeEnabled {
		mi := &file_events_event_model_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventNew) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventNew) ProtoMessage() {}

func (x *EventNew) ProtoReflect() protoreflect.Message {
	mi := &file_events_event_model_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventNew.ProtoReflect.Descriptor instead.
func (*EventNew) Descriptor() ([]byte, []int) {
	return file_events_event_model_proto_rawDescGZIP(), []int{0}
}

func (x *EventNew) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *EventNew) GetStartTime() *timestamp.Timestamp {
	if x != nil {
		return x.StartTime
	}
	return nil
}

func (x *EventNew) GetEndTime() *timestamp.Timestamp {
	if x != nil {
		return x.EndTime
	}
	return nil
}

func (x *EventNew) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *EventNew) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *EventNew) GetNotifyBeforeEvent() int64 {
	if x != nil {
		return x.NotifyBeforeEvent
	}
	return 0
}

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                string               `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Title             string               `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	StartTime         *timestamp.Timestamp `protobuf:"bytes,3,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	EndTime           *timestamp.Timestamp `protobuf:"bytes,4,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
	Description       string               `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`
	UserId            string               `protobuf:"bytes,6,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	NotifyBeforeEvent int64                `protobuf:"varint,7,opt,name=notify_before_event,json=notifyBeforeEvent,proto3" json:"notify_before_event,omitempty"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_events_event_model_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_events_event_model_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_events_event_model_proto_rawDescGZIP(), []int{1}
}

func (x *Event) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Event) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Event) GetStartTime() *timestamp.Timestamp {
	if x != nil {
		return x.StartTime
	}
	return nil
}

func (x *Event) GetEndTime() *timestamp.Timestamp {
	if x != nil {
		return x.EndTime
	}
	return nil
}

func (x *Event) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Event) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *Event) GetNotifyBeforeEvent() int64 {
	if x != nil {
		return x.NotifyBeforeEvent
	}
	return 0
}

var File_events_event_model_proto protoreflect.FileDescriptor

var file_events_event_model_proto_rawDesc = []byte{
	0x0a, 0x18, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x6d,
	0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xfd, 0x01, 0x0a, 0x08, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x4e, 0x65, 0x77, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x39, 0x0a, 0x0a,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x73, 0x74,
	0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x35, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x74,
	0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x20,
	0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x2e, 0x0a, 0x13, 0x6e, 0x6f, 0x74,
	0x69, 0x66, 0x79, 0x5f, 0x62, 0x65, 0x66, 0x6f, 0x72, 0x65, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x11, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x42, 0x65,
	0x66, 0x6f, 0x72, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x8a, 0x02, 0x0a, 0x05, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x39, 0x0a, 0x0a, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x54, 0x69, 0x6d, 0x65, 0x12, 0x35, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64,
	0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x17, 0x0a,
	0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x2e, 0x0a, 0x13, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x79,
	0x5f, 0x62, 0x65, 0x66, 0x6f, 0x72, 0x65, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x11, 0x6e, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x42, 0x65, 0x66, 0x6f, 0x72,
	0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x42, 0x08, 0x5a, 0x06, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x62,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_events_event_model_proto_rawDescOnce sync.Once
	file_events_event_model_proto_rawDescData = file_events_event_model_proto_rawDesc
)

func file_events_event_model_proto_rawDescGZIP() []byte {
	file_events_event_model_proto_rawDescOnce.Do(func() {
		file_events_event_model_proto_rawDescData = protoimpl.X.CompressGZIP(file_events_event_model_proto_rawDescData)
	})
	return file_events_event_model_proto_rawDescData
}

var file_events_event_model_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_events_event_model_proto_goTypes = []any{
	(*EventNew)(nil),            // 0: proto.events.EventNew
	(*Event)(nil),               // 1: proto.events.Event
	(*timestamp.Timestamp)(nil), // 2: google.protobuf.Timestamp
}
var file_events_event_model_proto_depIdxs = []int32{
	2, // 0: proto.events.EventNew.start_time:type_name -> google.protobuf.Timestamp
	2, // 1: proto.events.EventNew.end_time:type_name -> google.protobuf.Timestamp
	2, // 2: proto.events.Event.start_time:type_name -> google.protobuf.Timestamp
	2, // 3: proto.events.Event.end_time:type_name -> google.protobuf.Timestamp
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_events_event_model_proto_init() }
func file_events_event_model_proto_init() {
	if File_events_event_model_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_events_event_model_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*EventNew); i {
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
		file_events_event_model_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*Event); i {
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
			RawDescriptor: file_events_event_model_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_events_event_model_proto_goTypes,
		DependencyIndexes: file_events_event_model_proto_depIdxs,
		MessageInfos:      file_events_event_model_proto_msgTypes,
	}.Build()
	File_events_event_model_proto = out.File
	file_events_event_model_proto_rawDesc = nil
	file_events_event_model_proto_goTypes = nil
	file_events_event_model_proto_depIdxs = nil
}
