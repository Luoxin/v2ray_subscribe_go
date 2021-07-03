# Eutamias
v2ray_subscribe的golang版本

## 使用方法
### 获取二进制包
#### 使用编译包
- 下载最新的包
- `Eutamias [-c <config file path>]`
#### 自编译

- `git clone https://github.com/Luoxin/Eutamias.git`
- `cd Eutamias`
- `CGO_ENABLED=1 GO111MODULE=on go build -o eutamias ./cmd/Eutamias.go`

### 直接运行服务

`./eutamias [-c <config file path>]`

### 安装服务(需要管理员权限)

`./eutamias -s install`

### 卸载服务(需要管理员权限)

`./eutamias -s uninstall`

### 使用docker运行

- 下载镜像
  `docker pull luoxintt/eutamias`
- 运行
  `docker run -itd -p 7890:7890 luoxintt/eutamias`

### TODO

- [ ] docker支持
- [ ] 支持更多协议的抓取和检测
	- [x] ss
	- [x] ssr
	- [x] vmess
	- [x] trojan
	- [ ] http
	- [ ] socket
- [ ] 支持更多的订阅方式
	- [x] clash
	- [x] v2ray
	- [ ] surge
- [ ] 内置DNS支持
	- [x] dns查询
	- [x] dns服务
	- [x] clash支持内置dns
	- [ ] ipv6支持
	- [x] dot支持
	- [x] doh支持
	- [ ] 翻墙dns查询内置
- [ ] 优化订阅返回的节点策略以及排序
- [x] 本地代理支持
- [ ] 桌面客户端
	- [ ] 包括跨平台的支持
- [ ] pac的优化
	- [ ] 指定proxy地址
	- [ ] 对齐proxy rule
- [ ] clash订阅规则优化
	- [ ] 完善rule list
	- [ ] 支持订阅host
- [ ] 多数据库支持
	- [ ] memory
	- [x] sqlite
	- [x] mysql
	- [x] postgres
- [ ] 多缓存支持
	- [ ] memory
	- [ ] redis
	- [ ] etcd
	- [ ] sqlite
	- [ ] mysql
	- [ ] postgres
