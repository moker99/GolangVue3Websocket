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

	// swagger
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// 靜態資源(ex:顯示圖片)
	r.Static("/asset", "asset/")
	// r.LoadHTMLGlob("pages/**/*")

	// 首頁
	r.GET("/", service.GetIndex)
	r.GET("/index", service.GetIndex)
	r.GET("/toRegister", service.ToRegister)
	r.GET("/toChat", service.ToChat)
	r.GET("/chat", service.Chat)

	r.POST("/user/loadFriends", service.LoadFriends)

	// 用戶模塊
	r.POST("/user/getUserList", service.GetUserList)
	r.POST("/user/createUser", service.CreateUser)
	r.POST("/user/deleteUser", service.DeleteUser)
	r.POST("/user/updateUser", service.UpdateUser)
	r.POST("/user/findUserByNameAndPwd", service.FindUserByNameAndPwd)

	// 發送消息
	r.GET("/user/sendMsg", service.SendMsg)
	// 發送圖片
	r.POST("/user/upload", service.Upload)
	// 加好友
	r.POST("/user/addFriends", service.AddFriends)
	// 創建群
	r.POST("/user/createCommunity", service.CreateCommunity)

	r.POST("/user/loadCommunity", service.LoadCommunity)

	r.POST("/user/getChatHistory", service.GetChatHistory)

	return r
}
