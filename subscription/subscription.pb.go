// Code generated by protoc-gen-go. DO NOT EDIT.
// source: subscription.proto

package subscription // import "pinfire/subscription"

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

type ErrCode int32

const (
	ErrCode_Success ErrCode = 0
)

var ErrCode_name = map[int32]string{
	0: "Success",
}
var ErrCode_value = map[string]int32{
	"Success": 0,
}

func (x ErrCode) String() string {
	return proto.EnumName(ErrCode_name, int32(x))
}
func (ErrCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_subscription_a089406cf21c81e5, []int{0}
}

type CrawlType int32

const (
	CrawlType_CrawlTypeNil           CrawlType = 0
	CrawlType_CrawlTypeSubscription  CrawlType = 1
	CrawlType_CrawlTypeXpath         CrawlType = 2
	CrawlType_CrawlTypeFuzzyMatching CrawlType = 3
)

var CrawlType_name = map[int32]string{
	0: "CrawlTypeNil",
	1: "CrawlTypeSubscription",
	2: "CrawlTypeXpath",
	3: "CrawlTypeFuzzyMatching",
}
var CrawlType_value = map[string]int32{
	"CrawlTypeNil":           0,
	"CrawlTypeSubscription":  1,
	"CrawlTypeXpath":         2,
	"CrawlTypeFuzzyMatching": 3,
}

func (x CrawlType) String() string {
	return proto.EnumName(CrawlType_name, int32(x))
}
func (CrawlType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_subscription_a089406cf21c81e5, []int{1}
}

type ProxyNodeType int32

const (
	ProxyNodeType_ProxyNodeTypeNil    ProxyNodeType = 0
	ProxyNodeType_ProxyNodeTypeVmess  ProxyNodeType = 1
	ProxyNodeType_ProxyNodeTypeTrojan ProxyNodeType = 2
	ProxyNodeType_ProxyNodeTypeVless  ProxyNodeType = 3
)

var ProxyNodeType_name = map[int32]string{
	0: "ProxyNodeTypeNil",
	1: "ProxyNodeTypeVmess",
	2: "ProxyNodeTypeTrojan",
}
var ProxyNodeType_value = map[string]int32{
	"ProxyNodeTypeNil":    0,
	"ProxyNodeTypeVmess":  1,
	"ProxyNodeTypeTrojan": 2,
}

func (x ProxyNodeType) String() string {
	return proto.EnumName(ProxyNodeType_name, int32(x))
}
func (ProxyNodeType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_subscription_a089406cf21c81e5, []int{2}
}

