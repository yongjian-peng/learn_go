// Copyright 2022 The gRPC Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: gcp/observability/internal/logging/logging.proto

package logging

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// List of event types
type GrpcLogRecord_EventType int32

const (
	// Unknown event type
	GrpcLogRecord_GRPC_CALL_UNKNOWN GrpcLogRecord_EventType = 0
	// Header sent from client to server
	GrpcLogRecord_GRPC_CALL_REQUEST_HEADER GrpcLogRecord_EventType = 1
	// Header sent from server to client
	GrpcLogRecord_GRPC_CALL_RESPONSE_HEADER GrpcLogRecord_EventType = 2
	// Message sent from client to server
	GrpcLogRecord_GRPC_CALL_REQUEST_MESSAGE GrpcLogRecord_EventType = 3
	// Message sent from server to client
	GrpcLogRecord_GRPC_CALL_RESPONSE_MESSAGE GrpcLogRecord_EventType = 4
	// Trailer indicates the end of the gRPC call
	GrpcLogRecord_GRPC_CALL_TRAILER GrpcLogRecord_EventType = 5
	// A signal that client is done sending
	GrpcLogRecord_GRPC_CALL_HALF_CLOSE GrpcLogRecord_EventType = 6
	// A signal that the rpc is canceled
	GrpcLogRecord_GRPC_CALL_CANCEL GrpcLogRecord_EventType = 7
)

// Enum value maps for GrpcLogRecord_EventType.
var (
	GrpcLogRecord_EventType_name = map[int32]string{
		0: "GRPC_CALL_UNKNOWN",
		1: "GRPC_CALL_REQUEST_HEADER",
		2: "GRPC_CALL_RESPONSE_HEADER",
		3: "GRPC_CALL_REQUEST_MESSAGE",
		4: "GRPC_CALL_RESPONSE_MESSAGE",
		5: "GRPC_CALL_TRAILER",
		6: "GRPC_CALL_HALF_CLOSE",
		7: "GRPC_CALL_CANCEL",
	}
	GrpcLogRecord_EventType_value = map[string]int32{
		"GRPC_CALL_UNKNOWN":          0,
		"GRPC_CALL_REQUEST_HEADER":   1,
		"GRPC_CALL_RESPONSE_HEADER":  2,
		"GRPC_CALL_REQUEST_MESSAGE":  3,
		"GRPC_CALL_RESPONSE_MESSAGE": 4,
		"GRPC_CALL_TRAILER":          5,
		"GRPC_CALL_HALF_CLOSE":       6,
		"GRPC_CALL_CANCEL":           7,
	}
)

func (x GrpcLogRecord_EventType) Enum() *GrpcLogRecord_EventType {
	p := new(GrpcLogRecord_EventType)
	*p = x
	return p
}

