package main

import (
	"chord-backend/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	r.POST("/new", handlers.CreateNode)
	r.GET("/initialize", handlers.InitializeNode)
	r.GET("/quit", handlers.QuitNode)
	r.GET("/printstate", handlers.GetNodeState)
	r.POST("/storefile", handlers.StoreFile)
	r.POST("/getfile", handlers.GetFile)
	r.POST("/downloadfile", handlers.DownloadFile)

	r.Run(":8080")
}
