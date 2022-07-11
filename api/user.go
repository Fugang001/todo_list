package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"todo_list/service"
)

//新建用户操作的接口
// UserRegister @Tags USER
// @Summary 用户注册
// @Produce json
// @Accept json
// @Param data body service.UserService true "用户名, 密码"
// @Success 200 {object} serializer.ResponseUser "{"status":200,"data":{},"msg":"ok"}"
// @Failure 500  {object} serializer.ResponseUser "{"status":500,"data":{},"Msg":{},"Error":"error"}"
// @Router /user/register [post]
func UserRegister(c *gin.Context) { ////把上下文传递进来
	//写一个服务，服务在service里面创建
	var userRegisterService service.UserService                //相当于创建了一个UserRegisterService对象，调用这个对象中的Register方法。
	if err := c.ShouldBind(&userRegisterService); err == nil { //绑定对象就会把值传递过来，把c里面的值赋值给userRegisterService
		res := userRegisterService.Register() //执行注册方法，执行绑定
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
		fmt.Println("出错了")
		fmt.Println(err)
	}
}
func UserLogin(c *gin.Context) {
	//写一个服务，服务在service里面创建
	var userLoginService service.UserService                //相当于创建了一个UserRegisterService对象，调用这个对象中的Register方法。
	if err := c.ShouldBind(&userLoginService); err == nil { //绑定对象就会把值传递过来，把c里面的值赋值给userRegisterService
		res := userLoginService.Login() //执行登录方法
		c.JSON(200, res)
	} else {
		c.JSON(400, err)
	}
}
