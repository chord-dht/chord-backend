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
	r.GET("/printstate", handlers.PrintNodeState)
	r.POST("/storefile", handlers.StoreFile)
	r.POST("/getfile", handlers.GetFile)

	r.Run(":8080")
}
