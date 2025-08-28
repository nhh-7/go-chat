package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nhh-7/go-chat/internal/dto/request"
	"github.com/nhh-7/go-chat/internal/service/gorm"
	"github.com/nhh-7/go-chat/pkg/constants"
	"github.com/nhh-7/go-chat/utils/zlog"
)

func CreateGroup(c *gin.Context) {
	var createGroupReq request.CreateGroupRequest
	if err := c.BindJSON(&createGroupReq); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, ret := gorm.GroupInfoService.CreateGroup(createGroupReq)
	JsonBack(c, message, ret, nil)
}

func LoadMyGroup(c *gin.Context) {
	var loadMyGroupReq request.OwnListRequest
	if err := c.BindJSON(&loadMyGroupReq); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, groupList, ret := gorm.GroupInfoService.LoadMyGroup(loadMyGroupReq.OwnerId)
	JsonBack(c, message, ret, groupList)
}

func CheckGroupAddMode(c *gin.Context) {
	var req request.CheckGroupAddModeRequest
	if err := c.BindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, addMode, ret := gorm.GroupInfoService.CheckGroupAddMode(req.GroupId)
	JsonBack(c, message, ret, addMode)
}

// 传入群uuid和用户uuid，直接进入群聊
func EnterGroupDirectly(c *gin.Context) {
	var req request.EnterGroupDirectlyRequest
	if err := c.BindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, ret := gorm.GroupInfoService.EnterGroupDirectly(req.OwnerId, req.ContactId)
	JsonBack(c, message, ret, nil)
}

func LeaveGroup(c *gin.Context) {
	var req request.LeaveGroupRequest
	if err := c.BindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, ret := gorm.GroupInfoService.LeaveGroup(req.UserId, req.GroupId)
	JsonBack(c, message, ret, nil)
}

func DismissGroup(c *gin.Context) {
	var req request.DismissGroupRequest
	if err := c.BindJSON(&req); err != nil {
		zlog.Error(err.Error())
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": constants.SYSTEM_ERROR,
		})
		return
	}
	message, ret := gorm.GroupInfoService.DismissGroup(req.OwnerId, req.GroupId)
	JsonBack(c, message, ret, nil)
}
