package controllers

import (
	"fmt"
	"log"
	"net/http"
	"protocol"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

// CrashHandler handle crash
func CrashHandler(ctx *gin.Context) {
	if err := recover(); err != nil {
		log.Println(err)
		debug.PrintStack()

		res := &protocol.ResponseMsg{
			Code:  http.StatusInternalServerError,
			EnMsg: "server crash",
			CnMsg: "服务器内部出错",
			Data:  nil,
		}
		ctx.JSON(res.Code, res)
	}
}

// GetFilenameWithTimestamp 获取以时间戳命名的文件名
func GetFilenameWithTimestamp() string {
	return fmt.Sprintf("/var/tmp/%d.xlsx", time.Now().UnixNano())
}
