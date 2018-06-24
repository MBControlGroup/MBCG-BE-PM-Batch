package routes

import (
	"controllers"

	"github.com/gin-gonic/gin"
)

// Download download batch
func Download(ctx *gin.Context) {
	controllers.Download(ctx)
}
