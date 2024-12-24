# chord-backend

Example Usage:

```go
package main

import (
    "log"

    "github.com/chord-dht/chord-backend/router"

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

    router.SetupRouter(r)

    port := ":9000"
    log.Println("Starting server on", port)
    if err := r.Run(port); err != nil {
        log.Fatalf("Server stopped with error: %v", err)
    }
}
```
