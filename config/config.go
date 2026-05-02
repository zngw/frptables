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

package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

var cfgFile string
var Cfg Conf

// 配置文件结构体
type Conf struct {
	FrpsLog     string         `yaml:"frps_log,omitempty"`    // 监听的frps日志文件
	Logs        string         `yaml:"logs,omitempty"`        // 自身日志文件
	NamePort    map[string]int `yaml:"name_port,omitempty"`   // 名字端口对应表
	TablesType  string         `yaml:"tables_type,omitempty"` // 启用防火墙类型
	AllowIp     []string       `yaml:"allow_ip,omitempty"`    // ip白名单
	AllowPort   []int          `yaml:"allow_port,omitempty"`  // 端口白名单
	Rules       []CfgRules     `yaml:"rules,omitempty"`       // 规则访问
	RateMaxTime int64          `yaml:"-"`                     // IP频率中最高超时时间
}

type CfgRules struct {
	Port       int    `yaml:"port,omitempty"`       // 端口
	Country    string `yaml:"country,omitempty"`    // 国家
	RegionName string `yaml:"regionName,omitempty"` // 省
	City       string `yaml:"city,omitempty"`       // 城市
	Rules      string `yaml:"rules,omitempty"`      // 规则
	Time       int64  `yaml:"time,omitempty"`       // 时间
	Count      int    `yaml:"count,omitempty"`      // 次数
}

func (c *Conf) Load(file string) (err error) {
	refile, _ := filepath.Abs(file)
	yamlFile, err := ioutil.ReadFile(refile)
	if err != nil {
		return
	}

	// 临时序列化一次，校验配置语法。
	// 如果语法有问题不会对现有配置不影响。
	var tmp Conf
	err = yaml.Unmarshal(yamlFile, &tmp)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return
	}

	c.RateMaxTime = 0
	for _, v := range c.Rules {
		v.RegionName = strings.TrimSuffix(v.RegionName, "省")
		v.City = strings.TrimSuffix(v.City, "市")
		if c.RateMaxTime < v.Time {
			c.RateMaxTime = v.Time
		}
	}

	return
}

// Init 初始化配置
func Init(file string) (err error) {
	cfgFile = file
	err = Cfg.Load(cfgFile)
	return
}

// ReloadConfig 重新加载配置（供 signal handler 调用）
func ReloadConfig() {
	err := Cfg.Load(cfgFile)
	if err != nil {
		log.Printf("Failed to reload config: %v\n", err)
		return
	}
	log.Println("Config reloaded successfully")
}
