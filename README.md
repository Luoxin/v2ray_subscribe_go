# Eutamias
v2ray_subscribe的golang版本

## 使用方法
### 使用编译包
- 下载最新的包
- `Eutamias [-c <config file path>]`
### 自编译
- `git clone https://github.com/Luoxin/Eutamias.git`
- `cd Eutamias`
- `CGO_ENABLED=1 GO111MODULE=on go build -o Eutamias ./cmd/Eutamias.go`

### TODO
- [ ] 支持更多协议的抓取和检测
  - [x] ss
  - [x] ssr
  - [x] vmess
  - [x] trojan
- [ ] 支持更多的订阅方式
	- [x] clash
	- [x] v2ray
	- [ ] surge
- [ ] 优化订阅返回的节点策略以及排序
- [x] 本地代理支持
- [ ] 桌面客户端
  - [ ] 包括跨平台的支持
- [ ] pac的优化
