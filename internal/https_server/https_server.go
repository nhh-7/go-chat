package httpsserver

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	v1 "github.com/nhh-7/go-chat/api/v1"
)

var GE *gin.Engine

func init() {
	GE = gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	GE.Use(cors.New(corsConfig))

	GE.POST("/login", v1.Login)
	GE.POST("/register", v1.Register)
	GE.POST("/user/updateUserInfo", v1.UpdateUserInfo)
	GE.POST("/user/getUserInfo", v1.GetUserInfo)

	GE.POST("/group/createGroup", v1.CreateGroup)
	GE.POST("/group/loadMyGroup", v1.LoadMyGroup)
	GE.POST("/group/checkGroupAddMode", v1.CheckGroupAddMode)
	GE.POST("/group/enterGroupDirectly", v1.EnterGroupDirectly)
	GE.POST("/group/leaveGroup", v1.LeaveGroup)
	GE.POST("/group/dismissGroup", v1.DismissGroup)
}
