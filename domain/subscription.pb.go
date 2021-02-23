// Code generated by protoc-gen-go. DO NOT EDIT.
// source: domain/subscription.proto

package domain

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
	return fileDescriptor_subscription_7c0876797555260e, []int{0}
}

type ProxyNodeType int32

const (
	ProxyNodeType_ProxyNodeTypeNil    ProxyNodeType = 0
	ProxyNodeType_ProxyNodeTypeVmess  ProxyNodeType = 1
	ProxyNodeType_ProxyNodeTypeTrojan ProxyNodeType = 2
	ProxyNodeType_ProxyNodeTypeVless  ProxyNodeType = 3
	ProxyNodeType_ProxyNodeTypeSS     ProxyNodeType = 4
	ProxyNodeType_ProxyNodeTypeSSR    ProxyNodeType = 5
	ProxyNodeType_ProxyNodeTypeSocket ProxyNodeType = 6
	ProxyNodeType_ProxyNodeTypeHttp   ProxyNodeType = 7
	ProxyNodeType_ProxyNodeTypeSnell  ProxyNodeType = 8
)

var ProxyNodeType_name = map[int32]string{
	0: "ProxyNodeTypeNil",
	1: "ProxyNodeTypeVmess",
	2: "ProxyNodeTypeTrojan",
	3: "ProxyNodeTypeVless",
	4: "ProxyNodeTypeSS",
	5: "ProxyNodeTypeSSR",
	6: "ProxyNodeTypeSocket",
	7: "ProxyNodeTypeHttp",
	8: "ProxyNodeTypeSnell",
}
var ProxyNodeType_value = map[string]int32{
	"ProxyNodeTypeNil":    0,
	"ProxyNodeTypeVmess":  1,
	"ProxyNodeTypeTrojan": 2,
	"ProxyNodeTypeVless":  3,
	"ProxyNodeTypeSS":     4,
	"ProxyNodeTypeSSR":    5,
	"ProxyNodeTypeSocket": 6,
	"ProxyNodeTypeHttp":   7,
	"ProxyNodeTypeSnell":  8,
}

func (x ProxyNodeType) String() string {
	return proto.EnumName(ProxyNodeType_name, int32(x))
}
func (ProxyNodeType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_subscription_7c0876797555260e, []int{1}
}