type CrawlerConf struct {
	Id        uint64 `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	CreatedAt uint32 `protobuf:"varint,2,opt,name=created_at,json=createdAt" json:"created_at,omitempty"`
	UpdatedAt uint32 `protobuf:"varint,3,opt,name=updated_at,json=updatedAt" json:"updated_at,omitempty"`
	// @gorm: type:varchar(1000);unique_index:idx_crawl_url
	CrawlUrl  string `protobuf:"bytes,4,opt,name=crawl_url,json=crawlUrl" json:"crawl_url,omitempty" gorm:"type:varchar(1000);unique_index:idx_crawl_url"`
	CrawlType uint32 `protobuf:"varint,5,opt,name=crawl_type,json=crawlType" json:"crawl_type,omitempty"`
	// @gorm: type:json
	Rule *CrawlerConf_Rule `protobuf:"bytes,6,opt,name=rule" json:"rule,omitempty" gorm:"type:json"`
	// @grom: index: idx_next_crawl_at
	IsClosed bool `protobuf:"varint,7,opt,name=is_closed,json=isClosed" json:"is_closed,omitempty"`
	// @grom: index: idx_next_crawl_at
	NextAt   uint32 `protobuf:"varint,8,opt,name=next_at,json=nextAt" json:"next_at,omitempty"`
	Interval uint32 `protobuf:"varint,9,opt,name=interval" json:"interval,omitempty"`
	// @gorm: type:varchar(100)
	// @v: max=100
	Note                 string   `protobuf:"bytes,10,opt,name=note" json:"note,omitempty" gorm:"type:varchar(100)" validate:"max=100"`
	XXX_NoUnkeyedLiteral struct{} `json:"-" bson:"-" gorm:"-"`
	XXX_unrecognized     []byte   `json:"-" bson:"-" gorm:"-"`
	XXX_sizecache        int32    `json:"-" bson:"-" gorm:"-"`
}

func (m *CrawlerConf) Reset()         { *m = CrawlerConf{} }
func (m *CrawlerConf) String() string { return proto.CompactTextString(m) }
func (*CrawlerConf) ProtoMessage()    {}
func (*CrawlerConf) Descriptor() ([]byte, []int) {
	return fileDescriptor_subscription_a089406cf21c81e5, []int{0}
}
func (m *CrawlerConf) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CrawlerConf.Unmarshal(m, b)
}
func (m *CrawlerConf) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CrawlerConf.Marshal(b, m, deterministic)
}
func (dst *CrawlerConf) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CrawlerConf.Merge(dst, src)
}
func (m *CrawlerConf) XXX_Size() int {
	return xxx_messageInfo_CrawlerConf.Size(m)
}
func (m *CrawlerConf) XXX_DiscardUnknown() {
	xxx_messageInfo_CrawlerConf.DiscardUnknown(m)
}

var xxx_messageInfo_CrawlerConf proto.InternalMessageInfo

func (m *CrawlerConf) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *CrawlerConf) GetCreatedAt() uint32 {
	if m != nil {
		return m.CreatedAt
	}
	return 0
}

func (m *CrawlerConf) GetUpdatedAt() uint32 {
	if m != nil {
		return m.UpdatedAt
	}
	return 0
}

func (m *CrawlerConf) GetCrawlUrl() string {
	if m != nil {
		return m.CrawlUrl
	}
	return ""
}

func (m *CrawlerConf) GetCrawlType() uint32 {
	if m != nil {
		return m.CrawlType
	}
	return 0
}

func (m *CrawlerConf) GetRule() *CrawlerConf_Rule {
	if m != nil {
		return m.Rule
	}
	return nil
}

func (m *CrawlerConf) GetIsClosed() bool {
	if m != nil {
		return m.IsClosed
	}
	return false
}

func (m *CrawlerConf) GetNextAt() uint32 {
	if m != nil {
		return m.NextAt
	}
	return 0
}

func (m *CrawlerConf) GetInterval() uint32 {
	if m != nil {
		return m.Interval
	}
	return 0
}

func (m *CrawlerConf) GetNote() string {
	if m != nil {
		return m.Note
	}
	return ""
}

type CrawlerConf_Rule struct {
	UseProxy             bool     `protobuf:"varint,1,opt,name=use_proxy,json=useProxy" json:"use_proxy,omitempty"`
	Xpath                string   `protobuf:"bytes,2,opt,name=xpath" json:"xpath,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-" bson:"-" gorm:"-"`
	XXX_unrecognized     []byte   `json:"-" bson:"-" gorm:"-"`
	XXX_sizecache        int32    `json:"-" bson:"-" gorm:"-"`
}

func (m *CrawlerConf_Rule) Reset()         { *m = CrawlerConf_Rule{} }
func (m *CrawlerConf_Rule) String() string { return proto.CompactTextString(m) }
func (*CrawlerConf_Rule) ProtoMessage()    {}
func (*CrawlerConf_Rule) Descriptor() ([]byte, []int) {
	return fileDescriptor_subscription_a089406cf21c81e5, []int{0, 0}
}
func (m *CrawlerConf_Rule) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CrawlerConf_Rule.Unmarshal(m, b)
}
func (m *CrawlerConf_Rule) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CrawlerConf_Rule.Marshal(b, m, deterministic)
}
func (dst *CrawlerConf_Rule) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CrawlerConf_Rule.Merge(dst, src)
}
func (m *CrawlerConf_Rule) XXX_Size() int {
	return xxx_messageInfo_CrawlerConf_Rule.Size(m)
}
func (m *CrawlerConf_Rule) XXX_DiscardUnknown() {
	xxx_messageInfo_CrawlerConf_Rule.DiscardUnknown(m)
}

var xxx_messageInfo_CrawlerConf_Rule proto.InternalMessageInfo

func (m *CrawlerConf_Rule) GetUseProxy() bool {
	if m != nil {
		return m.UseProxy
	}
	return false
}

func (m *CrawlerConf_Rule) GetXpath() string {
	if m != nil {
		return m.Xpath
	}
	return ""
}

type ProxyNode struct {
	Id        uint64 `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	CreatedAt uint32 `protobuf:"varint,2,opt,name=created_at,json=createdAt" json:"created_at,omitempty"`
	UpdatedAt uint32 `protobuf:"varint,3,opt,name=updated_at,json=updatedAt" json:"updated_at,omitempty"`
	// HOST:PORT
	// @gorm: type:varchar(1000)
	// @v: max=1000
	Url           string `protobuf:"bytes,4,opt,name=url" json:"url,omitempty" gorm:"type:varchar(1000)" validate:"max=1000"`
	ProxyNodeType uint32 `protobuf:"varint,5,opt,name=proxy_node_type,json=proxyNodeType" json:"proxy_node_type,omitempty"`
	// @gorm: type:text
	NodeDetail        *ProxyNode_NodeDetail `protobuf:"bytes,6,opt,name=node_detail,json=nodeDetail" json:"node_detail,omitempty" gorm:"type:text"`
	ProxySpeed        float64               `protobuf:"fixed64,7,opt,name=proxy_speed,json=proxySpeed" json:"proxy_speed,omitempty"`
	ProxyNetworkDelay float64               `protobuf:"fixed64,8,opt,name=proxy_network_delay,json=proxyNetworkDelay" json:"proxy_network_delay,omitempty"`
	// @grom: index: idx_next_check_at
	NextCheckAt   uint32 `protobuf:"varint,9,opt,name=next_check_at,json=nextCheckAt" json:"next_check_at,omitempty"`
	CheckInterval uint32 `protobuf:"varint,10,opt,name=check_interval,json=checkInterval" json:"check_interval,omitempty"`
	CrawlId       uint64 `protobuf:"varint,11,opt,name=crawl_id,json=crawlId" json:"crawl_id,omitempty"`
	// @grom: index: idx_next_check_at
	DeathCount uint32 `protobuf:"varint,12,opt,name=death_count,json=deathCount" json:"death_count,omitempty"`
	// @grom: index: idx_next_check_at
	IsClose              bool     `protobuf:"varint,13,opt,name=is_close,json=isClose" json:"is_close,omitempty"`
	LastCrawlerAt        uint32   `protobuf:"varint,14,opt,name=last_crawler_at,json=lastCrawlerAt" json:"last_crawler_at,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-" bson:"-" gorm:"-"`
	XXX_unrecognized     []byte   `json:"-" bson:"-" gorm:"-"`
	XXX_sizecache        int32    `json:"-" bson:"-" gorm:"-"`
}

