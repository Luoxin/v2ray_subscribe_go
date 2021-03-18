// Code generated by protoc-gen-go. DO NOT EDIT.
// source: subscription.proto

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
	return fileDescriptor_subscription_042fa9a12bcaadd3, []int{0}
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
	return fileDescriptor_subscription_042fa9a12bcaadd3, []int{1}
}

type CrawlerConf struct {
	Id        uint64 `protobuf:"varint,1,opt,name=id" json:"id,omitempty" gorm:"primaryKey"`
	CreatedAt uint32 `protobuf:"varint,2,opt,name=created_at,json=createdAt" json:"created_at,omitempty" gorm:"autoUpdateTime:autoCreateTime"`
	UpdatedAt uint32 `protobuf:"varint,3,opt,name=updated_at,json=updatedAt" json:"updated_at,omitempty" gorm:"autoUpdateTime:autoUpdateTime"`
	CrawlUrl  string `protobuf:"bytes,4,opt,name=crawl_url,json=crawlUrl" json:"crawl_url,omitempty" gorm:"index:idx_crawl_url,type:text,comment:抓取的地址"`
	CrawlType uint32 `protobuf:"varint,5,opt,name=crawl_type,json=crawlType" json:"crawl_type,omitempty"`
	// @gorm: type:json
	Rule *CrawlerConf_Rule `protobuf:"bytes,6,opt,name=rule" json:"rule,omitempty"" gorm:"type:json"`
	// @grom: index: idx_next_crawl_at
	IsClosed bool `protobuf:"varint,7,opt,name=is_closed,json=isClosed" json:"is_closed,omitempty" gorm:"index:idx_next_crawl_at"`
	// @grom: index: idx_next_crawl_at
	NextAt   uint32 `protobuf:"varint,8,opt,name=next_at,json=nextAt" json:"next_at,omitempty" gorm:"index:idx_next_crawl_at"`
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
	return fileDescriptor_subscription_042fa9a12bcaadd3, []int{0}
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
	return fileDescriptor_subscription_042fa9a12bcaadd3, []int{0, 0}
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
	CreatedAt uint32 `protobuf:"varint,2,opt,name=created_at,json=createdAt" json:"created_at,omitempty" gorm:"autoUpdateTime:autoCreateTime"`
	UpdatedAt uint32 `protobuf:"varint,3,opt,name=updated_at,json=updatedAt" json:"updated_at,omitempty" gorm:"autoUpdateTime:autoUpdateTime`
	// HOST:PORT
	// @v: max=1000
	Url                  string   `protobuf:"bytes,4,opt,name=url" json:"url,omitempty" gorm:"index:idx_crawl_url,type:text,comment:节点的地址"`
	ProxyNodeType        uint32   `protobuf:"varint,5,opt,name=proxy_node_type,json=proxyNodeType" json:"proxy_node_type,omitempty"`
	ProxySpeed           float64  `protobuf:"fixed64,7,opt,name=proxy_speed,json=proxySpeed" json:"proxy_speed,omitempty"`
	ProxyNetworkDelay    float64  `protobuf:"fixed64,8,opt,name=proxy_network_delay,json=proxyNetworkDelay" json:"proxy_network_delay,omitempty"`
	NextCheckAt          uint32   `protobuf:"varint,9,opt,name=next_check_at,json=nextCheckAt" json:"next_check_at,omitempty" gorm:"index:idx_next_check_at"`
	CheckInterval        uint32   `protobuf:"varint,10,opt,name=check_interval,json=checkInterval" json:"check_interval,omitempty"`
	CrawlId              uint64   `protobuf:"varint,11,opt,name=crawl_id,json=crawlId" json:"crawl_id,omitempty"`
	DeathCount           uint32   `protobuf:"varint,12,opt,name=death_count,json=deathCount" json:"death_count,omitempty" gorm:"index:idx_alive"`
	IsClose              bool     `protobuf:"varint,13,opt,name=is_close,json=isClose" json:"is_close,omitempty" gorm:"index:idx_next_check_at;index:idx_alive"`
	LastCrawlerAt        uint32   `protobuf:"varint,14,opt,name=last_crawler_at,json=lastCrawlerAt" json:"last_crawler_at,omitempty"`
	AvailableCount       uint64   `protobuf:"varint,15,opt,name=available_count,json=availableCount" json:"available_count,omitempty" gorm:"index:idx_alive"`
	XXX_NoUnkeyedLiteral struct{} `json:"-" bson:"-" gorm:"-"`
	XXX_unrecognized     []byte   `json:"-" bson:"-" gorm:"-"`
	XXX_sizecache        int32    `json:"-" bson:"-" gorm:"-"`
}

