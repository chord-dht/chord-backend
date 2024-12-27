package router

import (
	"github.com/chord-dht/chord-backend/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	r.GET("/nodestatus", handlers.NodeStatus)
	r.POST("/new", handlers.CreateNode)
	r.GET("/initialize", handlers.InitializeNode)
	r.GET("/quit", handlers.QuitNode)
	r.GET("/printstate", handlers.GetNodeState)
	r.POST("/storefile", handlers.StoreFile)
	r.POST("/getfile", handlers.GetFile)
	r.POST("/downloadfile", handlers.DownloadFile)
}