func (m *ProxyNode) Reset()         { *m = ProxyNode{} }
func (m *ProxyNode) String() string { return proto.CompactTextString(m) }
func (*ProxyNode) ProtoMessage()    {}
func (*ProxyNode) Descriptor() ([]byte, []int) {
	return fileDescriptor_subscription_a089406cf21c81e5, []int{1}
}
func (m *ProxyNode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProxyNode.Unmarshal(m, b)
}
func (m *ProxyNode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProxyNode.Marshal(b, m, deterministic)
}
func (dst *ProxyNode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProxyNode.Merge(dst, src)
}
func (m *ProxyNode) XXX_Size() int {
	return xxx_messageInfo_ProxyNode.Size(m)
}
func (m *ProxyNode) XXX_DiscardUnknown() {
	xxx_messageInfo_ProxyNode.DiscardUnknown(m)
}

var xxx_messageInfo_ProxyNode proto.InternalMessageInfo

func (m *ProxyNode) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *ProxyNode) GetCreatedAt() uint32 {
	if m != nil {
		return m.CreatedAt
	}
	return 0
}

func (m *ProxyNode) GetUpdatedAt() uint32 {
	if m != nil {
		return m.UpdatedAt
	}
	return 0
}

func (m *ProxyNode) GetUrl() string {
	if m != nil {
		return m.Url
	}
	return ""
}

