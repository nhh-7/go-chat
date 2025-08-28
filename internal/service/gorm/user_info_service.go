package gorm

import (
	"errors"
	"fmt"
	"time"

	"github.com/nhh-7/go-chat/internal/dao"
	"github.com/nhh-7/go-chat/internal/dto/request"
	"github.com/nhh-7/go-chat/internal/dto/respond"
	"github.com/nhh-7/go-chat/internal/model"
	"github.com/nhh-7/go-chat/pkg/constants"
	"github.com/nhh-7/go-chat/pkg/enum/user_info/user_status_enum"
	"github.com/nhh-7/go-chat/utils/random"
	"github.com/nhh-7/go-chat/utils/zlog"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type userInfoService struct {
}

var UserInfoService = new(userInfoService)

func (u *userInfoService) Login(loginReq request.LoginRequest) (string, *respond.LoginRespond, int) {
	password := loginReq.Password
	var user model.UserInfo
	res := dao.GormDB.First(&user, "telephone = ?", loginReq.Telephone)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			message := "用户不存在"
			zlog.Error(message)
			return message, nil, -2
		}
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}
	if user.Password != password {
		message := "密码错误，请重新输入"
		zlog.Error(message)
		return message, nil, -2
	}

	loginRsp := &respond.LoginRespond{
		Uuid:      user.Uuid,
		Telephone: user.Telephone,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Avatar:    user.Avatar,
		Gender:    user.Gender,
		Birthday:  user.Birthday,
		Signature: user.Signature,
		IsAdmin:   user.IsAdmin,
		Status:    user.Status,
	}
	year, month, day := user.CreatedAt.Date()
	loginRsp.CreatedAt = fmt.Sprintf("%d.%d.%d", year, month, day)
	return "登录成功", loginRsp, 0
}

func (u *userInfoService) Register(registerReq request.RegisterRequest) (string, *respond.RegisterRespond, int) {
	if registerReq.SmsCode != "123456" {
		message := "验证码错误，请重新输入"
		zlog.Info(message)
		return message, nil, -2
	}

	message, ret := u.checkTelephoneExist(registerReq.Telephone)
	if ret != 0 {
		return message, nil, ret
	}
	var newUser model.UserInfo
	newUser.Uuid = "U" + random.GetNowAndLenRandomString(11)
	newUser.Telephone = registerReq.Telephone
	newUser.Password = registerReq.Password
	newUser.Nickname = registerReq.Nickname
	newUser.Avatar = "https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png"
	newUser.CreatedAt = time.Now()
	newUser.IsAdmin = u.checkUserIsAdminOrNot(newUser)
	newUser.Status = user_status_enum.NORMAL

	res := dao.GormDB.Create(&newUser)
	if res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}
	registerRsp := &respond.RegisterRespond{
		Uuid:      newUser.Uuid,
		Telephone: newUser.Telephone,
		Nickname:  newUser.Nickname,
		Email:     newUser.Email,
		Avatar:    newUser.Avatar,
		Gender:    newUser.Gender,
		Birthday:  newUser.Birthday,
		Signature: newUser.Signature,
		IsAdmin:   newUser.IsAdmin,
		Status:    newUser.Status,
	}
	year, month, day := newUser.CreatedAt.Date()
	registerRsp.CreatedAt = fmt.Sprintf("%d.%d.%d", year, month, day)
	return "注册成功", registerRsp, 0
}

func (u *userInfoService) UpdateUserInfo(updateUserInfoReq request.UpdateUserInfoRequest) (string, int) {
	var user model.UserInfo
	if res := dao.GormDB.First(&user, "uuid = ?", updateUserInfoReq.Uuid); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	if updateUserInfoReq.Email != "" {
		user.Email = updateUserInfoReq.Email
	}
	if updateUserInfoReq.Nickname != "" {
		user.Nickname = updateUserInfoReq.Nickname
	}
	if updateUserInfoReq.Birthday != "" {
		user.Birthday = updateUserInfoReq.Birthday
	}
	if updateUserInfoReq.Signature != "" {
		user.Signature = updateUserInfoReq.Signature
	}
	if updateUserInfoReq.Avatar != "" {
		user.Avatar = updateUserInfoReq.Avatar
	}
	if res := dao.GormDB.Save(&user); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	return "用户信息更新成功", 0
}

func (u *userInfoService) GetUserInfo(uuid string) (string, *respond.GetUserInfoRespond, int) {
	var user model.UserInfo
	if res := dao.GormDB.Where("uuid = ?", uuid).Find(&user); res.Error != nil {
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}
	rsp := respond.GetUserInfoRespond{
		Uuid:      user.Uuid,
		Telephone: user.Telephone,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Birthday:  user.Birthday,
		Email:     user.Email,
		Gender:    user.Gender,
		Signature: user.Signature,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		IsAdmin:   user.IsAdmin,
		Status:    user.Status,
	}
	return "获取用户信息成功", &rsp, 0
}

func (u *userInfoService) checkTelephoneExist(telephone string) (string, int) {
	var user model.UserInfo
	if res := dao.GormDB.Where("telephone = ?", telephone).First(&user); res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			zlog.Info("手机号不存在，可以注册", zap.String("tel:", telephone))
			return "", 0
		}
		zlog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, -1
	}
	zlog.Info("手机号已存在，无法注册")
	return "手机号已存在，无法注册", -2
}

func (u *userInfoService) checkUserIsAdminOrNot(user model.UserInfo) int8 {
	return user.IsAdmin
}
