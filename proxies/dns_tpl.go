package proxies

var baseDns = `dns:
  enable: true
  ipv6: true
  # listen: 0.0.0.0:53
  # enhanced-mode: redir-host # 或 fake-ip
  # # fake-ip-range: 198.18.0.1/16 # 如果你不知道这个参数的作用，请勿修改
  # # 实验性功能 hosts, 支持通配符 (例如 *.clash.dev 甚至 *.foo.*.example.com)
  # # 静态的域名 比 通配域名 具有更高的优先级 (foo.example.com 优先于 *.example.com)
  # # 注意: hosts 在 fake-ip 模式下不生效
  # hosts:
  #   '*.clash.dev': 127.0.0.1
  #   'alpha.clash.dev': '::1'

  nameserver:
     - 1.2.4.8
     - 223.5.5.5
     - 176.103.130.130
     - 114.114.114.114
     - 180.76.76.76
     - 119.29.29.29
     - tls://dns.rubyfish.cn:853
     - tls://dns.alidns.com:853
     - tls://dns.pub:853
     - tls://doh.pub:853
     #- https://dns.rubyfish.cn/dns-query

  fallback: # 与 nameserver 内的服务器列表同时发起请求，当规则符合 GEOIP 在 CN 以外时，fallback 列表内的域名服务器生效。
     - 8.8.8.8
     - 1.1.1.1
     - 176.103.130.130
     - 9.9.9.9
     - tls://dns.rubyfish.cn:853
     - tls://1.0.0.1:853
     - tls://dns.google:853

     #- https://dns.rubyfish.cn/dns-query
     #- https://cloudflare-dns.com/dns-query
     #- https://dns.google/dns-query`

var dnsTpl = `dns:
  enable: true
  ipv6: true
  # enhanced-mode: redir-host # 或 fake-ip
  # # fake-ip-range: 198.18.0.1/16 # 如果你不知道这个参数的作用，请勿修改
  # # 实验性功能 hosts, 支持通配符 (例如 *.clash.dev 甚至 *.foo.*.example.com)
  # # 静态的域名 比 通配域名 具有更高的优先级 (foo.example.com 优先于 *.example.com)
  # # 注意: hosts 在 fake-ip 模式下不生效
  # hosts:
  #   '*.clash.dev': 127.0.0.1
  #   'alpha.clash.dev': '::1'

  nameserver:
{{ range .DnsServiceList}}  - {{ .}}
{{ end}}
  fallback: # 与 nameserver 内的服务器列表同时发起请求，当规则符合 GEOIP 在 CN 以外时，fallback 列表内的域名服务器生效。
     - 8.8.8.8
     - 1.1.1.1
     - 176.103.130.130
     - 9.9.9.9
     - tls://dns.rubyfish.cn:853
     - tls://1.0.0.1:853
     - tls://dns.google:853`
