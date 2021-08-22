// @Title
// @Description $
// @Author  55
// @Date  2021/8/22
package rules

import "github.com/zngw/log"

func Check(text string) {
	err, ip, name, port := parse(text)
	if err != nil {
		if err.Error() != "not tcp link" {
			log.Trace("net", err.Error())
		}

		return
	}

	check(ip, name, port)

	return
}
