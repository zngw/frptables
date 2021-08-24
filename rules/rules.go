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

const RuleAllow = "Allow"
const RuleRefuse = "Refuse"
const RuleSkip = "Skip"

// 拦截IP-端口，因为验证ip有时间差，攻击频率太高会导致防火墙重复添加。
var RefuseMap = make(map[string]bool)

func Init() {
	rateInit()
}

func check(ip, name string, port int) {

	if checkAllow(ip, port) {
		log.Trace("link", name, ip, port, "is allow")
		return
	}

	rule, desc, p := checkLocation(ip, port)
	if rule == RuleAllow {
		log.Trace("link", fmt.Sprintf("location allow: [%s]%s:%d %s", name, ip, port, desc))
		return
	}

	if rule == RuleRefuse {
		refuse(ip, name, p)
		log.Trace("add", fmt.Sprintf("location refuse: [%s]%s:%d %s", name, ip, port, desc))
		return
	}

	r, c := checkRate(ip, port)
	if r {
		refuse(ip, name, p)
		log.Trace("add", fmt.Sprintf("rate refuse: [%s]%s:%d %s ->%d", name, ip, port, desc, c))
		return
	}

	log.Trace("link", fmt.Sprintf("rules allow: [%s]%s:%d %s", name, ip, port, desc))
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

func checkLocation(ip string, port int) (rule, desc string, p int) {
	ipInfo := util.GetIpInfo(ip)
	if ipInfo.Status != "success" {
		// 地址获取不成功，跳过
		rule = RuleSkip
		p = -1
		return
	}

	for _, v := range config.Cfg.Rules.Location {
		if v.Port != -1 && v.Port != port {
			// 端口不匹配
			continue
		}

		if v.Country != "" && v.Country != ipInfo.Country {
			// 国家不匹配
			continue
		}

		if v.RegionName != "" && v.Country != ipInfo.Region {
			// 省不匹配
			continue
		}

		if v.City != "" && v.Country != ipInfo.City {
			// 城市不匹配
			continue
		}

		desc = fmt.Sprintf("%s,%s,%s", ipInfo.Country, ipInfo.RegionName, ipInfo.City)
		rule = v.Rules
		p = v.Port
		return
	}

	return
}

func checkRate(ip string, port int) (refuse bool, count int) {
	info := getIpHistory(ip, port)

	for _, v := range config.Cfg.Rules.Rate {
		if v.Port != -1 && v.Port != port {
			// 端口不匹配
			continue
		}

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
