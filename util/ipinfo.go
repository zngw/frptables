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

package util

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// 配置文件结构体
type IpInfo struct {
	Status      string `json:"status,omitempty"`
	Country     string `json:"country,omitempty"`
	CountryCode string `json:"countryCode,omitempty"`
	Region      string `json:"region,omitempty"`
	RegionName  string `json:"regionName,omitempty"`
	City        string `json:"city,omitempty"`
	Zip         string `json:"zip,omitempty"`
	Lat         string `json:"lat,omitempty"`
	Lon         string `json:"lon,omitempty"`
	Timezone    string `json:"timezone,omitempty"`
	Isp         string `json:"isp,omitempty"`
	Org         string `json:"org,omitempty"`
	As          string `json:"as,omitempty"`
	Query       string `json:"query,omitempty"`
}

var ips = make(map[string]*IpInfo)

func GetIpInfo(ip string) (info *IpInfo) {
	if v, ok := ips[ip]; ok {
		info = v
		return
	}

	info = new(IpInfo)
	info.Status = "fail"

	const url = "http://ip-api.com/json/"
	resp, err := http.Get(url + ip + "?lang=zh-CN")
	if err != nil {
		// 获取不到地理位置，
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// 读取网页数据错误
		return
	}
	if resp.StatusCode == 200 {
		err = json.Unmarshal(body, &info)
		if err != nil {
			// 网页解析错误
			return
		}
	}

	ips[ip] = info

	return
}
