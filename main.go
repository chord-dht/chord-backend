package main

import (
	"chord-backend/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/new", handlers.CreateNode)
	r.GET("/initialize", handlers.InitializeNode)
	r.GET("/quit", handlers.QuitNode)
	r.GET("/printstate", handlers.GetNodeState)
	r.POST("/lookup", handlers.Lookup)
	r.POST("/storefile", handlers.StoreFile)
	r.POST("/getfile", handlers.GetFile)
	r.POST("/downloadfile", handlers.DownloadFile)

	r.Run(":8080")
}
