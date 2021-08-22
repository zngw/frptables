// @Title
// @Description $
// @Author  55
// @Date  2021/8/22
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

}