func (m *ProxyNode) GetProxyNodeType() uint32 {
	if m != nil {
		return m.ProxyNodeType
	}
	return 0
}

func (m *ProxyNode) GetNodeDetail() *ProxyNode_NodeDetail {
	if m != nil {
		return m.NodeDetail
	}
	return nil
}

func (m *ProxyNode) GetProxySpeed() float64 {
	if m != nil {
		return m.ProxySpeed
	}
	return 0
}

func (m *ProxyNode) GetProxyNetworkDelay() float64 {
	if m != nil {
		return m.ProxyNetworkDelay
	}
	return 0
}

func (m *ProxyNode) GetNextCheckAt() uint32 {
	if m != nil {
		return m.NextCheckAt
	}
	return 0
}

func (m *ProxyNode) GetCheckInterval() uint32 {
	if m != nil {
		return m.CheckInterval
	}
	return 0
}

func (m *ProxyNode) GetCrawlId() uint64 {
	if m != nil {
		return m.CrawlId
	}
	return 0
}

func (m *ProxyNode) GetDeathCount() uint32 {
	if m != nil {
		return m.DeathCount
	}
	return 0
}

func (m *ProxyNode) GetIsClose() bool {
	if m != nil {
		return m.IsClose
	}
	return false
}

func (m *ProxyNode) GetLastCrawlerAt() uint32 {
	if m != nil {
		return m.LastCrawlerAt
	}
	return 0
}

