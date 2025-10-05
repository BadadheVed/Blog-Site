package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yourname/blog-kafka/function"
	middleware "github.com/yourname/blog-kafka/middlewares"
)

func authRouter(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/signup", function.Signup)
		auth.POST("/login", function.Login)
		auth.GET("/users", middleware.AuthMiddleware(), function.GetAllUsers)
	}
}
