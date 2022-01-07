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
	"fmt"
	"io/ioutil"
	"net/http"
)

// 配置文件结构体
type IpInfo struct {
	Status     string `json:"status,omitempty"`
	Country    string `json:"country,omitempty"`
	RegionName string `json:"regionName,omitempty"`
	Region     string `json:"region,omitempty"`
	City       string `json:"city,omitempty"`
	Isp        string `json:"isp,omitempty"`
	Query      string `json:"query,omitempty"`
}

var ips = make(map[string]*IpInfo)

func GetIpInfo(ip string) (info *IpInfo) {
	if v, ok := ips[ip]; ok {
		info = v
		return
	}

	// 默认使用淘宝ip地址库
	info = getIpInfoTaoBao(ip)

	// 超出访问频率，改用 ipip.net 再请求一次
	if info.Status == "fail" {
		info = getIpInfoByIpIp(ip)
	}

	// 超出访问频率，改用 ip-api 再请求一次
	if info.Status == "fail" {
		info = getIpInfoByIpApi(ip)
	}

	return
}

// 淘宝ip地址库
// 限制频率:每个用户的访问频率需小于1qps
func getIpInfoTaoBao(ip string) (info *IpInfo) {
	info = new(IpInfo)
	info.Status = "fail"

	url := fmt.Sprintf("https://ip.taobao.com/outGetIpInfo?ip=%s&accessKey=alibaba-inc", ip)
	resp, err := http.Get(url)
	if err != nil {
		// 读取网页数据错误
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// 读取网页数据错误
		return
	}

	if resp.StatusCode == 200 {
		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			// 网页解析错误
			return
		}

		v,ok :=result["code"]
		if !ok || int(v.(float64)) != 0 {
			// 数据错误
			return
		}

		d := result["data"]
		data, err := json.Marshal(d)
		if err != nil {
			return
		}

		err = json.Unmarshal(data,&info)
		if err != nil {
			return
		}

		info.Status = "success"
		info.RegionName = info.Region
		info.Query = ip
	}

	return
}

// 通过ipip.net获取ip地理位置信息
// 有访问频率限制
func getIpInfoByIpIp(ip string) (info *IpInfo) {
	info = new(IpInfo)
	info.Status = "fail"

	url := fmt.Sprintf("http://freeapi.ipip.net/%s", ip)
	resp, err := http.Get(url + ip)
	if err != nil {
		// 读取网页数据错误
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// 读取网页数据错误
		return
	}

	if resp.StatusCode == 200 {
		var result []string
		err = json.Unmarshal(body, &result)
		if err != nil {
			// 网页解析错误
			return
		}

		if len(result) == 5 {
			info.Status = "success"
			info.Country = result[0]
			info.RegionName = result[1]
			info.City = result[2]
			info.Isp = result[4]
			info.Query = ip
		}
	}

	return
}

// 通过ip-api.com获取ip地理位置信息
// 由于ip-api.com是国外的网站，对国内市级ip位置有一定误差
func getIpInfoByIpApi(ip string) (info *IpInfo) {
	info = new(IpInfo)
	info.Status = "fail"

	url := fmt.Sprintf("http://ip-api.com/json/%s?lang=zh-CN", ip)
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
