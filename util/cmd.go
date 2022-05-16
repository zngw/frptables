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
	"github.com/zngw/log"
	"io/ioutil"
	"os/exec"
	"runtime"
)

func Command(arg ...string) (result string) {
	name := "/bin/bash"
	c := "-c"
	if runtime.GOOS == "windows" {
		name = "cmd"
		c = "/C"
	}
	arg = append([]string{c}, arg...)
	cmd := exec.Command(name, arg...)

	//创建获取命令输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Error("sys","Error:can not obtain stdout pipe for command:%v\n", err)
		return
	}

	//执行命令
	if err := cmd.Start(); err != nil {
		log.Error("sys","Error:The command is err,%v", err)
		return
	}

	//读取所有输出
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Error("sys","ReadAll Stdout:%v", err.Error())
		return
	}

	if err := cmd.Wait(); err != nil {
		log.Error("wait:", err.Error())
		return
	}

	result = string(bytes)
	return
}
