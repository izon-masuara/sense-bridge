package main

import (
	"server_lidm/db"
	"server_lidm/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db.ConnectDB("mongodb://mongo:27017")

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"Origin", "Content-Type"},
		ExposeHeaders: []string{"Content-Length"},
		// AllowCredentials: true,
		// MaxAge:           12 * time.Hour,
	}))

	r.GET("/start", handlers.GenerateCode)
	r.GET("/start/:code", handlers.FindCode)
	r.POST("/audio/:sec/:code", handlers.UploadAudio)
	r.POST("/image/:code", handlers.UploadImage)
	r.GET("/summary/:code", handlers.Summary)
	r.Run(":8080")
}
