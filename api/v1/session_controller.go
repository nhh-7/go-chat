package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nhh-7/go-chat/internal/dto/request"
	"github.com/nhh-7/go-chat/internal/service/gorm"
	"github.com/nhh-7/go-chat/pkg/constants"
	"github.com/nhh-7/go-chat/utils/zlog"
)

func OpenSession(c *gin.Context) {
	var openSessionReq request.OpensessionRequest
	if err := c.BindJSON(&openSessionReq); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	messsage, sessionId, ret := gorm.SessionService.OpenSession(openSessionReq)
	JsonBack(c, messsage, ret, sessionId)
}

func GetUserSessionList(c *gin.Context) {
	var getUserSessionListReq request.OwnListRequest
	if err := c.BindJSON(&getUserSessionListReq); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, sessionList, ret := gorm.SessionService.GetUserSessionList(getUserSessionListReq.OwnerId)
	JsonBack(c, message, ret, sessionList)
}

func GetGroupSessionList(c *gin.Context) {
	var req request.OwnListRequest
	if err := c.BindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, list, ret := gorm.SessionService.GetGroupSessionList(req.OwnerId)
	JsonBack(c, message, ret, list)
}

func DeleteSession(c *gin.Context) {
	var req request.DeleteSessionRequest
	if err := c.BindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, ret := gorm.SessionService.DeleteSession(req.OwnerId, req.SessionId)
	JsonBack(c, message, ret, nil)
}

func CheckOpenSessionAllowed(c *gin.Context) {
	var req request.CreateSessionRequest
	if err := c.BindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, res, ret := gorm.SessionService.CheckOpenSessionAllowed(req.SendId, req.ReceiveId)
	JsonBack(c, message, ret, res)
}
