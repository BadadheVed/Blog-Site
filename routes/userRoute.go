package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yourname/blog-kafka/function"
	middleware "github.com/yourname/blog-kafka/middlewares"
)

func userRouter(r *gin.Engine) {
	user := r.Group("/user")
	{

		user.GET("/notifications", middleware.AuthMiddleware(), function.GetUserNotifications)
	}
}
