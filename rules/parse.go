// @Title
// @Description $
// @Author  55
// @Date  2021/8/22
package rules

import (
	"fmt"
	"github.com/zngw/frptables/config"
	"strings"
)

func parse(text string) (err error, ip, name string, port int) {
	// 从frp日志中获取tcp连接信息
	// 2021/08/21 15:35:29 [I] [proxy.go:162] [f1aec30e84827422] [ZNGW] get a user connection [210.0.159.76:32832]
	if !strings.Contains(text, "get a user connection") {
		err = fmt.Errorf("not tcp link")
		return
	}

	// 获取ip
	begin := strings.LastIndex(text, "[")
	end := strings.LastIndex(text, "]")
	if begin < 0 || end <= begin {
		err = fmt.Errorf("formt error")
		return
	}

	linker := text[begin+1 : end]
	linkerArray := strings.Split(linker, ":")
	if len(linkerArray) != 2 {
		err = fmt.Errorf("formt error")
		return
	}

	ip = linkerArray[0]

	// 获取 Name
	tmp := text[0:begin]
	begin = strings.LastIndex(tmp, "[")
	end = strings.LastIndex(tmp, "]")
	if begin < 0 || end <= begin {
		err = fmt.Errorf("formt error")
		return
	}

	name = text[begin+1 : end]

	if v, ok := config.Cfg.NamePort[name]; ok {
		port = v
	} else {
		port = -1
	}

	return
}
