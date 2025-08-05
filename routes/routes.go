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

		appGroup := api.Group("/app")
		{
			// 一言接口
			appGroup.GET("/hitokoto", handlers.GetHitokoto)
		}

		// 需要JWT认证的路由
		secureGroup := api.Group("/")
		secureGroup.Use(middleware.JWTAuth())
		{
			// 获取所有活跃用户
			secureGroup.GET("/getUsers", handlers.GetUsers)
			// 获取所有失效用户
			secureGroup.GET("/getInactiveUsers", handlers.GetInactiveUsers)
			// 获取当前登录用户
			secureGroup.GET("/getUser", handlers.GetCurrentUser)
			// 获取指定id的用户
			secureGroup.GET("/getUser/:userId", handlers.GetUserByID)
			// 创建用户
			secureGroup.POST("/createUser", handlers.CreateUser)
			// 更新指定id用户状态
			secureGroup.POST("/updateUserStatus/:userId", handlers.UpdateUserStatus)
			// 更新当前登录用户信息
			secureGroup.POST("/updateUser", handlers.UpdateCurrentUser)
			// 更新指定id用户信息
			secureGroup.POST("/updateUser/:userId", handlers.UpdateUser)
			// 删除指定id用户
			secureGroup.DELETE("/deleteUser/:userId", handlers.DeleteUser)
			// 其他需要保护的路由...
		}
	}
}
