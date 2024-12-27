package router

import (
	"path"

	"github.com/chord-dht/chord-backend/handlers"

	"github.com/gin-gonic/gin"
)

func SetupAPIRouter(prefixPath string, r *gin.Engine) {
	basePath := "/" + prefixPath

	r.GET(path.Join(basePath, "nodestatus"), handlers.NodeStatus)
	r.POST(path.Join(basePath, "new"), handlers.CreateNode)
	r.GET(path.Join(basePath, "initialize"), handlers.InitializeNode)
	r.GET(path.Join(basePath, "quit"), handlers.QuitNode)
	r.GET(path.Join(basePath, "printstate"), handlers.GetNodeState)
	r.POST(path.Join(basePath, "storefile"), handlers.StoreFile)
	r.POST(path.Join(basePath, "getfile"), handlers.GetFile)
	r.POST(path.Join(basePath, "downloadfile"), handlers.DownloadFile)
}
