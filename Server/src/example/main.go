package main

import (
	"Src/libs/config"
	"Src/libs/db"
	"Src/libs/logger"
	"Src/libs/mmo"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"


)

func main() {

	//是否以守护进程启动 参数-d
	args := os.Args
	daemon := false
	for k, v := range args {
		if v == "-d" {
			daemon = true
			args[k] = ""
		}
	}

	if daemon {
		Daemonize(args...)
		return
	}

	//参数解析及提示
	if len(os.Args) < 2 {

		fmt.Println("Args Error")
		fmt.Printf("USAGE: %s config.json\n", os.Args[0])
		os.Exit(1) //退出os 状态码0表示成功，非0表示出错。
	}

	//全局配置初始化
	if !config.GlobalConfig.InitConfig(os.Args[1]) {
		time.Sleep(2 * time.Second) //等两秒钟，不然还没写入日志，程序就结束了
		return
	}

	//日志初始化
	logDir := config.GlobalConfig.LogDir
	fmt.Println("logdir",logDir)
	if !logger.Init("flow.log", logDir, 40*1024*1024, 10, logger.ERROR_LEVEL, false, logger.WRITE_FILE) {
		fmt.Println("Logger Init Fail")
		return
	}
	logger.Error("Initialization Start")

	logger.Error("Logger Init Succ")

	logger.Error("Load Config File:%s Succ", os.Args[1])
	//数据库初始化 10.0.4.7数据库只能内网访问
	if !db.Db.InitMysql() {

		logger.Error("MysqlDB Init Fail")
		time.Sleep(2 * time.Second) //等两秒钟，不然还没写入日志，程序就结束了
		return
	}
	logger.Error("MysqlDB Init Succ,IP:%s", config.GlobalConfig.DBConfig.Mysql.IP)

	// 随机种子初始化 种子可变
	rand.Seed(time.Now().UnixNano())

	//GameServer 初始化
	if !mmo.GetGameServer().Init() {
		logger.Error("mmo GameServer Init Fail")
		time.Sleep(3 * time.Second) //等两秒钟，不然还没写入日志，程序就结束了
		return
	}
	logger.Error("GameServer Init Succ")
	logger.Error("Initialization Completed")

	//创建pid文件
	Pidfile := fmt.Sprintf("%s.pid", os.Args[0])
	file, err := os.Create(Pidfile)
	if err != nil {
		fmt.Println("create pid file fail")
	}
	writeString := fmt.Sprintf("%d\n", os.Getpid())
	file.WriteString(writeString)
	defer func() {
		// 2.关闭文件
		file.Close()
		fmt.Printf("file close")
	}()

	mmo.GetGameServer().Run()

}

func Daemonize(args ...string) {
	var arg []string
	if len(args) > 1 {
		arg = args[1:]
	}
	cmd := exec.Command(args[0], arg...)
	cmd.Env = os.Environ()
	cmd.Start()
}
