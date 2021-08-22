// @Title
// @Description $
// @Author  55
// @Date  2021/8/22
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

func GetIpInfo(ip string) (info IpInfo) {
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

	return
}
