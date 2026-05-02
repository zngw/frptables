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
	"regexp"
	"strings"

	"github.com/zngw/frptables/config"
)

// 解析日志
func parse(text string) (err error, ip, name string, port int) {
	// 从frp日志中获取tcp连接信息
	// 2026-05-02 00:50:55.239 [I] [proxy/proxy.go:236] [5e6adf3e09951a2f] [家电脑.TCP-家] get a user connection [39.182.4.210:26652]
	// 2026-05-02 01:11:16.477 [I] [proxy/proxy.go:236] [f1a4c9ca8a1e69b2] [小电脑.TCP-小] get a user connection [[2409:8a28:14e9:740:2905:c2ba:6c98:3ec5]:4754]
	// IPv4: ... [39.182.4.210:26652]
	// IPv6: ... [[2409:8a28:14e9:740:2905:c2ba:6c98:3ec5]:4754]
	if !strings.Contains(text, "get a user connection") {
		err = fmt.Errorf("not tcp link")
		return
	}

	// 正则表达式获取转发名和地址
	// 日志格式: ... [name] get a user connection [ip:port] 或 [[ipv6]:port]
	compileRegex := regexp.MustCompile(`\[I\] \[.*?\] \[.*?\] \[(.*?)\] get a user connection (.*)`)
	matchArr := compileRegex.FindStringSubmatch(text)

	if len(matchArr) <= 2 {
		err = fmt.Errorf("not tcp link")
		return
	}

	// 转发名
	name = matchArr[1]
	addr := strings.TrimSpace(matchArr[2]) // "[39.182.4.210:26652]" 或 "[[2409:8a28:14e9:740:2905:c2ba:6c98:3ec5]:4754]"

	// 解析 IP 地址（支持 IPv4 和 IPv6）
	ip, err = extractIP(addr)
	if err != nil {
		return
	}

	if v, ok := config.Cfg.NamePort[name]; ok {
		port = v
	} else {
		port = -1
	}

	return
}

// extractIP 从地址字符串中提取 IP（支持 IPv4 和 IPv6）
// IPv4 格式: [ip:port]
// IPv6 格式: [[ipv6]:port]
func extractIP(addr string) (string, error) {
	// 尝试 IPv6 格式：[[ipv6]:port]
	ipv6Regex := regexp.MustCompile(`^\[\[([^\]]+)\]:\d+\]$`)
	if match := ipv6Regex.FindStringSubmatch(addr); len(match) > 0 {
		return match[1], nil
	}

	// 尝试 IPv4 格式：[ip:port]
	ipv4Regex := regexp.MustCompile(`^\[([^:]+):\d+\]$`)
	if match := ipv4Regex.FindStringSubmatch(addr); len(match) > 0 {
		return match[1], nil
	}

	return "", fmt.Errorf("invalid addr: %s", addr)
}
