package routes

import (
	"github.com/gin-contrib/sessions"
	_ "github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"todo_list/api"
	"todo_list/middleware"
)

//写路由，返回gin的引擎
func NewRouter() *gin.Engine {
	r := gin.Default() //生成了一个Web应用程序实例
	store := cookie.NewStore([]byte("something-very-secret"))
	r.Use(sessions.Sessions("mysession", store)) //对session进行存储
	//写基础路由
	v1 := r.Group("api/v1")
	{
		//用户操作，用户操作，调用的是api用户注册接口
		v1.POST("user/register", api.UserRegister)
		v1.POST("user/login", api.UserLogin)
		authed := v1.Group("/")      //需要登陆保护
		authed.Use(middleware.JWT()) //JWT鉴权，JWT是个中间件，下面分组用户只有登录才能操作
		{
			//任务操作  对备忘录的增删改查
			authed.GET("tasks", api.ListTasks)
			authed.POST("task", api.CreateTask)
			authed.GET("task/:id", api.ShowTask)
			authed.PUT("task/:id", api.UpdateTask)
			authed.POST("search", api.SearchTasks)
			authed.DELETE("task/:id", api.DeleteTask)
		}
	}
	return r
}
