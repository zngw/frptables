//MIT License
//
//Copyright (c) 2021 zngw
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//SOFTWARE.

package rules

import (
	"fmt"
	"github.com/zngw/frptables/config"
	"github.com/zngw/frptables/util"
	"github.com/zngw/log"
)

// 拦截IP-端口，因为验证ip有时间差，攻击频率太高会导致防火墙重复添加。
var RefuseMap = make(map[string]bool)

func Init() {
	rateInit()
}

func rules(ip, name string, port int) {
	// 检查白名单
	if checkAllow(ip, port) {
		log.Trace("link", name, ip, port, "is allow")
		return
	}

	result, desc, p, count := checkRules(ip, port)
	if result {
		refuse(ip, name, p)
		log.Trace("add", fmt.Sprintf("refuse: [%s]%s:%d ->%d, %s", name, ip, port, count, desc))
		return
	} else {
		log.Trace("link", fmt.Sprintf("allow: [%s]%s:%d ->%d, %s", name, ip, port, count, desc))
		return
	}
}

func checkAllow(ip string, port int) bool {
	for _, v := range config.Cfg.AllowIp {
		if v == ip {
			return true
		}
	}

	for _, v := range config.Cfg.AllowPort {
		if v == port {
			return true
		}
	}

	return false
}

func checkRules(ip string, port int) (refuse bool, desc string, p, count int) {
	ipInfo := util.GetIpInfo(ip)
	if ipInfo.Status != "success" {
		// 地址获取不成功，跳过
		refuse = false
		p = -1
		return
	}

	for _, v := range config.Cfg.Rules {
		if v.Port != -1 && v.Port != port {
			// 端口不匹配
			continue
		}

		if v.Country != "" && v.Country != ipInfo.Country {
			// 国家不匹配
			continue
		}

		if v.RegionName != "" && v.RegionName != ipInfo.Region {
			// 省不匹配
			continue
		}

		if v.City != "" && v.City != ipInfo.City {
			// 城市不匹配
			continue
		}

		p = v.Port
		desc = fmt.Sprintf("%s,%s,%s", ipInfo.Country, ipInfo.RegionName, ipInfo.City)

		// 跳过
		if v.Count < 0 {
			refuse = false
			return
		}

		// 拒绝访问
		if v.Count == 0 {
			refuse = true
			return
		}

		info := getIpHistory(ip, port)
		count = info.Count(v.Time)
		if count <= v.Count {
			refuse = false
			return
		}

		refuse = true
		delIpHistory(ip, port)
		return
	}

	return
}

func refuse(ip, name string, port int) {
	// 检测IP是否有添加记录
	if _, ok := RefuseMap[ip]; ok {
		// ip添加
		return
	}

	// 检测IP:Port是否有添加记录
	if port != -1 {
		key := fmt.Sprintf("%s:%d", ip, port)
		if _, ok := RefuseMap[key]; ok {
			return
		}

		RefuseMap[key] = true
	} else {
		RefuseMap[ip] = true
	}

	cmd := ""
	switch config.Cfg.TablesType {
	case "iptables":
		if port == -1 {
			cmd = fmt.Sprintf("iptables -I INPUT -s %s -j DROP", ip)
		} else {
			cmd = fmt.Sprintf("iptables -I INPUT -s %s -ptcp --dport %d -j DROP", ip, port)
		}
	case "firewall":
		if port == -1 {
			cmd = fmt.Sprintf("firewall-cmd --permanent --add-rich-rule=\"rule family=\"ipv4\" source address=\"%s\" reject\"", ip)
		} else {
			cmd = fmt.Sprintf("firewall-cmd --permanent --add-rich-rule=\"rule family=\"ipv4\" source address=\"%s\" port protocol=\"tcp\" port=\"%d\" reject\"", ip, port)
		}
		cmd += "&& firewall-cmd --reload"
	case "md":
		if port == -1 {
			cmd = fmt.Sprintf("netsh advfirewall firewall add rule name=%s-%s dir=in action=block protocol=TCP remoteip=%s", name, ip, ip)
		} else {
			cmd = fmt.Sprintf("netsh advfirewall firewall add rule name=%s-%s-%d dir=in action=block protocol=TCP remoteip=%s localport=%d", name, ip, port, ip, port)
		}
	}

	result := util.Command(cmd)
	if result != "" {
		log.Trace("sys", result)
	}
}
