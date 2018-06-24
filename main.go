package main

import (
	"dao"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"routes"

	"github.com/gin-gonic/gin"
)

const configFilename = "./config.json"

type serverConfig struct {
	IPAddr string `json:"ip_addr"`
	Port   string `json:"port"`
}

type myConfig struct {
	ServerCfg serverConfig `json:"server_cfg"`
}

var globalCfg myConfig
var listenAddr string

func init() {
	log.SetFlags(log.Llongfile | log.LstdFlags)

	// 加载配置信息
	configData, err := ioutil.ReadFile(configFilename)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(configData, &globalCfg)
	if err != nil {
		log.Fatal(err)
	}

	listenAddr = fmt.Sprintf("%s:%s", globalCfg.ServerCfg.IPAddr, globalCfg.ServerCfg.Port)
}

func main() {
	defer func() {
		dao.MyDB.Close()
	}()

	router := gin.Default()

	// 批量导入
	router.POST("/upload", routes.Upload)

	// 批量导出
	router.GET("/download", routes.Download)

	router.Run(listenAddr)
}
