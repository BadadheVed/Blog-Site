package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/yourname/blog-kafka/function"
	middleware "github.com/yourname/blog-kafka/middlewares"
)

func ChannelRoutes(r *gin.Engine) {
	ch := r.Group("/channels")
	ch.Use(middleware.AuthMiddleware())
	{
		ch.POST("/", function.CreateChannel)
		ch.POST("/:id/subscribe", function.SubscribeChannel)

	}
}
