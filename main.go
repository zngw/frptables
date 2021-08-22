// @Title
// @Description $
// @Author  55
// @Date  2021/8/21
package main

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/hpcloud/tail"
	"github.com/zngw/frptables/config"
	"github.com/zngw/frptables/rules"
	"github.com/zngw/log"
)

func main() {
	// 获取程序自身路径
	dir, file := filepath.Split(os.Args[0])

	// 读取命令行配置文件参数
	c := flag.String("c", dir+"/config.yml", "默认配置为 config.yml")
	flag.Parse()

	// 初始化日志
	err := log.Init(dir+"/logs/"+file, []string{"add", "link", "net"})
	if err != nil {
		panic(err)
	}

	// 初始化配置
	err = config.Cfg.Init(*c)
	if err != nil {
		panic(err)
	}

	// 初始化规则
	rules.Init()

	// 启动用tail监听
	tails, err := tail.TailFile(config.Cfg.FrpsLog, tail.Config{
		ReOpen:    true,                                 // 重新打开
		Follow:    true,                                 // 是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 从文件的哪个地方开始读
		MustExist: false,                                // 文件不存在不报错
		Poll:      true,
	})

	if err != nil {
		log.Error("tail file failed, err:", err)
		return
	}
	var line *tail.Line
	var ok bool

	for {
		line, ok = <-tails.Lines
		if !ok {
			log.Error("tail file close reopen, filename:%s\n", tails.Filename)
			time.Sleep(time.Second)
			continue
		}

		rules.Check(line.Text)
	}
}