type CrawlerConf struct {
	Id        uint64 `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	CreatedAt uint32 `protobuf:"varint,2,opt,name=created_at,json=createdAt" json:"created_at,omitempty"`
	UpdatedAt uint32 `protobuf:"varint,3,opt,name=updated_at,json=updatedAt" json:"updated_at,omitempty"`
	// @gorm: type:varchar(1000);unique_index:idx_crawl_url
	CrawlUrl  string `protobuf:"bytes,4,opt,name=crawl_url,json=crawlUrl" json:"crawl_url,omitempty"`
	CrawlType uint32 `protobuf:"varint,5,opt,name=crawl_type,json=crawlType" json:"crawl_type,omitempty"`
	// @gorm: type:json
	Rule *CrawlerConf_Rule `protobuf:"bytes,6,opt,name=rule" json:"rule,omitempty"`
	// @grom: index: idx_next_crawl_at
	IsClosed bool `protobuf:"varint,7,opt,name=is_closed,json=isClosed" json:"is_closed,omitempty"`
	// @grom: index: idx_next_crawl_at
	NextAt   uint32 `protobuf:"varint,8,opt,name=next_at,json=nextAt" json:"next_at,omitempty"`
	Interval uint32 `protobuf:"varint,9,opt,name=interval" json:"interval,omitempty"`
	// @gorm: type:varchar(100)
	// @v: max=100
	Note                 string   `protobuf:"bytes,10,opt,name=note" json:"note,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-" bson:"-" gorm:"-"`
	XXX_unrecognized     []byte   `json:"-" bson:"-" gorm:"-"`
	XXX_sizecache        int32    `json:"-" bson:"-" gorm:"-"`
}

func (m *CrawlerConf) Reset()         { *m = CrawlerConf{} }
func (m *CrawlerConf) String() string { return proto.CompactTextString(m) }
func (*CrawlerConf) ProtoMessage()    {}
func (*CrawlerConf) Descriptor() ([]byte, []int) {
	return fileDescriptor_subscription_7c0876797555260e, []int{0}
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
	return fileDescriptor_subscription_7c0876797555260e, []int{0, 0}
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
	// @gorm: type:varchar(1000); unique_index: idx_proxy_node_url
	// @v: max=1000
	Url           string `protobuf:"bytes,4,opt,name=url" json:"url,omitempty"`
	ProxyNodeType uint32 `protobuf:"varint,5,opt,name=proxy_node_type,json=proxyNodeType" json:"proxy_node_type,omitempty"`
	// @gorm: type:text
	NodeDetail        *ProxyNode_NodeDetail `protobuf:"bytes,6,opt,name=node_detail,json=nodeDetail" json:"node_detail,omitempty"`
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
	AvailableCount       uint64   `protobuf:"varint,15,opt,name=available_count,json=availableCount" json:"available_count,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-" bson:"-" gorm:"-"`
	XXX_unrecognized     []byte   `json:"-" bson:"-" gorm:"-"`
	XXX_sizecache        int32    `json:"-" bson:"-" gorm:"-"`
}

func (m *ProxyNode) Reset()         { *m = ProxyNode{} }
func (m *ProxyNode) String() string { return proto.CompactTextString(m) }
func (*ProxyNode) ProtoMessage()    {}
func (*ProxyNode) Descriptor() ([]byte, []int) {
	return fileDescriptor_subscription_7c0876797555260e, []int{1}
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

func (m *ProxyNode) GetAvailableCount() uint64 {
	if m != nil {
		return m.AvailableCount
	}
	return 0
}

type ProxyNode_NodeDetail struct {
	Buf                  string   `protobuf:"bytes,1,opt,name=buf" json:"buf,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-" bson:"-" gorm:"-"`
	XXX_unrecognized     []byte   `json:"-" bson:"-" gorm:"-"`
	XXX_sizecache        int32    `json:"-" bson:"-" gorm:"-"`
}

func (m *ProxyNode_NodeDetail) Reset()         { *m = ProxyNode_NodeDetail{} }
func (m *ProxyNode_NodeDetail) String() string { return proto.CompactTextString(m) }
func (*ProxyNode_NodeDetail) ProtoMessage()    {}
func (*ProxyNode_NodeDetail) Descriptor() ([]byte, []int) {
	return fileDescriptor_subscription_7c0876797555260e, []int{1, 0}
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

func init() {
	proto.RegisterType((*CrawlerConf)(nil), "domain.CrawlerConf")
	proto.RegisterType((*CrawlerConf_Rule)(nil), "domain.CrawlerConf.Rule")
	proto.RegisterType((*ProxyNode)(nil), "domain.ProxyNode")
	proto.RegisterType((*ProxyNode_NodeDetail)(nil), "domain.ProxyNode.NodeDetail")
	proto.RegisterEnum("domain.ErrCode", ErrCode_name, ErrCode_value)
	proto.RegisterEnum("domain.ProxyNodeType", ProxyNodeType_name, ProxyNodeType_value)
}

func init() {
	proto.RegisterFile("domain/subscription.proto", fileDescriptor_subscription_7c0876797555260e)
}

var fileDescriptor_subscription_7c0876797555260e = []byte{
	// 646 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x54, 0x4f, 0x4f, 0xdb, 0x4e,
	0x10, 0xc5, 0xc1, 0xc4, 0xf6, 0xf8, 0x97, 0xc4, 0x2c, 0xfc, 0xc0, 0xd0, 0x7f, 0x11, 0x52, 0xdb,
	0x08, 0x55, 0xa9, 0x44, 0x4f, 0x3d, 0xf4, 0x90, 0x86, 0x4a, 0xe5, 0x82, 0xd0, 0x9a, 0xf6, 0xd0,
	0x8b, 0xb5, 0xb1, 0x17, 0xe1, 0xb2, 0x78, 0xad, 0xf5, 0x1a, 0xc8, 0x17, 0xeb, 0xf7, 0xe9, 0xa5,
	0x9f, 0xa3, 0xda, 0x59, 0x27, 0x10, 0xe5, 0xda, 0x8b, 0x35, 0xf3, 0xde, 0xce, 0xce, 0xf3, 0x9b,
	0xd1, 0xc2, 0x41, 0x2e, 0x6f, 0x59, 0x51, 0xbe, 0xaf, 0x9b, 0x59, 0x9d, 0xa9, 0xa2, 0xd2, 0x85,
	0x2c, 0xc7, 0x95, 0x92, 0x5a, 0x92, 0xae, 0xa5, 0x8e, 0xfe, 0x74, 0x20, 0x9c, 0x2a, 0x76, 0x2f,
	0xb8, 0x9a, 0xca, 0xf2, 0x8a, 0xf4, 0xa1, 0x53, 0xe4, 0xb1, 0x33, 0x74, 0x46, 0x2e, 0xed, 0x14,
	0x39, 0x79, 0x01, 0x90, 0x29, 0xce, 0x34, 0xcf, 0x53, 0xa6, 0xe3, 0xce, 0xd0, 0x19, 0xf5, 0x68,
	0xd0, 0x22, 0x13, 0x6d, 0xe8, 0xa6, 0xca, 0x17, 0xf4, 0xa6, 0xa5, 0x5b, 0x64, 0xa2, 0xc9, 0x33,
	0x08, 0x32, 0x73, 0x79, 0xda, 0x28, 0x11, 0xbb, 0x43, 0x67, 0x14, 0x50, 0x1f, 0x81, 0x6f, 0x4a,
	0xd8, 0xab, 0x0d, 0xa9, 0xe7, 0x15, 0x8f, 0xb7, 0x16, 0x57, 0xb3, 0x7b, 0x71, 0x39, 0xaf, 0x38,
	0x79, 0x07, 0xae, 0x6a, 0x04, 0x8f, 0xbb, 0x43, 0x67, 0x14, 0x9e, 0xc4, 0x63, 0x2b, 0x78, 0xfc,
	0x44, 0xec, 0x98, 0x36, 0x82, 0x53, 0x3c, 0x65, 0x3a, 0x15, 0x75, 0x9a, 0x09, 0x59, 0xf3, 0x3c,
	0xf6, 0x86, 0xce, 0xc8, 0xa7, 0x7e, 0x51, 0x4f, 0x31, 0x27, 0xfb, 0xe0, 0x95, 0xfc, 0x41, 0x1b,
	0x89, 0x3e, 0xb6, 0xe9, 0x9a, 0x74, 0xa2, 0xc9, 0x21, 0xf8, 0x45, 0xa9, 0xb9, 0xba, 0x63, 0x22,
	0x0e, 0x90, 0x59, 0xe6, 0x84, 0x80, 0x5b, 0x4a, 0xcd, 0x63, 0x40, 0xd9, 0x18, 0x1f, 0x7e, 0x04,
	0x97, 0xb6, 0xdd, 0x9a, 0x9a, 0xa7, 0x95, 0x92, 0x0f, 0x73, 0x34, 0xcb, 0xa7, 0x7e, 0x53, 0xf3,
	0x0b, 0x93, 0x93, 0x5d, 0xd8, 0x7a, 0xa8, 0x98, 0xbe, 0x46, 0xb7, 0x02, 0x6a, 0x93, 0xa3, 0x5f,
	0x2e, 0x04, 0xc8, 0x9f, 0xcb, 0x9c, 0xff, 0x63, 0x9b, 0x23, 0xd8, 0x7c, 0x34, 0xd8, 0x84, 0xe4,
	0x0d, 0x0c, 0x50, 0x5c, 0x5a, 0xca, 0x9c, 0x3f, 0x35, 0xb8, 0x57, 0x2d, 0x34, 0xa0, 0xc9, 0x9f,
	0x20, 0xc4, 0x13, 0x39, 0xd7, 0xac, 0x10, 0xad, 0xd7, 0xcf, 0x17, 0x5e, 0x2f, 0xf5, 0x8e, 0xcd,
	0xe7, 0x14, 0xcf, 0x50, 0x28, 0x97, 0x31, 0x79, 0x05, 0xa1, 0x6d, 0x53, 0x57, 0xbc, 0xf5, 0xdd,
	0xa1, 0x80, 0x50, 0x62, 0x10, 0x32, 0x86, 0x9d, 0x56, 0x07, 0xd7, 0xf7, 0x52, 0xdd, 0xa4, 0x39,
	0x17, 0x6c, 0x8e, 0x53, 0x70, 0xe8, 0xb6, 0xd5, 0x62, 0x99, 0x53, 0x43, 0x90, 0x23, 0xe8, 0xe1,
	0xa4, 0xb2, 0x6b, 0x9e, 0xdd, 0x98, 0x7f, 0xb5, 0x53, 0x09, 0x0d, 0x38, 0x35, 0xd8, 0x44, 0x93,
	0xd7, 0xd0, 0xb7, 0xf4, 0x72, 0x74, 0x60, 0x7f, 0x0d, 0xd1, 0xb3, 0xc5, 0xfc, 0x0e, 0xc0, 0xae,
	0x5a, 0x5a, 0xe4, 0x71, 0x88, 0x46, 0x7b, 0x98, 0x9f, 0xe5, 0x46, 0x76, 0xce, 0x99, 0xbe, 0x4e,
	0x33, 0xd9, 0x94, 0x3a, 0xfe, 0x0f, 0xcb, 0x01, 0xa1, 0xa9, 0x41, 0x4c, 0xed, 0x62, 0x9b, 0xe2,
	0x1e, 0x8e, 0xd7, 0x6b, 0x97, 0xc9, 0x38, 0x2b, 0x58, 0xad, 0xd3, 0xcc, 0xee, 0xa1, 0xd1, 0xd8,
	0xb7, 0xed, 0x0d, 0xdc, 0x6e, 0xe7, 0x44, 0x93, 0xb7, 0x30, 0x60, 0x77, 0xac, 0x10, 0x6c, 0x26,
	0x78, 0xdb, 0x67, 0x80, 0x2a, 0xfa, 0x4b, 0x18, 0x7b, 0x1d, 0xbe, 0x04, 0x78, 0x74, 0xd7, 0x8c,
	0x72, 0xd6, 0x5c, 0xe1, 0x66, 0x04, 0xd4, 0x84, 0xc7, 0x7b, 0xe0, 0x7d, 0x51, 0x6a, 0x6a, 0xb6,
	0x26, 0x04, 0x2f, 0x69, 0xb2, 0x8c, 0xd7, 0x75, 0xb4, 0x71, 0xfc, 0xdb, 0x81, 0xde, 0xc5, 0xca,
	0x30, 0x77, 0x21, 0x5a, 0x01, 0xce, 0x0b, 0x11, 0x6d, 0x90, 0x3d, 0x20, 0x2b, 0xe8, 0xf7, 0x5b,
	0x53, 0xef, 0x90, 0x7d, 0xd8, 0x59, 0xc1, 0x2f, 0x95, 0xfc, 0xc9, 0xca, 0xa8, 0xb3, 0x5e, 0x20,
	0x4c, 0xc1, 0x26, 0xd9, 0x81, 0xc1, 0x0a, 0x9e, 0x24, 0x91, 0xbb, 0xd6, 0x33, 0x49, 0x68, 0xb4,
	0xb5, 0x76, 0x77, 0x22, 0xb3, 0x1b, 0xae, 0xa3, 0x2e, 0xf9, 0x1f, 0xb6, 0x57, 0x88, 0xaf, 0x5a,
	0x57, 0x91, 0xb7, 0xd6, 0x32, 0x29, 0xb9, 0x10, 0x91, 0x7f, 0xe2, 0x43, 0xfb, 0x4e, 0x7d, 0xf6,
	0x7f, 0xb4, 0xd1, 0xac, 0x8b, 0x0f, 0xd8, 0x87, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x72, 0xd5,
	0x04, 0xff, 0xdd, 0x04, 0x00, 0x00,
}
