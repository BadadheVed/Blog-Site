package routes

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	r := gin.Default()

	authRouter(r)

	ChannelRoutes(r)
	blogRouter(r)
	userRouter(r)
	return r
}
