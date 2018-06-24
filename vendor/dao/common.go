package dao

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	MySQL "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	// gorm mysql connection config
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// MyDB
var MyDB *gorm.DB

// RecordNotFoundErrMsg record not found error message
const RecordNotFoundErrMsg = "record not found"

const configFilename = "./config.json"

type dbConfig struct {
	User   string            `json:"user"`
	Passwd string            `json:"passwd"`
	DBName string            `json:"db_name"`
	Net    string            `json:"net"`
	IPAddr string            `json:"ip_addr"`
	Port   string            `json:"port"`
	Params map[string]string `json:"params"`
}

type myConfig struct {
	DBCfg dbConfig `json:"db_cfg"`
}

var globalConfig myConfig
var dbAddr string

func init() {
	log.Println("init dao ...")
	var err error

	// 加载配置
	configData, err := ioutil.ReadFile(configFilename)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(configData, &globalConfig)
	if err != nil {
		log.Fatal(err)
	}

	dbAddr = fmt.Sprintf("%s:%s", globalConfig.DBCfg.IPAddr, globalConfig.DBCfg.Port)

	cfg := &MySQL.Config{
		User:                 globalConfig.DBCfg.User,
		Passwd:               globalConfig.DBCfg.Passwd,
		DBName:               globalConfig.DBCfg.DBName,
		Net:                  globalConfig.DBCfg.Net,
		Addr:                 dbAddr,
		Params:               globalConfig.DBCfg.Params,
		Loc:                  time.Local,
		ParseTime:            true,
		AllowNativePasswords: true,
		ReadTimeout:          31536000 * time.Second,
		WriteTimeout:         31536000 * time.Second,
	}

	// 连接DB
	MyDB, err = gorm.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	MyDB.DB().SetConnMaxLifetime(31536000 * time.Second)
}
