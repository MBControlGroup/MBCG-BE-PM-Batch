package main

import (
	"dao"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"routes"

	"github.com/gin-contrib/cors"

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

// CorsHeadersMiddleWare 跨域头设置中间件
func CorsHeadersMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		ctx.Next()
	}
}

func main() {
	defer func() {
		dao.MyDB.Close()
	}()

	router := gin.Default()
	router.Use(cors.Default())
	router.Use(CorsHeadersMiddleWare()) // 允许所有源的cors配置

	// OPTIONS 跨域配置
	router.OPTIONS("/*action", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, struct{}{})
	})

	// 批量导入
	router.POST("/upload", routes.Upload)

	// 批量导出
	router.GET("/download", routes.Download)

	router.Run(listenAddr)
}
