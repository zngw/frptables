package util

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/zngw/golib/log"
)

func GetIpInfo(ip string) (ok bool, Country, Region, City string) {
	// 通过 https://ip.zengwu.com.cn?ip= 接口获取ip归属地信息
	// 这里也可以切换成自己的ip查询
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://ip.zengwu.com.cn?ip="+ip, nil)
	resp, err := client.Do(req)
	if err != nil {
		// 读取网页数据错误
		log.Error(err.Error())
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// 读取网页数据错误
		log.Error(err.Error())
		return
	}

	if resp.StatusCode != 200 {
		log.Error("http status code = %d", resp.StatusCode)
		return
	}

	var jsonInfo struct {
		Result   int32  `json:"result,omitempty"`   // 状态
		Country  string `json:"country,omitempty"`  // 国家
		Province string `json:"province,omitempty"` // 省
		City     string `json:"city,omitempty"`     // 城市
		Isp      string `json:"isp,omitempty"`      // 运营商
		Query    string `json:"query,omitempty"`    // 查询IP
		Time     int64  `json:"time,omitempty"`     // 查询时间
	}

	err = json.Unmarshal(body, &jsonInfo)
	if err != nil {
		log.Error(err.Error())
		return
	}

	if jsonInfo.Result != 0 {
		log.Error("get ip info err = %d", jsonInfo.Result)
		return
	}

	return true, jsonInfo.Country, jsonInfo.Province, jsonInfo.City
}
