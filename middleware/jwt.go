package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	"todo_list/pkg/utils"
)

//JWT token验证中间件 鉴权
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		code = 200
		//一般从请求头里
		token := c.GetHeader("Authorization")
		if token == "" {
			code = 404
		} else {
			//token不为空就解析token
			claims, err := utils.ParseToken(token)
			if err != nil {
				//解析token出错
				code = 403 //无权限，说明token是假的
			} else if time.Now().Unix() > claims.ExpiresAt {
				//说明token过期
				code = 401 //token无效
			}
		}
		if code != 200 {
			c.JSON(400, gin.H{
				"status": code,
				"msg":    "Token解析错误",
			})
			//如果授权失败（例如：密码不匹配），则调用Abort以确保不调用此请求的其余处理程序
			c.Abort() //Abort 阻止等待的handlers被调用. 注意: 这不会停止当前handler
			return
		}
		//在调用的handler 中,执行链路上等待的handlers
		//调用Next函数会立刻自信后续的handler，反之，等本中间件执行完毕后，再执行后续的handler
		//立刻终止handler的执行和后续待执行的handler的执行的相关函数
		c.Next()
		fmt.Println("到这了")
	}
}
