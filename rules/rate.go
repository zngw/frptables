// @Title
// @Description $
// @Author  55
// @Date  2021/8/22
package rules

import (
	"fmt"
	"github.com/zngw/frptables/config"
	"sync"
	"time"
)

// 用线程安装的map保存ip记录
var ips = sync.Map{}

type history struct {
	Ip   string
	Port int
	List []int64
}

func (h *history) Init(ip string, port int) {
	h.Ip = ip
	h.Port = port
	h.Add()
}

func (h *history) Add() {
	h.List = append(h.List, time.Now().Unix())
}

func (h *history) Count(interval int64) (count int) {
	count = 0
	now := time.Now().Unix()
	for _, v := range h.List {
		if now < v+interval {
			count++
		}
	}
	return
}

func rateInit() {
	// 异步清理
	go func() {
		for {
			delTime := time.Now().Unix() - config.Cfg.Rules.RateMaxTime
			var del []string
			ips.Range(func(k, v interface{}) bool {
				h := v.(*history)
				for i := 0; i < len(h.List); {
					if h.List[i] < delTime {
						h.List = append(h.List[:i], h.List[i+1:]...)
					} else {
						i++
					}
				}

				if len(h.List) == 0 {
					del = append(del, k.(string))
				}

				return true
			})

			for _, k := range del {
				ips.Delete(k)
			}

			time.Sleep(time.Second * time.Duration(config.Cfg.Rules.RateMaxTime))
		}
	}()
}

func getIpHistory(ip string, port int) (h *history) {
	key := fmt.Sprintf("%s:%d", ip, port)
	tmp, ok := ips.Load(key)
	if !ok {
		h = new(history)
		h.Init(ip, port)
		ips.Store(key, h)
		return
	}

	h = tmp.(*history)
	h.Add()
	return
}

func delIpHistory(ip string, port int) {
	key := fmt.Sprintf("%s:%d", ip, port)
	ips.Delete(key)
}
