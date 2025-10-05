package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/yourname/blog-kafka/function"
	middleware "github.com/yourname/blog-kafka/middlewares"
)

func ChannelRoutes(r *gin.Engine) {
	ch := r.Group("/channel")
	ch.Use(middleware.AuthMiddleware())
	{
		ch.POST("/create", function.CreateChannel)
		ch.POST("/:id/subscribe", function.SubscribeChannel)
	}
}
