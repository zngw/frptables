// @Title
// @Description $
// @Author  55
// @Date  2021/8/22
package tables

import (
	"github.com/zngw/log"
)

func Check(text string) {
	err, ip, name, port := parse(text)
	if err != nil {
		log.Trace("net", err.Error())
		return
	}

	check(ip, name, port)

	return
}
