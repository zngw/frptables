[comment]: <> (dtapps)
[![GitHub Org's stars](https://img.shields.io/github/stars/zngw)](https://github.com/zngw)

[comment]: <> (go)
[![godoc](https://pkg.go.dev/badge/github.com/zngw/frptables?status.svg)](https://pkg.go.dev/github.com/zngw/frptables)
[![oproxy.cn](https://goproxy.cn/stats/github.com/zngw/frptables/badges/download-count.svg)](https://goproxy.cn/stats/github.com/zngw/frptables)
[![goreportcard.com](https://goreportcard.com/badge/github.com/zngw/frptables)](https://goreportcard.com/report/github.com/zngw/frptables)
[![deps.dev](https://img.shields.io/badge/deps-go-red.svg)](https://deps.dev/go/github.com%2Fdtapps%2Fgo-ssh-tunnel)

[comment]: <> (github.com)
[![watchers](https://badgen.net/github/watchers/zngw/frptables)](https://github.com/zngw/frptables/watchers)
[![stars](https://badgen.net/github/stars/zngw/frptables)](https://github.com/zngw/frptables/stargazers)
[![forks](https://badgen.net/github/forks/zngw/frptables)](https://github.com/zngw/frptables/network/members)
[![issues](https://badgen.net/github/issues/zngw/frptables)](https://github.com/zngw/frptables/issues)
[![branches](https://badgen.net/github/branches/zngw/frptables)](https://github.com/zngw/frptables/branches)
[![releases](https://badgen.net/github/releases/zngw/frptables)](https://github.com/zngw/frptables/releases)
[![tags](https://badgen.net/github/tags/zngw/frptables)](https://github.com/zngw/frptables/tags)
[![license](https://badgen.net/github/license/zngw/frptables)](https://github.com/zngw/frptables/blob/master/LICENSE)
[![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/zngw/frptables)](https://github.com/zngw/frptables)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/zngw/frptables)](https://github.com/zngw/frptables/releases)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/zngw/frptables)](https://github.com/zngw/frptables/tags)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/zngw/frptables)](https://github.com/zngw/frptables/pulls)
[![GitHub issues](https://img.shields.io/github/issues/zngw/frptables)](https://github.com/zngw/frptables/issues)
[![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/zngw/frptables)](https://github.com/zngw/frptables)
[![GitHub language count](https://img.shields.io/github/languages/count/zngw/frptables)](https://github.com/zngw/frptables)
[![GitHub search hit counter](https://img.shields.io/github/search/zngw/frptables/go)](https://github.com/zngw/frptables)
[![GitHub top language](https://img.shields.io/github/languages/top/zngw/frptables)](https://github.com/zngw/frptables)

# frptables
监控frps日志文件，自定义规则，使用系统自带的防火墙(iptables、firewall、Microsoft Defender)拦截tcp连接的ip，防止暴力破解

## 前提
1、因为需要用到命令行修改系统防火墙，所以运行程序需要root或管理员权限

2、需要frps在0.36及以上版本，并且开启日志功能，日志输出等级为info。
```ini
#日志输出，可以设置为具体的日志文件或者console
log_file = /usr/local/frp/log/frps.log

#日志记录等级，有trace, debug, info, warn, error
log_level = info
```

tcp连接日志格式为
`2021/08/21 15:35:29 [I] [proxy.go:162] [f1aec30e84827422] [ZNGW] get a user connection [210.0.159.76:32832]`

如果日志格式发生了调整，需要修改`rules/parse.go`中的解析规则

## 注意
本程序只对frp日志进行分析，添加防火墙阻止消息，若有误加入防火墙的，需要手动删除。

## 配置

```yaml
# frps日志文件
frps_log: ./log/frps.log

# 输出日志目录
logs: ./tlog/

# frps 名字端口对应配置
name_port:
  "ZNGW": 3389

# 启用防火墙类型 iptables / firewall / md (Microsoft Defender)
tables_type: iptables

# ip白名单:
allow_ip:
  - 127.0.0.1

# 端口白名单
allow_port:
  - 80
  - 443

# 规则访问
rules:

  # 按数组顺序来，匹配到了就按匹配的规则执行，跳过此规则。
  # 地区 country-国家， regionName-省名，名字中不带省字， city-市名，名字中也不带市字
  # 端口: -1 所有端口
  # time: 时间区间
  # count: 访问次数，-1不限，0限制。其他为 time时间内访问count次，超出频率就限制

  - # 中国上海IP允许
    port: -1
    country: 中国
    regionName: 上海
    city: 上海
    time: 1
    count: -1

  - # 中国地区IP 10分钟3次，超出这频率添加防火墙
    port: -1
    country: 中国
    regionName: 浙江
    city:
    time: 600
    count: 3

  - # 其他地区IP 直接加入防火墙
    port: -1
    country:
    regionName:
    city:
    time: 1
    count: 0
```

## 启动
直接使用`nohup ./frptables -c config.yml &`启动

也可以新建`/etc/systemd/system/frptables.service`文件加入系统,以服务方式启动
```ini
[Unit]
Description=frps daemon
After=syslog.target  network.target
Wants=network.target

[Service]
Type=simple
ExecStart=/usr/local/frp/frptables -c /usr/local/frp/config.yml
Restart= always
RestartSec=1min

[Install]
WantedBy=multi-user.target

```

