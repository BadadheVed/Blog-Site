package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yourname/blog-kafka/function"
	middleware "github.com/yourname/blog-kafka/middlewares"
)

func blogRouter(r *gin.Engine) {
	blog := r.Group("/blog")
	{
		blog.POST("/:channelId/create", middleware.AuthMiddleware(), function.CreateBlog)
		blog.PUT("/:id/edit", middleware.AuthMiddleware(), function.EditBlog)
		blog.DELETE("/:id/delete", middleware.AuthMiddleware(), function.DeleteBlog)
		blog.GET("/", middleware.AuthMiddleware(), function.GetBlogs)
	}
}
