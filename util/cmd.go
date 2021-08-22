// @Title
// @Description $
// @Author  55
// @Date  2021/8/23
package util

import (
	"github.com/zngw/log"
	"io/ioutil"
	"os/exec"
)

func Command(name string, arg ...string) (result string) {
	cmd := exec.Command(name, arg...)

	//创建获取命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Error("Error:can not obtain stdout pipe for command:%s\n", err)
		return
	}

	//执行命令
	if err := cmd.Start(); err != nil {
		log.Error("Error:The command is err,", err)
		return
	}

	//读取所有输出
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Error("ReadAll Stdout:", err.Error())
		return
	}

	if err := cmd.Wait(); err != nil {
		log.Error("wait:", err.Error())
		return
	}

	result = string(bytes)
	return
}
