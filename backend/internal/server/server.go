package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func GetRouter() *gin.Engine {
	router := gin.Default()

	config := cors.DefaultConfig()

	// TODO: Set a less permissive "allow-origins". Potentially even configure via an setting.
	config.AllowOrigins = []string{"*"}

	router.Use(cors.New(config))

	router.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	return router
}
