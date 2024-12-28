# chord-backend

API Document: [API Doc](API.md)

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
    port := 21776

    gin.SetMode(gin.ReleaseMode)
    r := gin.New()
    r.Use(gin.Logger())
    r.Use(gin.Recovery())
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost", "http://127.0.0.1"},
        AllowMethods:     []string{"GET", "POST"},
        AllowHeaders:     []string{"Content-Type"},
        AllowCredentials: true,
    }))

    router.SetupAPIRouter("api", r)

    log.Println("Starting server on", port)
    if err := r.Run("localhost:" + strconv.Itoa(port)); err != nil {
        log.Fatalf("Server stopped with error: %v", err)
    }
}
```
