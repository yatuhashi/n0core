// Code generated by protoc-gen-go. DO NOT EDIT.
// source: v0/component.proto

package pn0stack // import "github.com/n0stack/proto.go/v0"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Pendingおよび存在しない場合はメッセージ自体が追加されないので問題ない
type Component_ComponentState int32

const (
	Component_OK   Component_ComponentState = 0
	Component_FAIL Component_ComponentState = 1
)

var Component_ComponentState_name = map[int32]string{
	0: "OK",
	1: "FAIL",
}
var Component_ComponentState_value = map[string]int32{
	"OK":   0,
	"FAIL": 1,
}

func (x Component_ComponentState) String() string {
	return proto.EnumName(Component_ComponentState_name, int32(x))
}
func (Component_ComponentState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_component_30a356b0e71ef218, []int{0, 0}
}

// ComponentとはモデルをWatchすることで機能を提供するアドオンのようなもの
// Modelにつき０個以上のComponentが関連づく
// すべてのComponentがOKの場合、正常に動作していると見ることができる
type Component struct {
	Service              string                   `protobuf:"bytes,1,opt,name=service,proto3" json:"service,omitempty"`
	Annotations          map[string]string        `protobuf:"bytes,2,rep,name=annotations,proto3" json:"annotations,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	State                Component_ComponentState `protobuf:"varint,3,opt,name=state,proto3,enum=n0stack.Component_ComponentState" json:"state,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *Component) Reset()         { *m = Component{} }
func (m *Component) String() string { return proto.CompactTextString(m) }
func (*Component) ProtoMessage()    {}
func (*Component) Descriptor() ([]byte, []int) {
	return fileDescriptor_component_30a356b0e71ef218, []int{0}
}
func (m *Component) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Component.Unmarshal(m, b)
}
func (m *Component) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Component.Marshal(b, m, deterministic)
}
func (dst *Component) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Component.Merge(dst, src)
}
func (m *Component) XXX_Size() int {
	return xxx_messageInfo_Component.Size(m)
}
func (m *Component) XXX_DiscardUnknown() {
	xxx_messageInfo_Component.DiscardUnknown(m)
}

var xxx_messageInfo_Component proto.InternalMessageInfo

func (m *Component) GetService() string {
	if m != nil {
		return m.Service
	}
	return ""
}

func (m *Component) GetAnnotations() map[string]string {
	if m != nil {
		return m.Annotations
	}
	return nil
}

func (m *Component) GetState() Component_ComponentState {
	if m != nil {
		return m.State
	}
	return Component_OK
}

func init() {
	proto.RegisterType((*Component)(nil), "n0stack.Component")
	proto.RegisterMapType((map[string]string)(nil), "n0stack.Component.AnnotationsEntry")
	proto.RegisterEnum("n0stack.Component_ComponentState", Component_ComponentState_name, Component_ComponentState_value)
}

func init() { proto.RegisterFile("v0/component.proto", fileDescriptor_component_30a356b0e71ef218) }

var fileDescriptor_component_30a356b0e71ef218 = []byte{
	// 245 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x12, 0x2a, 0x33, 0xd0, 0x4f,
	0xce, 0xcf, 0x2d, 0xc8, 0xcf, 0x4b, 0xcd, 0x2b, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62,
	0xcf, 0x33, 0x28, 0x2e, 0x49, 0x4c, 0xce, 0x56, 0xea, 0x64, 0xe2, 0xe2, 0x74, 0x86, 0x49, 0x0a,
	0x49, 0x70, 0xb1, 0x17, 0xa7, 0x16, 0x95, 0x65, 0x26, 0xa7, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0x70,
	0x06, 0xc1, 0xb8, 0x42, 0xae, 0x5c, 0xdc, 0x89, 0x79, 0x79, 0xf9, 0x25, 0x89, 0x25, 0x99, 0xf9,
	0x79, 0xc5, 0x12, 0x4c, 0x0a, 0xcc, 0x1a, 0xdc, 0x46, 0xca, 0x7a, 0x50, 0x63, 0xf4, 0xe0, 0x46,
	0xe8, 0x39, 0x22, 0x54, 0xb9, 0xe6, 0x95, 0x14, 0x55, 0x06, 0x21, 0xeb, 0x13, 0x32, 0xe7, 0x62,
	0x2d, 0x2e, 0x49, 0x2c, 0x49, 0x95, 0x60, 0x56, 0x60, 0xd4, 0xe0, 0x33, 0x52, 0xc4, 0x62, 0x00,
	0x9c, 0x15, 0x0c, 0x52, 0x18, 0x04, 0x51, 0x2f, 0x65, 0xc7, 0x25, 0x80, 0x6e, 0xb2, 0x90, 0x00,
	0x17, 0x73, 0x76, 0x6a, 0x25, 0xd4, 0xa5, 0x20, 0xa6, 0x90, 0x08, 0x17, 0x6b, 0x59, 0x62, 0x4e,
	0x69, 0xaa, 0x04, 0x13, 0x58, 0x0c, 0xc2, 0xb1, 0x62, 0xb2, 0x60, 0x54, 0x52, 0xe2, 0xe2, 0x43,
	0x35, 0x58, 0x88, 0x8d, 0x8b, 0xc9, 0xdf, 0x5b, 0x80, 0x41, 0x88, 0x83, 0x8b, 0xc5, 0xcd, 0xd1,
	0xd3, 0x47, 0x80, 0xd1, 0x49, 0x33, 0x4a, 0x3d, 0x3d, 0xb3, 0x24, 0xa3, 0x34, 0x49, 0x2f, 0x39,
	0x3f, 0x57, 0x1f, 0xea, 0x32, 0x7d, 0x70, 0x80, 0xe9, 0xa5, 0xe7, 0xeb, 0x97, 0x19, 0x58, 0x17,
	0x40, 0x05, 0x93, 0xd8, 0xc0, 0xa2, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x4a, 0xa6, 0xf1,
	0x8f, 0x5c, 0x01, 0x00, 0x00,
}
