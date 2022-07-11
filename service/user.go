package service

import (
	"github.com/jinzhu/gorm"
	"todo_list/model"
	"todo_list/pkg/utils"
	"todo_list/serializer"
)

// UserService 用户注册服务
type UserService struct {
	//UserName最小三位，最长15位Password最短5位，最长16位
	UserName string `form:"user_name" json:"user_name" binding:"required,min=3,max=15" example:"FanOne"`
	Password string `form:"password" json:"password" binding:"required,min=5,max=16" example:"FanOne666"`
}

//实现方法，返回json格式
func (service *UserService) Register() serializer.Response {
	var user model.User
	var count int64
	//                                                 在数据库里面查，验证看是否存在，存在就会返回错误
	model.DB.Model(&model.User{}).Where("user_name=?", service.UserName).First(&user).Count(&count)
	//表单验证  count=1说明表单里面已经有这个人了
	if count == 1 {
		return serializer.Response{
			Status: 400,
			Msg:    "已经有这个人了，无需再注册",
		}
	}
	//说明没这个人，下面就进行注册
	user.UserName = service.UserName
	//对密码进行加密
	if err := user.SetPassword(service.Password); err != nil {
		return serializer.Response{ //不为空说明加密失败
			Status: 400,
			Msg:    err.Error(),
		}
	}
	//走到这说明成功了，创建用户
	if err := model.DB.Create(&user).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "数据库操作错误",
		}
	}
	//执行到这说明用户创建成功
	return serializer.Response{
		Status: 200,
		Msg:    "用户创建成功",
	}
}

//Login 用户登陆函数
func (service *UserService) Login() serializer.Response {
	var user model.User
	//先去找一下这个user，看看数据库有没有这个人
	if err := model.DB.Where("user_name=?", service.UserName).First(&user).Error; err != nil {
		//如果查询不到，返回相应的错误
		if gorm.IsRecordNotFoundError(err) {
			return serializer.Response{
				Status: 400,
				Msg:    "用户不存在，请先登录",
			}
		}
		//如果不是用户不存在，是其它不可抗拒的因素导致的错误
		return serializer.Response{
			Status: 500,
			Msg:    "数据库错误",
		}
	}
	//如果有这个用户，那就对密码进行验证
	if user.CheckPassword(service.Password) == false {
		return serializer.Response{
			Status: 400,
			Msg:    "密码错误",
		}
	}
	//发一个token,为了其他功能需要身份验证所给前端存储的。
	//创建一个备忘录，这个功能就要token，不然都不知道是谁创建的备忘录。
	token, err := utils.GenerateToken(user.ID, service.UserName, service.Password)
	if err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "Token签发错误",
		}
	}
	return serializer.Response{
		Status: 200,
		Data:   serializer.TokenData{User: serializer.BuildUser(user), Token: token},
		Msg:    "登录成功",
	}
}