type ProxyNode_VmessNode struct {
	Path                 string   `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	Tls                  string   `protobuf:"bytes,2,opt,name=tls" json:"tls,omitempty"`
	V                    string   `protobuf:"bytes,3,opt,name=v" json:"v,omitempty"`
	Aid                  string   `protobuf:"bytes,4,opt,name=aid" json:"aid,omitempty"`
	Net                  string   `protobuf:"bytes,5,opt,name=net" json:"net,omitempty"`
	Id                   string   `protobuf:"bytes,6,opt,name=id" json:"id,omitempty"`
	Type                 string   `protobuf:"bytes,7,opt,name=type" json:"type,omitempty"`
	Host                 string   `protobuf:"bytes,8,opt,name=host" json:"host,omitempty"`
	Ps                   string   `protobuf:"bytes,9,opt,name=ps" json:"ps,omitempty"`
	Add                  string   `protobuf:"bytes,10,opt,name=add" json:"add,omitempty"`
	Port                 string   `protobuf:"bytes,11,opt,name=port" json:"port,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-" bson:"-" gorm:"-"`
	XXX_unrecognized     []byte   `json:"-" bson:"-" gorm:"-"`
	XXX_sizecache        int32    `json:"-" bson:"-" gorm:"-"`
}

func (m *ProxyNode_VmessNode) Reset()         { *m = ProxyNode_VmessNode{} }
func (m *ProxyNode_VmessNode) String() string { return proto.CompactTextString(m) }
func (*ProxyNode_VmessNode) ProtoMessage()    {}
func (*ProxyNode_VmessNode) Descriptor() ([]byte, []int) {
	return fileDescriptor_subscription_a089406cf21c81e5, []int{1, 0}
}
func (m *ProxyNode_VmessNode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProxyNode_VmessNode.Unmarshal(m, b)
}
func (m *ProxyNode_VmessNode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProxyNode_VmessNode.Marshal(b, m, deterministic)
}
func (dst *ProxyNode_VmessNode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProxyNode_VmessNode.Merge(dst, src)
}
func (m *ProxyNode_VmessNode) XXX_Size() int {
	return xxx_messageInfo_ProxyNode_VmessNode.Size(m)
}
func (m *ProxyNode_VmessNode) XXX_DiscardUnknown() {
	xxx_messageInfo_ProxyNode_VmessNode.DiscardUnknown(m)
}

var xxx_messageInfo_ProxyNode_VmessNode proto.InternalMessageInfo

func (m *ProxyNode_VmessNode) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *ProxyNode_VmessNode) GetTls() string {
	if m != nil {
		return m.Tls
	}
	return ""
}

func (m *ProxyNode_VmessNode) GetV() string {
	if m != nil {
		return m.V
	}
	return ""
}

func (m *ProxyNode_VmessNode) GetAid() string {
	if m != nil {
		return m.Aid
	}
	return ""
}

func (m *ProxyNode_VmessNode) GetNet() string {
	if m != nil {
		return m.Net
	}
	return ""
}

func (m *ProxyNode_VmessNode) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ProxyNode_VmessNode) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *ProxyNode_VmessNode) GetHost() string {
	if m != nil {
		return m.Host
	}
	return ""
}

func (m *ProxyNode_VmessNode) GetPs() string {
	if m != nil {
		return m.Ps
	}
	return ""
}

func (m *ProxyNode_VmessNode) GetAdd() string {
	if m != nil {
		return m.Add
	}
	return ""
}

func (m *ProxyNode_VmessNode) GetPort() string {
	if m != nil {
		return m.Port
	}
	return ""
}

type ProxyNode_NodeDetail struct {
	Buf                  string               `protobuf:"bytes,1,opt,name=buf" json:"buf,omitempty"`
	VmessNode            *ProxyNode_VmessNode `protobuf:"bytes,2,opt,name=vmess_node,json=vmessNode" json:"vmess_node,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-" bson:"-" gorm:"-"`
	XXX_unrecognized     []byte               `json:"-" bson:"-" gorm:"-"`
	XXX_sizecache        int32                `json:"-" bson:"-" gorm:"-"`
}

func (m *ProxyNode_NodeDetail) Reset()         { *m = ProxyNode_NodeDetail{} }
func (m *ProxyNode_NodeDetail) String() string { return proto.CompactTextString(m) }
func (*ProxyNode_NodeDetail) ProtoMessage()    {}
func (*ProxyNode_NodeDetail) Descriptor() ([]byte, []int) {
	return fileDescriptor_subscription_a089406cf21c81e5, []int{1, 1}
}
func (m *ProxyNode_NodeDetail) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProxyNode_NodeDetail.Unmarshal(m, b)
}
func (m *ProxyNode_NodeDetail) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProxyNode_NodeDetail.Marshal(b, m, deterministic)
}
func (dst *ProxyNode_NodeDetail) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProxyNode_NodeDetail.Merge(dst, src)
}
func (m *ProxyNode_NodeDetail) XXX_Size() int {
	return xxx_messageInfo_ProxyNode_NodeDetail.Size(m)
}
func (m *ProxyNode_NodeDetail) XXX_DiscardUnknown() {
	xxx_messageInfo_ProxyNode_NodeDetail.DiscardUnknown(m)
}

var xxx_messageInfo_ProxyNode_NodeDetail proto.InternalMessageInfo

func (m *ProxyNode_NodeDetail) GetBuf() string {
	if m != nil {
		return m.Buf
	}
	return ""
}

func (m *ProxyNode_NodeDetail) GetVmessNode() *ProxyNode_VmessNode {
	if m != nil {
		return m.VmessNode
	}
	return nil
}

func init() {
	proto.RegisterType((*CrawlerConf)(nil), "subscription.CrawlerConf")
	proto.RegisterType((*CrawlerConf_Rule)(nil), "subscription.CrawlerConf.Rule")
	proto.RegisterType((*ProxyNode)(nil), "subscription.ProxyNode")
	proto.RegisterType((*ProxyNode_VmessNode)(nil), "subscription.ProxyNode.VmessNode")
	proto.RegisterType((*ProxyNode_NodeDetail)(nil), "subscription.ProxyNode.NodeDetail")
	proto.RegisterEnum("subscription.ErrCode", ErrCode_name, ErrCode_value)
	proto.RegisterEnum("subscription.CrawlType", CrawlType_name, CrawlType_value)
	proto.RegisterEnum("subscription.ProxyNodeType", ProxyNodeType_name, ProxyNodeType_value)
}

func init() { proto.RegisterFile("subscription.proto", fileDescriptor_subscription_a089406cf21c81e5) }

var fileDescriptor_subscription_a089406cf21c81e5 = []byte{
	// 750 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x54, 0xc1, 0x6e, 0xe3, 0x36,
	0x10, 0x8d, 0x1c, 0xc7, 0xb6, 0x46, 0xb6, 0xab, 0x32, 0xa9, 0xa3, 0xb8, 0x68, 0xeb, 0x1a, 0x68,
	0x61, 0xe4, 0xe0, 0x02, 0xee, 0xa9, 0xb7, 0xba, 0x4a, 0x0b, 0xe4, 0xd0, 0x20, 0x60, 0xd2, 0xa0,
	0xe8, 0x45, 0x55, 0x44, 0x66, 0xad, 0x8d, 0x56, 0x12, 0x48, 0xca, 0x89, 0xf3, 0x47, 0xfb, 0x39,
	0xfb, 0x1f, 0xfb, 0x11, 0x8b, 0x19, 0xca, 0x8a, 0x7d, 0xd8, 0xdb, 0x5e, 0x8c, 0x99, 0x37, 0x22,
	0x67, 0xe6, 0xbd, 0x67, 0x02, 0xd3, 0xd5, 0xbd, 0x4e, 0x54, 0x5a, 0x9a, 0xb4, 0xc8, 0xe7, 0xa5,
	0x2a, 0x4c, 0xc1, 0xfa, 0xbb, 0xd8, 0xf4, 0x63, 0x0b, 0xbc, 0x50, 0xc5, 0x4f, 0x99, 0x54, 0x61,
	0x91, 0x3f, 0xb0, 0x21, 0xb4, 0x52, 0x11, 0x38, 0x13, 0x67, 0xd6, 0xe6, 0xad, 0x54, 0xb0, 0xef,
	0x00, 0x12, 0x25, 0x63, 0x23, 0x45, 0x14, 0x9b, 0xa0, 0x35, 0x71, 0x66, 0x03, 0xee, 0xd6, 0xc8,
	0xd2, 0x60, 0xb9, 0x2a, 0xc5, 0xb6, 0x7c, 0x68, 0xcb, 0x35, 0xb2, 0x34, 0xec, 0x5b, 0x70, 0x13,
	0xbc, 0x3c, 0xaa, 0x54, 0x16, 0xb4, 0x27, 0xce, 0xcc, 0xe5, 0x3d, 0x02, 0xfe, 0x51, 0x99, 0xbd,
	0x1a, 0x8b, 0x66, 0x53, 0xca, 0xe0, 0x68, 0x7b, 0x75, 0xfc, 0x94, 0xdd, 0x6e, 0x4a, 0xc9, 0x16,
	0xd0, 0x56, 0x55, 0x26, 0x83, 0xce, 0xc4, 0x99, 0x79, 0x8b, 0xef, 0xe7, 0x7b, 0xab, 0xec, 0x8c,
	0x3c, 0xe7, 0x55, 0x26, 0x39, 0x7d, 0x8b, 0xfd, 0x52, 0x1d, 0x25, 0x59, 0xa1, 0xa5, 0x08, 0xba,
	0x13, 0x67, 0xd6, 0xe3, 0xbd, 0x54, 0x87, 0x94, 0xb3, 0x53, 0xe8, 0xe6, 0xf2, 0xd9, 0xe0, 0xa0,
	0x3d, 0x6a, 0xd6, 0xc1, 0x74, 0x69, 0xd8, 0x18, 0x7a, 0x69, 0x6e, 0xa4, 0x5a, 0xc7, 0x59, 0xe0,
	0x52, 0xa5, 0xc9, 0x19, 0x83, 0x76, 0x5e, 0x18, 0x19, 0x00, 0x0d, 0x4f, 0xf1, 0xf8, 0x37, 0x68,
	0xf3, 0xba, 0x5b, 0xa5, 0x65, 0x54, 0xaa, 0xe2, 0x79, 0x43, 0x94, 0xf5, 0x78, 0xaf, 0xd2, 0xf2,
	0x1a, 0x73, 0x76, 0x02, 0x47, 0xcf, 0x65, 0x6c, 0x56, 0xc4, 0x99, 0xcb, 0x6d, 0x32, 0x7d, 0xdf,
	0x01, 0x97, 0xea, 0x57, 0x85, 0x90, 0x5f, 0x98, 0x6c, 0x1f, 0x0e, 0x5f, 0x69, 0xc6, 0x90, 0xfd,
	0x0c, 0x5f, 0xd1, 0x70, 0x51, 0x5e, 0x08, 0xb9, 0x4b, 0xf3, 0xa0, 0xdc, 0xce, 0x40, 0x54, 0x87,
	0xe0, 0xd1, 0x17, 0x42, 0x9a, 0x38, 0xcd, 0x6a, 0xc6, 0xa7, 0xfb, 0x8c, 0x37, 0x53, 0xcf, 0xf1,
	0xe7, 0x82, 0xbe, 0xe4, 0x90, 0x37, 0x31, 0xfb, 0x01, 0x3c, 0xdb, 0x4c, 0x97, 0xb2, 0x66, 0xdf,
	0xe1, 0x40, 0xd0, 0x0d, 0x22, 0x6c, 0x0e, 0xc7, 0xf5, 0x34, 0xd2, 0x3c, 0x15, 0xea, 0x31, 0x12,
	0x32, 0x8b, 0x37, 0xa4, 0x85, 0xc3, 0xbf, 0xb6, 0x13, 0xd9, 0xca, 0x05, 0x16, 0xd8, 0x14, 0x06,
	0xa4, 0x57, 0xb2, 0x92, 0xc9, 0x23, 0x6e, 0x6c, 0xb5, 0xf1, 0x10, 0x0c, 0x11, 0x5b, 0x1a, 0xf6,
	0x13, 0x0c, 0x6d, 0xb9, 0x11, 0x10, 0xec, 0x82, 0x84, 0x5e, 0x6e, 0x55, 0x3c, 0x03, 0x6b, 0xbb,
	0x28, 0x15, 0x81, 0x47, 0x74, 0x77, 0x29, 0xbf, 0x14, 0x38, 0xb6, 0x90, 0xb1, 0x59, 0x45, 0x49,
	0x51, 0xe5, 0x26, 0xe8, 0xd3, 0x71, 0x20, 0x28, 0x44, 0x04, 0xcf, 0x6e, 0x3d, 0x15, 0x0c, 0x48,
	0xe4, 0x6e, 0x6d, 0x29, 0xe4, 0x37, 0x8b, 0xb5, 0x89, 0x12, 0xeb, 0x46, 0x9c, 0x71, 0x68, 0xdb,
	0x23, 0x5c, 0x7b, 0x74, 0x69, 0xc6, 0x1f, 0x1c, 0x70, 0xef, 0xde, 0x49, 0xad, 0x49, 0x75, 0x06,
	0x6d, 0x32, 0x86, 0x63, 0x2d, 0x85, 0x31, 0x6a, 0x67, 0x32, 0x5d, 0x7b, 0x05, 0x43, 0xd6, 0x07,
	0x67, 0x4d, 0x1a, 0xbb, 0xdc, 0x59, 0x63, 0x3d, 0x4e, 0xc5, 0x56, 0xdb, 0x38, 0x15, 0x88, 0xe4,
	0xd2, 0x90, 0x9e, 0x2e, 0xc7, 0xb0, 0x76, 0x53, 0x87, 0x00, 0x74, 0x13, 0x83, 0x36, 0x49, 0xde,
	0xb5, 0x7d, 0x30, 0x46, 0x6c, 0x55, 0x68, 0xfb, 0x07, 0x70, 0x39, 0xc5, 0x78, 0xae, 0xd4, 0x44,
	0xae, 0xcb, 0x5b, 0xa5, 0xa6, 0x5e, 0x42, 0xd4, 0x8e, 0xc7, 0x90, 0x26, 0x2e, 0x94, 0x21, 0xea,
	0x70, 0xe2, 0x42, 0x99, 0xf1, 0xff, 0x00, 0xaf, 0x46, 0xc0, 0x33, 0xf7, 0xd5, 0x43, 0xbd, 0x12,
	0x86, 0xec, 0x77, 0x80, 0x35, 0xae, 0x4c, 0xde, 0xa3, 0xc5, 0xbc, 0xc5, 0x8f, 0x9f, 0xb3, 0x54,
	0x43, 0x0e, 0x77, 0xd7, 0xdb, 0xf0, 0x7c, 0x04, 0xdd, 0x3f, 0x95, 0x0a, 0x91, 0x32, 0x0f, 0xba,
	0x37, 0x55, 0x92, 0x48, 0xad, 0xfd, 0x83, 0xf3, 0x15, 0xb8, 0x61, 0xf3, 0x4a, 0xf8, 0xd0, 0x6f,
	0x92, 0xab, 0x34, 0xf3, 0x0f, 0xd8, 0x19, 0x7c, 0xd3, 0x20, 0x37, 0x3b, 0xed, 0x7c, 0x87, 0x31,
	0x18, 0x36, 0xa5, 0x7f, 0x91, 0x77, 0xbf, 0xc5, 0xc6, 0x30, 0x6a, 0xb0, 0xbf, 0xaa, 0x97, 0x97,
	0xcd, 0xdf, 0xb1, 0x49, 0x56, 0x69, 0xfe, 0xc6, 0x3f, 0x3c, 0xbf, 0x83, 0xc1, 0xf5, 0xde, 0x1f,
	0xe5, 0x04, 0xfc, 0x3d, 0xc0, 0x76, 0x1c, 0x01, 0xdb, 0x43, 0x69, 0x1b, 0xdf, 0x61, 0xa7, 0x70,
	0xbc, 0x87, 0xdf, 0xaa, 0xe2, 0x6d, 0x9c, 0xfb, 0xad, 0xc5, 0x10, 0xf6, 0x1e, 0xe1, 0x3f, 0x46,
	0xff, 0x9d, 0x94, 0x69, 0xfe, 0x90, 0x2a, 0xf9, 0xcb, 0x2e, 0x7e, 0xdf, 0xa1, 0x17, 0xfb, 0xd7,
	0x4f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x53, 0x40, 0xf4, 0x4b, 0xc7, 0x05, 0x00, 0x00,
}
