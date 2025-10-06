package httpv1

import "github.com/gin-gonic/gin"

func SetUpRoutes(v1 *gin.RouterGroup) {
	auth := v1.Group("/auth")
	{
		auth.POST("/login")
	}
}
