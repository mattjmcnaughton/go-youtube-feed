package server

import (
	"context"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/mattjmcnaughton/go-youtube-feed/internal/youtube"
)

func GetRouter(youtubeClient *youtube.YoutubeClient) *gin.Engine {
	router := gin.Default()

	config := cors.DefaultConfig()
	// TODO: Set a less permissive "allow-origins". Potentially even configure via an setting.
	config.AllowOrigins = []string{"*"}
	router.Use(cors.New(config))

	router.GET("/status", statusGET)

	v1 := router.Group("v1")
	v1.GET("feed/:handle", feedGETClosure(youtubeClient))

	return router
}

func statusGET(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func feedGETClosure(youtubeClient *youtube.YoutubeClient) func(c *gin.Context) {
	feedGET := func(c *gin.Context) {
		ctx := context.Background()
		feedURL, err := youtubeClient.GenerateAtomFeedURL(ctx, c.Param("handle"))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"feed": gin.H{"url": feedURL}})
	}

	return feedGET
}
