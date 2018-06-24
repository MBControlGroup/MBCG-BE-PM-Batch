package routes

import (
	"controllers"

	"github.com/gin-gonic/gin"
)

// Upload upload batch
func Upload(ctx *gin.Context) {
	controllers.Upload(ctx)
}
