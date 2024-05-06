package router

import (
	"ginchat/docs"
	"ginchat/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	//SWAG
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	//首页
	r.GET("/index", service.GetIndex)
	r.GET("/", service.GetIndex)
	r.GET("/toRegister", service.ToRegister)
	r.GET("/toChat", service.ToChat)
	r.POST("/searchFriends", service.SearchFriends)
	r.GET("/chat", service.Chat)

	// 静态资源
	r.Static("/asset", "asset/")
	r.LoadHTMLGlob("views/**/*")

	//用户模块
	r.GET("/user/getUserList", service.GetUserList)
	r.POST("/user/createUser", service.CreateUser)
	r.GET("/user/deleteUser", service.DeleteUser)
	r.POST("/user/updateUser", service.UpdateUser)
	r.POST("/user/findUserByNameAndPwd", service.FindUserByNameAndPwd)
	//发送
	r.GET("/user/sendMsg", service.SendMsg)
	r.GET("/user/sendUserMsg", service.SendUserMsg)
	r.POST("/user/redisMsg", service.RedisMsg)
	return r
}