func (m *ProxyNode) Reset()         { *m = ProxyNode{} }
func (m *ProxyNode) String() string { return proto.CompactTextString(m) }
func (*ProxyNode) ProtoMessage()    {}
func (*ProxyNode) Descriptor() ([]byte, []int) {
	return fileDescriptor_subscription_042fa9a12bcaadd3, []int{1}
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

type TohruFeed struct {
	Id        uint64 `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	CreatedAt uint32 `protobuf:"varint,2,opt,name=created_at,json=createdAt" json:"created_at,omitempty"`
	UpdatedAt uint32 `protobuf:"varint,3,opt,name=updated_at,json=updatedAt" json:"updated_at,omitempty"`
	// 用户的唯一标识
	UserId string `protobuf:"bytes,4,opt,name=user_id,json=userId" json:"user_id,omitempty" gorm:"uniqueIndexidx_user_id;not null"`
	// 用户上传的数量
	UpCount              uint64   `protobuf:"varint,5,opt,name=up_count,json=upCount" json:"up_count,omitempty"`
	LastIp               string   `protobuf:"bytes,6,opt,name=last_ip,json=lastIp" json:"last_ip,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-" bson:"-" gorm:"-"`
	XXX_unrecognized     []byte   `json:"-" bson:"-" gorm:"-"`
	XXX_sizecache        int32    `json:"-" bson:"-" gorm:"-"`
}

func (m *TohruFeed) Reset()         { *m = TohruFeed{} }
func (m *TohruFeed) String() string { return proto.CompactTextString(m) }
func (*TohruFeed) ProtoMessage()    {}
func (*TohruFeed) Descriptor() ([]byte, []int) {
	return fileDescriptor_subscription_042fa9a12bcaadd3, []int{2}
}
func (m *TohruFeed) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TohruFeed.Unmarshal(m, b)
}
func (m *TohruFeed) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TohruFeed.Marshal(b, m, deterministic)
}
func (dst *TohruFeed) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TohruFeed.Merge(dst, src)
}
func (m *TohruFeed) XXX_Size() int {
	return xxx_messageInfo_TohruFeed.Size(m)
}
func (m *TohruFeed) XXX_DiscardUnknown() {
	xxx_messageInfo_TohruFeed.DiscardUnknown(m)
}

var xxx_messageInfo_TohruFeed proto.InternalMessageInfo

func (m *TohruFeed) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *TohruFeed) GetCreatedAt() uint32 {
	if m != nil {
		return m.CreatedAt
	}
	return 0
}

func (m *TohruFeed) GetUpdatedAt() uint32 {
	if m != nil {
		return m.UpdatedAt
	}
	return 0
}

func (m *TohruFeed) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *TohruFeed) GetUpCount() uint64 {
	if m != nil {
		return m.UpCount
	}
	return 0
}

func (m *TohruFeed) GetLastIp() string {
	if m != nil {
		return m.LastIp
	}
	return ""
}

func init() {
	proto.RegisterType((*CrawlerConf)(nil), "domain.CrawlerConf")
	proto.RegisterType((*CrawlerConf_Rule)(nil), "domain.CrawlerConf.Rule")
	proto.RegisterType((*ProxyNode)(nil), "domain.ProxyNode")
	proto.RegisterType((*TohruFeed)(nil), "domain.TohruFeed")
	proto.RegisterEnum("domain.ErrCode", ErrCode_name, ErrCode_value)
	proto.RegisterEnum("domain.ProxyNodeType", ProxyNodeType_name, ProxyNodeType_value)
}

func init() { proto.RegisterFile("subscription.proto", fileDescriptor_subscription_042fa9a12bcaadd3) }

var fileDescriptor_subscription_042fa9a12bcaadd3 = []byte{
	// 697 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x54, 0xcd, 0x6e, 0xda, 0x4c,
	0x14, 0x8d, 0x09, 0x60, 0xfb, 0xf2, 0x01, 0xce, 0x24, 0x5f, 0xf0, 0x97, 0xef, 0x0f, 0x21, 0xb5,
	0x45, 0x51, 0xc5, 0x22, 0x5d, 0x75, 0xd1, 0x05, 0x25, 0xad, 0xca, 0x26, 0x8a, 0xc6, 0x69, 0x17,
	0xdd, 0x58, 0x83, 0x3d, 0x11, 0x6e, 0x26, 0x1e, 0x6b, 0x3c, 0x4e, 0xc2, 0x0b, 0xf5, 0x11, 0xfa,
	0x3e, 0xdd, 0xf4, 0x39, 0xaa, 0x3b, 0x63, 0x48, 0x10, 0xdb, 0x6c, 0xd0, 0xdc, 0x73, 0xe6, 0xce,
	0x3d, 0x9c, 0x7b, 0x64, 0x20, 0x65, 0xb5, 0x28, 0x13, 0x95, 0x15, 0x3a, 0x93, 0xf9, 0xa4, 0x50,
	0x52, 0x4b, 0xd2, 0x4e, 0xe5, 0x2d, 0xcb, 0xf2, 0xd1, 0xaf, 0x06, 0x74, 0x66, 0x8a, 0xdd, 0x0b,
	0xae, 0x66, 0x32, 0xbf, 0x26, 0x3d, 0x68, 0x64, 0x69, 0xe8, 0x0c, 0x9d, 0x71, 0x93, 0x36, 0xb2,
	0x94, 0xfc, 0x0b, 0x90, 0x28, 0xce, 0x34, 0x4f, 0x63, 0xa6, 0xc3, 0xc6, 0xd0, 0x19, 0x77, 0xa9,
	0x5f, 0x23, 0x53, 0x8d, 0x74, 0x55, 0xa4, 0x6b, 0x7a, 0xdf, 0xd2, 0x35, 0x32, 0xd5, 0xe4, 0x6f,
	0xf0, 0x13, 0x7c, 0x3c, 0xae, 0x94, 0x08, 0x9b, 0x43, 0x67, 0xec, 0x53, 0xcf, 0x00, 0x9f, 0x95,
	0xb0, 0x4f, 0x23, 0xa9, 0x57, 0x05, 0x0f, 0x5b, 0xeb, 0xa7, 0xd9, 0xbd, 0xb8, 0x5a, 0x15, 0x9c,
	0xbc, 0x86, 0xa6, 0xaa, 0x04, 0x0f, 0xdb, 0x43, 0x67, 0xdc, 0x39, 0x0b, 0x27, 0x56, 0xf0, 0xe4,
	0x89, 0xd8, 0x09, 0xad, 0x04, 0xa7, 0xe6, 0x16, 0x4e, 0xca, 0xca, 0x38, 0x11, 0xb2, 0xe4, 0x69,
	0xe8, 0x0e, 0x9d, 0xb1, 0x47, 0xbd, 0xac, 0x9c, 0x99, 0x9a, 0x0c, 0xc0, 0xcd, 0xf9, 0x83, 0x46,
	0x89, 0x9e, 0x19, 0xd3, 0xc6, 0x72, 0xaa, 0xc9, 0x09, 0x78, 0x59, 0xae, 0xb9, 0xba, 0x63, 0x22,
	0xf4, 0x0d, 0xb3, 0xa9, 0x09, 0x81, 0x66, 0x2e, 0x35, 0x0f, 0xc1, 0xc8, 0x36, 0xe7, 0x93, 0xb7,
	0xd0, 0xa4, 0xf5, 0xb4, 0xaa, 0xe4, 0x71, 0xa1, 0xe4, 0xc3, 0xca, 0x98, 0xe5, 0x51, 0xaf, 0x2a,
	0xf9, 0x25, 0xd6, 0xe4, 0x08, 0x5a, 0x0f, 0x05, 0xd3, 0x4b, 0xe3, 0x96, 0x4f, 0x6d, 0x31, 0xfa,
	0xd1, 0x04, 0xdf, 0xf0, 0x17, 0x32, 0xe5, 0xcf, 0x6c, 0x73, 0x00, 0xfb, 0x8f, 0x06, 0xe3, 0x91,
	0xbc, 0x84, 0xbe, 0x11, 0x17, 0xe7, 0x32, 0xe5, 0x4f, 0x0d, 0xee, 0x16, 0x6b, 0x0d, 0xc6, 0xe4,
	0x77, 0xd0, 0x31, 0x37, 0x52, 0xae, 0x59, 0x26, 0x6a, 0xaf, 0xff, 0x59, 0x7b, 0xbd, 0xd1, 0x3b,
	0xc1, 0x9f, 0x73, 0x73, 0x87, 0x42, 0xbe, 0x39, 0x93, 0xff, 0xa1, 0x63, 0xc7, 0x94, 0x05, 0xaf,
	0x7d, 0x77, 0x28, 0x18, 0x28, 0x42, 0x84, 0x4c, 0xe0, 0xb0, 0xd6, 0xc1, 0xf5, 0xbd, 0x54, 0x37,
	0x71, 0xca, 0x05, 0x5b, 0x99, 0x2d, 0x38, 0xf4, 0xc0, 0x6a, 0xb1, 0xcc, 0x39, 0x12, 0x64, 0x04,
	0x5d, 0xb3, 0xa9, 0x64, 0xc9, 0x93, 0x1b, 0xfc, 0xaf, 0x76, 0x2b, 0x1d, 0x04, 0x67, 0x88, 0x4d,
	0x35, 0x79, 0x01, 0x3d, 0x4b, 0x6f, 0x56, 0x07, 0xf6, 0xaf, 0x19, 0x74, 0xbe, 0xde, 0xdf, 0x5f,
	0x60, 0xa3, 0x16, 0x67, 0x69, 0xd8, 0x31, 0x46, 0xbb, 0xa6, 0x9e, 0xa7, 0x28, 0x3b, 0xe5, 0x4c,
	0x2f, 0xe3, 0x44, 0x56, 0xb9, 0x0e, 0xff, 0x30, 0xed, 0x60, 0xa0, 0x19, 0x22, 0xd8, 0xbb, 0x4e,
	0x53, 0xd8, 0x35, 0xeb, 0x75, 0xeb, 0x30, 0xa1, 0xb3, 0x82, 0x95, 0x3a, 0x4e, 0x6c, 0x0e, 0x51,
	0x63, 0xcf, 0x8e, 0x47, 0xb8, 0x4e, 0xe7, 0x54, 0x93, 0x57, 0xd0, 0x67, 0x77, 0x2c, 0x13, 0x6c,
	0x21, 0x78, 0x3d, 0xa7, 0x6f, 0x54, 0xf4, 0x36, 0xb0, 0x99, 0x75, 0xf2, 0x1f, 0xc0, 0xa3, 0xbb,
	0xb8, 0xca, 0x45, 0x75, 0x6d, 0x92, 0xe1, 0x53, 0x3c, 0x8e, 0xbe, 0x3b, 0xe0, 0x5f, 0xc9, 0xa5,
	0xaa, 0x3e, 0xa2, 0xa1, 0xcf, 0x1b, 0x9c, 0x01, 0xb8, 0x55, 0xc9, 0x15, 0x5a, 0x64, 0xc3, 0xd3,
	0xc6, 0x72, 0x9e, 0xa2, 0x01, 0x55, 0x51, 0xcb, 0x6e, 0x59, 0xf3, 0xaa, 0xc2, 0x7a, 0x33, 0x00,
	0xd7, 0x18, 0x90, 0x15, 0x26, 0x2e, 0x3e, 0x6d, 0x63, 0x39, 0x2f, 0x4e, 0x8f, 0xc1, 0xfd, 0xa0,
	0xd4, 0x0c, 0xe3, 0xdd, 0x01, 0x37, 0xaa, 0x92, 0x84, 0x97, 0x65, 0xb0, 0x77, 0xfa, 0xd3, 0x81,
	0xee, 0xe5, 0x56, 0xea, 0x8e, 0x20, 0xd8, 0x02, 0x2e, 0x32, 0x11, 0xec, 0x91, 0x63, 0x20, 0x5b,
	0xe8, 0x97, 0x5b, 0xec, 0x77, 0xc8, 0x00, 0x0e, 0xb7, 0xf0, 0x2b, 0x25, 0xbf, 0xb1, 0x3c, 0x68,
	0xec, 0x36, 0x08, 0x6c, 0xd8, 0x27, 0x87, 0xd0, 0xdf, 0xc2, 0xa3, 0x28, 0x68, 0xee, 0xcc, 0x8c,
	0x22, 0x1a, 0xb4, 0x76, 0xde, 0x8e, 0x64, 0x72, 0xc3, 0x75, 0xd0, 0x26, 0x7f, 0xc2, 0xc1, 0x16,
	0xf1, 0x49, 0xeb, 0x22, 0x70, 0x77, 0x46, 0x46, 0x39, 0x17, 0x22, 0xf0, 0xce, 0x3c, 0xa8, 0x3f,
	0xa8, 0xef, 0xbd, 0xaf, 0xf5, 0x69, 0xd1, 0x36, 0x5f, 0xda, 0x37, 0xbf, 0x03, 0x00, 0x00, 0xff,
	0xff, 0x76, 0x28, 0x4c, 0x73, 0x7f, 0x05, 0x00, 0x00,
}