func (x GrpcLogRecord_EventType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (GrpcLogRecord_EventType) Descriptor() protoreflect.EnumDescriptor {
	return file_gcp_observability_internal_logging_logging_proto_enumTypes[0].Descriptor()
}

func (GrpcLogRecord_EventType) Type() protoreflect.EnumType {
	return &file_gcp_observability_internal_logging_logging_proto_enumTypes[0]
}

func (x GrpcLogRecord_EventType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use GrpcLogRecord_EventType.Descriptor instead.
func (GrpcLogRecord_EventType) EnumDescriptor() ([]byte, []int) {
	return file_gcp_observability_internal_logging_logging_proto_rawDescGZIP(), []int{0, 0}
}

// The entity that generates the log entry
type GrpcLogRecord_EventLogger int32

const (
	GrpcLogRecord_LOGGER_UNKNOWN GrpcLogRecord_EventLogger = 0
	GrpcLogRecord_LOGGER_CLIENT  GrpcLogRecord_EventLogger = 1
	GrpcLogRecord_LOGGER_SERVER  GrpcLogRecord_EventLogger = 2
)

// Enum value maps for GrpcLogRecord_EventLogger.
var (
	GrpcLogRecord_EventLogger_name = map[int32]string{
		0: "LOGGER_UNKNOWN",
		1: "LOGGER_CLIENT",
		2: "LOGGER_SERVER",
	}
	GrpcLogRecord_EventLogger_value = map[string]int32{
		"LOGGER_UNKNOWN": 0,
		"LOGGER_CLIENT":  1,
		"LOGGER_SERVER":  2,
	}
)

func (x GrpcLogRecord_EventLogger) Enum() *GrpcLogRecord_EventLogger {
	p := new(GrpcLogRecord_EventLogger)
	*p = x
	return p
}

func (x GrpcLogRecord_EventLogger) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (GrpcLogRecord_EventLogger) Descriptor() protoreflect.EnumDescriptor {
	return file_gcp_observability_internal_logging_logging_proto_enumTypes[1].Descriptor()
}

func (GrpcLogRecord_EventLogger) Type() protoreflect.EnumType {
	return &file_gcp_observability_internal_logging_logging_proto_enumTypes[1]
}

func (x GrpcLogRecord_EventLogger) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use GrpcLogRecord_EventLogger.Descriptor instead.
func (GrpcLogRecord_EventLogger) EnumDescriptor() ([]byte, []int) {
	return file_gcp_observability_internal_logging_logging_proto_rawDescGZIP(), []int{0, 1}
}

// The log severity level of the log entry
type GrpcLogRecord_LogLevel int32

const (
	GrpcLogRecord_LOG_LEVEL_UNKNOWN  GrpcLogRecord_LogLevel = 0
	GrpcLogRecord_LOG_LEVEL_TRACE    GrpcLogRecord_LogLevel = 1
	GrpcLogRecord_LOG_LEVEL_DEBUG    GrpcLogRecord_LogLevel = 2
	GrpcLogRecord_LOG_LEVEL_INFO     GrpcLogRecord_LogLevel = 3
	GrpcLogRecord_LOG_LEVEL_WARN     GrpcLogRecord_LogLevel = 4
	GrpcLogRecord_LOG_LEVEL_ERROR    GrpcLogRecord_LogLevel = 5
	GrpcLogRecord_LOG_LEVEL_CRITICAL GrpcLogRecord_LogLevel = 6
)

// Enum value maps for GrpcLogRecord_LogLevel.
var (
	GrpcLogRecord_LogLevel_name = map[int32]string{
		0: "LOG_LEVEL_UNKNOWN",
		1: "LOG_LEVEL_TRACE",
		2: "LOG_LEVEL_DEBUG",
		3: "LOG_LEVEL_INFO",
		4: "LOG_LEVEL_WARN",
		5: "LOG_LEVEL_ERROR",
		6: "LOG_LEVEL_CRITICAL",
	}
	GrpcLogRecord_LogLevel_value = map[string]int32{
		"LOG_LEVEL_UNKNOWN":  0,
		"LOG_LEVEL_TRACE":    1,
		"LOG_LEVEL_DEBUG":    2,
		"LOG_LEVEL_INFO":     3,
		"LOG_LEVEL_WARN":     4,
		"LOG_LEVEL_ERROR":    5,
		"LOG_LEVEL_CRITICAL": 6,
	}
)

func (x GrpcLogRecord_LogLevel) Enum() *GrpcLogRecord_LogLevel {
	p := new(GrpcLogRecord_LogLevel)
	*p = x
	return p
}

func (x GrpcLogRecord_LogLevel) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (GrpcLogRecord_LogLevel) Descriptor() protoreflect.EnumDescriptor {
	return file_gcp_observability_internal_logging_logging_proto_enumTypes[2].Descriptor()
}

func (GrpcLogRecord_LogLevel) Type() protoreflect.EnumType {
	return &file_gcp_observability_internal_logging_logging_proto_enumTypes[2]
}

func (x GrpcLogRecord_LogLevel) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use GrpcLogRecord_LogLevel.Descriptor instead.
func (GrpcLogRecord_LogLevel) EnumDescriptor() ([]byte, []int) {
	return file_gcp_observability_internal_logging_logging_proto_rawDescGZIP(), []int{0, 2}
}

type GrpcLogRecord_Address_Type int32

const (
	GrpcLogRecord_Address_TYPE_UNKNOWN GrpcLogRecord_Address_Type = 0
	GrpcLogRecord_Address_TYPE_IPV4    GrpcLogRecord_Address_Type = 1 // in 1.2.3.4 form
	GrpcLogRecord_Address_TYPE_IPV6    GrpcLogRecord_Address_Type = 2 // IPv6 canonical form (RFC5952 section 4)
	GrpcLogRecord_Address_TYPE_UNIX    GrpcLogRecord_Address_Type = 3 // UDS string
)

// Enum value maps for GrpcLogRecord_Address_Type.
var (
	GrpcLogRecord_Address_Type_name = map[int32]string{
		0: "TYPE_UNKNOWN",
		1: "TYPE_IPV4",
		2: "TYPE_IPV6",
		3: "TYPE_UNIX",
	}
	GrpcLogRecord_Address_Type_value = map[string]int32{
		"TYPE_UNKNOWN": 0,
		"TYPE_IPV4":    1,
		"TYPE_IPV6":    2,
		"TYPE_UNIX":    3,
	}
)

func (x GrpcLogRecord_Address_Type) Enum() *GrpcLogRecord_Address_Type {
	p := new(GrpcLogRecord_Address_Type)
	*p = x
	return p
}

func (x GrpcLogRecord_Address_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (GrpcLogRecord_Address_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_gcp_observability_internal_logging_logging_proto_enumTypes[3].Descriptor()
}

func (GrpcLogRecord_Address_Type) Type() protoreflect.EnumType {
	return &file_gcp_observability_internal_logging_logging_proto_enumTypes[3]
}

func (x GrpcLogRecord_Address_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use GrpcLogRecord_Address_Type.Descriptor instead.
func (GrpcLogRecord_Address_Type) EnumDescriptor() ([]byte, []int) {
	return file_gcp_observability_internal_logging_logging_proto_rawDescGZIP(), []int{0, 2, 0}
}

type GrpcLogRecord struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The timestamp of the log event
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// Uniquely identifies a call. The value must not be 0 in order to disambiguate
	// from an unset value.
	// Each call may have several log entries. They will all have the same rpc_id.
	// Nothing is guaranteed about their value other than they are unique across
	// different RPCs in the same gRPC process.
	RpcId       string                    `protobuf:"bytes,2,opt,name=rpc_id,json=rpcId,proto3" json:"rpc_id,omitempty"`
	EventType   GrpcLogRecord_EventType   `protobuf:"varint,3,opt,name=event_type,json=eventType,proto3,enum=grpc.observability.logging.v1.GrpcLogRecord_EventType" json:"event_type,omitempty"`         // one of the above EventType enum
	EventLogger GrpcLogRecord_EventLogger `protobuf:"varint,4,opt,name=event_logger,json=eventLogger,proto3,enum=grpc.observability.logging.v1.GrpcLogRecord_EventLogger" json:"event_logger,omitempty"` // one of the above EventLogger enum
	// the name of the service
	ServiceName string `protobuf:"bytes,5,opt,name=service_name,json=serviceName,proto3" json:"service_name,omitempty"`
	// the name of the RPC method
	MethodName string                 `protobuf:"bytes,6,opt,name=method_name,json=methodName,proto3" json:"method_name,omitempty"`
	LogLevel   GrpcLogRecord_LogLevel `protobuf:"varint,7,opt,name=log_level,json=logLevel,proto3,enum=grpc.observability.logging.v1.GrpcLogRecord_LogLevel" json:"log_level,omitempty"` // one of the above LogLevel enum
	// Peer address information. On client side, peer is logged on server
	// header event or trailer event (if trailer-only). On server side, peer
	// is always logged on the client header event.
	PeerAddress *GrpcLogRecord_Address `protobuf:"bytes,8,opt,name=peer_address,json=peerAddress,proto3" json:"peer_address,omitempty"`
	// the RPC timeout value
	Timeout *durationpb.Duration `protobuf:"bytes,11,opt,name=timeout,proto3" json:"timeout,omitempty"`
	// A single process may be used to run multiple virtual servers with
	// different identities.
	// The authority is the name of such a server identify. It is typically a
	// portion of the URI in the form of <host> or <host>:<port>.
	Authority string `protobuf:"bytes,12,opt,name=authority,proto3" json:"authority,omitempty"`
	// Size of the message or metadata, depending on the event type,
	// regardless of whether the full message or metadata is being logged
	// (i.e. could be truncated or omitted).
	PayloadSize uint32 `protobuf:"varint,13,opt,name=payload_size,json=payloadSize,proto3" json:"payload_size,omitempty"`
	// true if message or metadata field is either truncated or omitted due
	// to config options
	PayloadTruncated bool `protobuf:"varint,14,opt,name=payload_truncated,json=payloadTruncated,proto3" json:"payload_truncated,omitempty"`
	// Used by header event or trailer event
	Metadata *GrpcLogRecord_Metadata `protobuf:"bytes,15,opt,name=metadata,proto3" json:"metadata,omitempty"`
	// The entry sequence ID for this call. The first message has a value of 1,
	// to disambiguate from an unset value. The purpose of this field is to
	// detect missing entries in environments where durability or ordering is
	// not guaranteed.
	SequenceId uint64 `protobuf:"varint,16,opt,name=sequence_id,json=sequenceId,proto3" json:"sequence_id,omitempty"`
	// Used by message event
	Message []byte `protobuf:"bytes,17,opt,name=message,proto3" json:"message,omitempty"`
	// The gRPC status code
	StatusCode uint32 `protobuf:"varint,18,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	// The gRPC status message
	StatusMessage string `protobuf:"bytes,19,opt,name=status_message,json=statusMessage,proto3" json:"status_message,omitempty"`
	// The value of the grpc-status-details-bin metadata key, if any.
	// This is always an encoded google.rpc.Status message
	StatusDetails []byte `protobuf:"bytes,20,opt,name=status_details,json=statusDetails,proto3" json:"status_details,omitempty"`
}

func (x *GrpcLogRecord) Reset() {
	*x = GrpcLogRecord{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gcp_observability_internal_logging_logging_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GrpcLogRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GrpcLogRecord) ProtoMessage() {}

func (x *GrpcLogRecord) ProtoReflect() protoreflect.Message {
	mi := &file_gcp_observability_internal_logging_logging_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GrpcLogRecord.ProtoReflect.Descriptor instead.
func (*GrpcLogRecord) Descriptor() ([]byte, []int) {
	return file_gcp_observability_internal_logging_logging_proto_rawDescGZIP(), []int{0}
}

func (x *GrpcLogRecord) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *GrpcLogRecord) GetRpcId() string {
	if x != nil {
		return x.RpcId
	}
	return ""
}

func (x *GrpcLogRecord) GetEventType() GrpcLogRecord_EventType {
	if x != nil {
		return x.EventType
	}
	return GrpcLogRecord_GRPC_CALL_UNKNOWN
}

func (x *GrpcLogRecord) GetEventLogger() GrpcLogRecord_EventLogger {
	if x != nil {
		return x.EventLogger
	}
	return GrpcLogRecord_LOGGER_UNKNOWN
}

func (x *GrpcLogRecord) GetServiceName() string {
	if x != nil {
		return x.ServiceName
	}
	return ""
}

func (x *GrpcLogRecord) GetMethodName() string {
	if x != nil {
		return x.MethodName
	}
	return ""
}

func (x *GrpcLogRecord) GetLogLevel() GrpcLogRecord_LogLevel {
	if x != nil {
		return x.LogLevel
	}
	return GrpcLogRecord_LOG_LEVEL_UNKNOWN
}

func (x *GrpcLogRecord) GetPeerAddress() *GrpcLogRecord_Address {
	if x != nil {
		return x.PeerAddress
	}
	return nil
}

func (x *GrpcLogRecord) GetTimeout() *durationpb.Duration {
	if x != nil {
		return x.Timeout
	}
	return nil
}

func (x *GrpcLogRecord) GetAuthority() string {
	if x != nil {
		return x.Authority
	}
	return ""
}

func (x *GrpcLogRecord) GetPayloadSize() uint32 {
	if x != nil {
		return x.PayloadSize
	}
	return 0
}

func (x *GrpcLogRecord) GetPayloadTruncated() bool {
	if x != nil {
		return x.PayloadTruncated
	}
	return false
}

func (x *GrpcLogRecord) GetMetadata() *GrpcLogRecord_Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *GrpcLogRecord) GetSequenceId() uint64 {
	if x != nil {
		return x.SequenceId
	}
	return 0
}

func (x *GrpcLogRecord) GetMessage() []byte {
	if x != nil {
		return x.Message
	}
	return nil
}

func (x *GrpcLogRecord) GetStatusCode() uint32 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

func (x *GrpcLogRecord) GetStatusMessage() string {
	if x != nil {
		return x.StatusMessage
	}
	return ""
}

func (x *GrpcLogRecord) GetStatusDetails() []byte {
	if x != nil {
		return x.StatusDetails
	}
	return nil
}

// A list of metadata pairs
type GrpcLogRecord_Metadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Entry []*GrpcLogRecord_MetadataEntry `protobuf:"bytes,1,rep,name=entry,proto3" json:"entry,omitempty"`
}

func (x *GrpcLogRecord_Metadata) Reset() {
	*x = GrpcLogRecord_Metadata{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gcp_observability_internal_logging_logging_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GrpcLogRecord_Metadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GrpcLogRecord_Metadata) ProtoMessage() {}

func (x *GrpcLogRecord_Metadata) ProtoReflect() protoreflect.Message {
	mi := &file_gcp_observability_internal_logging_logging_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GrpcLogRecord_Metadata.ProtoReflect.Descriptor instead.
func (*GrpcLogRecord_Metadata) Descriptor() ([]byte, []int) {
	return file_gcp_observability_internal_logging_logging_proto_rawDescGZIP(), []int{0, 0}
}

func (x *GrpcLogRecord_Metadata) GetEntry() []*GrpcLogRecord_MetadataEntry {
	if x != nil {
		return x.Entry
	}
	return nil
}

// One metadata key value pair
type GrpcLogRecord_MetadataEntry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *GrpcLogRecord_MetadataEntry) Reset() {
	*x = GrpcLogRecord_MetadataEntry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gcp_observability_internal_logging_logging_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GrpcLogRecord_MetadataEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GrpcLogRecord_MetadataEntry) ProtoMessage() {}

func (x *GrpcLogRecord_MetadataEntry) ProtoReflect() protoreflect.Message {
	mi := &file_gcp_observability_internal_logging_logging_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GrpcLogRecord_MetadataEntry.ProtoReflect.Descriptor instead.
func (*GrpcLogRecord_MetadataEntry) Descriptor() ([]byte, []int) {
	return file_gcp_observability_internal_logging_logging_proto_rawDescGZIP(), []int{0, 1}
}

func (x *GrpcLogRecord_MetadataEntry) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *GrpcLogRecord_MetadataEntry) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

// Address information
type GrpcLogRecord_Address struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type    GrpcLogRecord_Address_Type `protobuf:"varint,1,opt,name=type,proto3,enum=grpc.observability.logging.v1.GrpcLogRecord_Address_Type" json:"type,omitempty"`
	Address string                     `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty"`
	// only for TYPE_IPV4 and TYPE_IPV6
	IpPort uint32 `protobuf:"varint,3,opt,name=ip_port,json=ipPort,proto3" json:"ip_port,omitempty"`
}

func (x *GrpcLogRecord_Address) Reset() {
	*x = GrpcLogRecord_Address{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gcp_observability_internal_logging_logging_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GrpcLogRecord_Address) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GrpcLogRecord_Address) ProtoMessage() {}

func (x *GrpcLogRecord_Address) ProtoReflect() protoreflect.Message {
	mi := &file_gcp_observability_internal_logging_logging_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GrpcLogRecord_Address.ProtoReflect.Descriptor instead.
func (*GrpcLogRecord_Address) Descriptor() ([]byte, []int) {
	return file_gcp_observability_internal_logging_logging_proto_rawDescGZIP(), []int{0, 2}
}

func (x *GrpcLogRecord_Address) GetType() GrpcLogRecord_Address_Type {
	if x != nil {
		return x.Type
	}
	return GrpcLogRecord_Address_TYPE_UNKNOWN
}

func (x *GrpcLogRecord_Address) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *GrpcLogRecord_Address) GetIpPort() uint32 {
	if x != nil {
		return x.IpPort
	}
	return 0
}

var File_gcp_observability_internal_logging_logging_proto protoreflect.FileDescriptor

var file_gcp_observability_internal_logging_logging_proto_rawDesc = []byte{
	0x0a, 0x30, 0x67, 0x63, 0x70, 0x2f, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x62, 0x69, 0x6c,
	0x69, 0x74, 0x79, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x6c, 0x6f, 0x67,
	0x67, 0x69, 0x6e, 0x67, 0x2f, 0x6c, 0x6f, 0x67, 0x67, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x1d, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61,
	0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x2e, 0x6c, 0x6f, 0x67, 0x67, 0x69, 0x6e, 0x67, 0x2e, 0x76,
	0x31, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xe5, 0x0d, 0x0a, 0x0d, 0x47, 0x72, 0x70, 0x63, 0x4c, 0x6f, 0x67, 0x52, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x15,
	0x0a, 0x06, 0x72, 0x70, 0x63, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x72, 0x70, 0x63, 0x49, 0x64, 0x12, 0x55, 0x0a, 0x0a, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x36, 0x2e, 0x67, 0x72, 0x70, 0x63,
	0x2e, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x2e, 0x6c,
	0x6f, 0x67, 0x67, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x72, 0x70, 0x63, 0x4c, 0x6f,
	0x67, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x5b, 0x0a, 0x0c,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x6c, 0x6f, 0x67, 0x67, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x38, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76,
	0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x2e, 0x6c, 0x6f, 0x67, 0x67, 0x69, 0x6e, 0x67, 0x2e,
	0x76, 0x31, 0x2e, 0x47, 0x72, 0x70, 0x63, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64,
	0x2e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x4c, 0x6f, 0x67, 0x67, 0x65, 0x72, 0x52, 0x0b, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x4c, 0x6f, 0x67, 0x67, 0x65, 0x72, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x0b,
	0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x52, 0x0a,
	0x09, 0x6c, 0x6f, 0x67, 0x5f, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x35, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x62,
	0x69, 0x6c, 0x69, 0x74, 0x79, 0x2e, 0x6c, 0x6f, 0x67, 0x67, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31,
	0x2e, 0x47, 0x72, 0x70, 0x63, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x4c,
	0x6f, 0x67, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x52, 0x08, 0x6c, 0x6f, 0x67, 0x4c, 0x65, 0x76, 0x65,
	0x6c, 0x12, 0x57, 0x0a, 0x0c, 0x70, 0x65, 0x65, 0x72, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x34, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x6f,
	0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x2e, 0x6c, 0x6f, 0x67,
	0x67, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x72, 0x70, 0x63, 0x4c, 0x6f, 0x67, 0x52,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x52, 0x0b, 0x70,
	0x65, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x33, 0x0a, 0x07, 0x74, 0x69,
	0x6d, 0x65, 0x6f, 0x75, 0x74, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75,
	0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x74, 0x69, 0x6d, 0x65, 0x6f, 0x75, 0x74, 0x12,
	0x1c, 0x0a, 0x09, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x18, 0x0c, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x74, 0x79, 0x12, 0x21, 0x0a,
	0x0c, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x0d, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x0b, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x69, 0x7a, 0x65,
	0x12, 0x2b, 0x0a, 0x11, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x5f, 0x74, 0x72, 0x75, 0x6e,
	0x63, 0x61, 0x74, 0x65, 0x64, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x08, 0x52, 0x10, 0x70, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x54, 0x72, 0x75, 0x6e, 0x63, 0x61, 0x74, 0x65, 0x64, 0x12, 0x51, 0x0a,
	0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x35, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x62, 0x69,
	0x6c, 0x69, 0x74, 0x79, 0x2e, 0x6c, 0x6f, 0x67, 0x67, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x2e,
	0x47, 0x72, 0x70, 0x63, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18,
	0x10, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x49,
	0x64, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x11, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x12, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x25, 0x0a, 0x0e,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x13,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x12, 0x25, 0x0a, 0x0e, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x64, 0x65,
	0x74, 0x61, 0x69, 0x6c, 0x73, 0x18, 0x14, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0d, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x1a, 0x5c, 0x0a, 0x08, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x50, 0x0a, 0x05, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x3a, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x6f, 0x62, 0x73,
	0x65, 0x72, 0x76, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x2e, 0x6c, 0x6f, 0x67, 0x67, 0x69,
	0x6e, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x72, 0x70, 0x63, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x05, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x1a, 0x37, 0x0a, 0x0d, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x1a, 0xd2, 0x01, 0x0a, 0x07, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x4d, 0x0a,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x39, 0x2e, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79,
	0x2e, 0x6c, 0x6f, 0x67, 0x67, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x72, 0x70, 0x63,
	0x4c, 0x6f, 0x67, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x17, 0x0a, 0x07, 0x69, 0x70, 0x5f, 0x70, 0x6f, 0x72,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x69, 0x70, 0x50, 0x6f, 0x72, 0x74, 0x22,
	0x45, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x10, 0x0a, 0x0c, 0x54, 0x59, 0x50, 0x45, 0x5f,
	0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x54, 0x59, 0x50,
	0x45, 0x5f, 0x49, 0x50, 0x56, 0x34, 0x10, 0x01, 0x12, 0x0d, 0x0a, 0x09, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x49, 0x50, 0x56, 0x36, 0x10, 0x02, 0x12, 0x0d, 0x0a, 0x09, 0x54, 0x59, 0x50, 0x45, 0x5f,
	0x55, 0x4e, 0x49, 0x58, 0x10, 0x03, 0x22, 0xe5, 0x01, 0x0a, 0x09, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x15, 0x0a, 0x11, 0x47, 0x52, 0x50, 0x43, 0x5f, 0x43, 0x41, 0x4c,
	0x4c, 0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x1c, 0x0a, 0x18, 0x47,
	0x52, 0x50, 0x43, 0x5f, 0x43, 0x41, 0x4c, 0x4c, 0x5f, 0x52, 0x45, 0x51, 0x55, 0x45, 0x53, 0x54,
	0x5f, 0x48, 0x45, 0x41, 0x44, 0x45, 0x52, 0x10, 0x01, 0x12, 0x1d, 0x0a, 0x19, 0x47, 0x52, 0x50,
	0x43, 0x5f, 0x43, 0x41, 0x4c, 0x4c, 0x5f, 0x52, 0x45, 0x53, 0x50, 0x4f, 0x4e, 0x53, 0x45, 0x5f,
	0x48, 0x45, 0x41, 0x44, 0x45, 0x52, 0x10, 0x02, 0x12, 0x1d, 0x0a, 0x19, 0x47, 0x52, 0x50, 0x43,
	0x5f, 0x43, 0x41, 0x4c, 0x4c, 0x5f, 0x52, 0x45, 0x51, 0x55, 0x45, 0x53, 0x54, 0x5f, 0x4d, 0x45,
	0x53, 0x53, 0x41, 0x47, 0x45, 0x10, 0x03, 0x12, 0x1e, 0x0a, 0x1a, 0x47, 0x52, 0x50, 0x43, 0x5f,
	0x43, 0x41, 0x4c, 0x4c, 0x5f, 0x52, 0x45, 0x53, 0x50, 0x4f, 0x4e, 0x53, 0x45, 0x5f, 0x4d, 0x45,
	0x53, 0x53, 0x41, 0x47, 0x45, 0x10, 0x04, 0x12, 0x15, 0x0a, 0x11, 0x47, 0x52, 0x50, 0x43, 0x5f,
	0x43, 0x41, 0x4c, 0x4c, 0x5f, 0x54, 0x52, 0x41, 0x49, 0x4c, 0x45, 0x52, 0x10, 0x05, 0x12, 0x18,
	0x0a, 0x14, 0x47, 0x52, 0x50, 0x43, 0x5f, 0x43, 0x41, 0x4c, 0x4c, 0x5f, 0x48, 0x41, 0x4c, 0x46,
	0x5f, 0x43, 0x4c, 0x4f, 0x53, 0x45, 0x10, 0x06, 0x12, 0x14, 0x0a, 0x10, 0x47, 0x52, 0x50, 0x43,
	0x5f, 0x43, 0x41, 0x4c, 0x4c, 0x5f, 0x43, 0x41, 0x4e, 0x43, 0x45, 0x4c, 0x10, 0x07, 0x22, 0x47,
	0x0a, 0x0b, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x4c, 0x6f, 0x67, 0x67, 0x65, 0x72, 0x12, 0x12, 0x0a,
	0x0e, 0x4c, 0x4f, 0x47, 0x47, 0x45, 0x52, 0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10,
	0x00, 0x12, 0x11, 0x0a, 0x0d, 0x4c, 0x4f, 0x47, 0x47, 0x45, 0x52, 0x5f, 0x43, 0x4c, 0x49, 0x45,
	0x4e, 0x54, 0x10, 0x01, 0x12, 0x11, 0x0a, 0x0d, 0x4c, 0x4f, 0x47, 0x47, 0x45, 0x52, 0x5f, 0x53,
	0x45, 0x52, 0x56, 0x45, 0x52, 0x10, 0x02, 0x22, 0xa0, 0x01, 0x0a, 0x08, 0x4c, 0x6f, 0x67, 0x4c,
	0x65, 0x76, 0x65, 0x6c, 0x12, 0x15, 0x0a, 0x11, 0x4c, 0x4f, 0x47, 0x5f, 0x4c, 0x45, 0x56, 0x45,
	0x4c, 0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x13, 0x0a, 0x0f, 0x4c,
	0x4f, 0x47, 0x5f, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x54, 0x52, 0x41, 0x43, 0x45, 0x10, 0x01,
	0x12, 0x13, 0x0a, 0x0f, 0x4c, 0x4f, 0x47, 0x5f, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x44, 0x45,
	0x42, 0x55, 0x47, 0x10, 0x02, 0x12, 0x12, 0x0a, 0x0e, 0x4c, 0x4f, 0x47, 0x5f, 0x4c, 0x45, 0x56,
	0x45, 0x4c, 0x5f, 0x49, 0x4e, 0x46, 0x4f, 0x10, 0x03, 0x12, 0x12, 0x0a, 0x0e, 0x4c, 0x4f, 0x47,
	0x5f, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x57, 0x41, 0x52, 0x4e, 0x10, 0x04, 0x12, 0x13, 0x0a,
	0x0f, 0x4c, 0x4f, 0x47, 0x5f, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52,
	0x10, 0x05, 0x12, 0x16, 0x0a, 0x12, 0x4c, 0x4f, 0x47, 0x5f, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f,
	0x43, 0x52, 0x49, 0x54, 0x49, 0x43, 0x41, 0x4c, 0x10, 0x06, 0x42, 0x77, 0x0a, 0x1d, 0x69, 0x6f,
	0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x62, 0x69, 0x6c,
	0x69, 0x74, 0x79, 0x2e, 0x6c, 0x6f, 0x67, 0x67, 0x69, 0x6e, 0x67, 0x42, 0x19, 0x4f, 0x62, 0x73,
	0x65, 0x72, 0x76, 0x61, 0x62, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x4c, 0x6f, 0x67, 0x67, 0x69, 0x6e,
	0x67, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x39, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x2e, 0x6f, 0x72, 0x67, 0x2f, 0x67, 0x72, 0x70, 0x63,
	0x2f, 0x67, 0x63, 0x70, 0x2f, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x61, 0x62, 0x69, 0x6c, 0x69,
	0x74, 0x79, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x6c, 0x6f, 0x67, 0x67,
	0x69, 0x6e, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gcp_observability_internal_logging_logging_proto_rawDescOnce sync.Once
	file_gcp_observability_internal_logging_logging_proto_rawDescData = file_gcp_observability_internal_logging_logging_proto_rawDesc
)

func file_gcp_observability_internal_logging_logging_proto_rawDescGZIP() []byte {
	file_gcp_observability_internal_logging_logging_proto_rawDescOnce.Do(func() {
		file_gcp_observability_internal_logging_logging_proto_rawDescData = protoimpl.X.CompressGZIP(file_gcp_observability_internal_logging_logging_proto_rawDescData)
	})
	return file_gcp_observability_internal_logging_logging_proto_rawDescData
}

var file_gcp_observability_internal_logging_logging_proto_enumTypes = make([]protoimpl.EnumInfo, 4)
var file_gcp_observability_internal_logging_logging_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_gcp_observability_internal_logging_logging_proto_goTypes = []interface{}{
	(GrpcLogRecord_EventType)(0),        // 0: grpc.observability.logging.v1.GrpcLogRecord.EventType
	(GrpcLogRecord_EventLogger)(0),      // 1: grpc.observability.logging.v1.GrpcLogRecord.EventLogger
	(GrpcLogRecord_LogLevel)(0),         // 2: grpc.observability.logging.v1.GrpcLogRecord.LogLevel
	(GrpcLogRecord_Address_Type)(0),     // 3: grpc.observability.logging.v1.GrpcLogRecord.Address.Type
	(*GrpcLogRecord)(nil),               // 4: grpc.observability.logging.v1.GrpcLogRecord
	(*GrpcLogRecord_Metadata)(nil),      // 5: grpc.observability.logging.v1.GrpcLogRecord.Metadata
	(*GrpcLogRecord_MetadataEntry)(nil), // 6: grpc.observability.logging.v1.GrpcLogRecord.MetadataEntry
	(*GrpcLogRecord_Address)(nil),       // 7: grpc.observability.logging.v1.GrpcLogRecord.Address
	(*timestamppb.Timestamp)(nil),       // 8: google.protobuf.Timestamp
	(*durationpb.Duration)(nil),         // 9: google.protobuf.Duration
}
var file_gcp_observability_internal_logging_logging_proto_depIdxs = []int32{
	8, // 0: grpc.observability.logging.v1.GrpcLogRecord.timestamp:type_name -> google.protobuf.Timestamp
	0, // 1: grpc.observability.logging.v1.GrpcLogRecord.event_type:type_name -> grpc.observability.logging.v1.GrpcLogRecord.EventType
	1, // 2: grpc.observability.logging.v1.GrpcLogRecord.event_logger:type_name -> grpc.observability.logging.v1.GrpcLogRecord.EventLogger
	2, // 3: grpc.observability.logging.v1.GrpcLogRecord.log_level:type_name -> grpc.observability.logging.v1.GrpcLogRecord.LogLevel
	7, // 4: grpc.observability.logging.v1.GrpcLogRecord.peer_address:type_name -> grpc.observability.logging.v1.GrpcLogRecord.Address
	9, // 5: grpc.observability.logging.v1.GrpcLogRecord.timeout:type_name -> google.protobuf.Duration
	5, // 6: grpc.observability.logging.v1.GrpcLogRecord.metadata:type_name -> grpc.observability.logging.v1.GrpcLogRecord.Metadata
	6, // 7: grpc.observability.logging.v1.GrpcLogRecord.Metadata.entry:type_name -> grpc.observability.logging.v1.GrpcLogRecord.MetadataEntry
	3, // 8: grpc.observability.logging.v1.GrpcLogRecord.Address.type:type_name -> grpc.observability.logging.v1.GrpcLogRecord.Address.Type
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	9, // [9:9] is the sub-list for extension type_name
	9, // [9:9] is the sub-list for extension extendee
	0, // [0:9] is the sub-list for field type_name
}

func init() { file_gcp_observability_internal_logging_logging_proto_init() }
func file_gcp_observability_internal_logging_logging_proto_init() {
	if File_gcp_observability_internal_logging_logging_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gcp_observability_internal_logging_logging_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GrpcLogRecord); i {
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
		file_gcp_observability_internal_logging_logging_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GrpcLogRecord_Metadata); i {
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
		file_gcp_observability_internal_logging_logging_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GrpcLogRecord_MetadataEntry); i {
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
		file_gcp_observability_internal_logging_logging_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GrpcLogRecord_Address); i {
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
			RawDescriptor: file_gcp_observability_internal_logging_logging_proto_rawDesc,
			NumEnums:      4,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_gcp_observability_internal_logging_logging_proto_goTypes,
		DependencyIndexes: file_gcp_observability_internal_logging_logging_proto_depIdxs,
		EnumInfos:         file_gcp_observability_internal_logging_logging_proto_enumTypes,
		MessageInfos:      file_gcp_observability_internal_logging_logging_proto_msgTypes,
	}.Build()
	File_gcp_observability_internal_logging_logging_proto = out.File
	file_gcp_observability_internal_logging_logging_proto_rawDesc = nil
	file_gcp_observability_internal_logging_logging_proto_goTypes = nil
	file_gcp_observability_internal_logging_logging_proto_depIdxs = nil
}
