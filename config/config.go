// @Title
// @Description $
// @Author  55
// @Date  2021/8/22
package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

var Cfg Conf

// 配置文件结构体
type Conf struct {
	FrpsLog    string         `yaml:"frps_log,omitempty"`    // 监听的frps日志文件
	Logs       string         `yaml:"logs,omitempty"`        // 自身日志文件
	NamePort   map[string]int `yaml:"name_port,omitempty"`   // 名字端口对应表
	TablesType string         `yaml:"tables_type,omitempty"` // 启用防火墙类型
	AllowIp    []string       `yaml:"allow_ip,omitempty"`    // ip白名单
	AllowPort  []int          `yaml:"allow_port,omitempty"`  // 端口白名单
	Rules      CfgRules       `yaml:"rules,omitempty"`       // 规则访问
}

type CfgRules struct {
	Location    []CfgLocation `yaml:"location,omitempty"` // IP区域规则
	Rate        []CfgRate     `yaml:"rate,omitempty"`     // IP频率访问频率
	RateMaxTime int64         `yaml:"-"`                  // IP频率中最高超时时间
}

type CfgLocation struct {
	Port       int    `yaml:"port,omitempty"`       // 端口
	Country    string `yaml:"country,omitempty"`    // 国家
	RegionName string `yaml:"regionName,omitempty"` // 省
	City       string `yaml:"city,omitempty"`       // 城市
	Rules      string `yaml:"rules,omitempty"`      // 规则
}

type CfgRate struct {
	Port  int   `yaml:"port,omitempty"`  // 端口
	Time  int64 `yaml:"time,omitempty"`  // 时间
	Count int   `yaml:"count,omitempty"` // 次数
}

func (c *Conf) Init(file string) (err error) {
	refile, _ := filepath.Abs(file)
	yamlFile, err := ioutil.ReadFile(refile)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return
	}

	c.Rules.RateMaxTime = 0
	for _, v := range c.Rules.Rate {
		if c.Rules.RateMaxTime < v.Time {
			c.Rules.RateMaxTime = v.Time
		}
	}

	return
}
