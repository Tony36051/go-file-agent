package routers

import (
	"github.com/Tony36051/go-file-agent/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/download/:filename", handlers.DownloadFile)
}
