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

package main

import (
	"encoding/json"
	"flag"
	"github.com/hpcloud/tail"
	"github.com/zngw/frptables/config"
	"github.com/zngw/frptables/rules"
	"github.com/zngw/log"
	"github.com/zngw/zipinfo/ipinfo"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// 读取命令行配置文件参数
	c := flag.String("c", "./config.yml", "默认配置为 config.yml")
	s := flag.String("s", "", "默认配置为空")
	flag.Parse()

	if *s == "reload" {
		// 如果是reload，发送reload指令后退出
		config.SendReload()
		return
	}

	// 初始化配置
	err := config.Init(*c)
	if err != nil {
		panic(err)
	}

	// 初始化日志
	_, file := filepath.Split(os.Args[0])
	logFile, _ := filepath.Abs(config.Cfg.Logs)
	logFile = filepath.Join(config.Cfg.Logs, file)
	log.InitLog("all", logFile, "trace", 7, true, []string{"add", "link", "net", "sys"})

	var ipCfg []interface{}
	err = json.Unmarshal([]byte(config.Cfg.IpInfo), &ipCfg)
	if err == nil {
		ipinfo.Init(ipCfg)
	}

	// 初始化规则
	rules.Init()

	// 启动用tail监听
	frpLog, _ := filepath.Abs(config.Cfg.FrpsLog)
	tails, err := tail.TailFile(frpLog, tail.Config{
		ReOpen:    true,                                 // 重新打开
		Follow:    true,                                 // 是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件的哪个地方开始读
		MustExist: false,                                // 文件不存在不报错
		Poll:      true,
	})

	if err != nil {
		log.Error("sys", "tail file failed, err:%v", err)
		return
	}

	log.Trace("sys", "frptables 已启动，正在监听日志文件：%s", frpLog)
	var line *tail.Line
	var ok bool

	for {
		line, ok = <-tails.Lines
		if !ok {
			log.Error("sys", "tail file close reopen, filename:%s\n", tails.Filename)
			time.Sleep(time.Second)
			continue
		}

		rules.Check(line.Text)
	}
}
