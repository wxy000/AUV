package routes

import (
	"AUV/handlers"
	"AUV/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		// 公开路由
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/login", handlers.Login)
			authGroup.POST("/refresh", handlers.RefreshToken)
		}

		// 一言接口
		api.GET("/hitokoto", handlers.GetHitokoto)

		// 需要JWT认证的路由
		secureGroup := api.Group("/")
		secureGroup.Use(middleware.JWTAuth())
		{
			secureGroup.GET("/getUsers", handlers.GetUsers)
			secureGroup.GET("/getInactiveUsers", handlers.GetInactiveUsers)
			secureGroup.GET("/getUser/:userId", handlers.GetUser)
			secureGroup.POST("/createUser", handlers.CreateUser)
			secureGroup.POST("/updateUserStatus/:userId", handlers.UpdateUserStatus)
			secureGroup.POST("/updateUser/:userId", handlers.UpdateUser)
			secureGroup.DELETE("/deleteUser/:userId", handlers.DeleteUser)
			// 其他需要保护的路由...
		}
	}
}
