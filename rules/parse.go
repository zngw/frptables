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
