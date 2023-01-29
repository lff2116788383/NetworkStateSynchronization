package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Redis struct {
	IP   string `json:"ip" `
	Port int    `json:"port" `
}
type Mysql struct {
	IP           string `json:"ip" `
	Port         int    `json:"port" `
	User         string `json:"user" `
	Password     string `json:"password" `
	Databasename string `json:"databasename" `
}

type DBConfig struct {
	Redis *Redis `json:"redis" `
	Mysql *Mysql `json:"mysql" `
}

type ServerConfig struct {
	ServerId      int       `json:"serverid" `
	Host          string    `json:"host" `
	LogDir 		  string    `json:"logdir"` 	// 日志存储路径
	LocalSavePath string    `json:"localsavepath"` //! 本地存储路径
	DBConfig      *DBConfig `json:"databases" `
	HttpHost      string    `json:"httphost" `
}

var GlobalConfig = ServerConfig{} //服务配置

func (c *ServerConfig) InitConfig(configPath string) bool {
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Println("error")
		return false
	}
	err = json.Unmarshal(configFile, &c)
	if err != nil {
		fmt.Println("error")
		return false
	}
	return true
}
