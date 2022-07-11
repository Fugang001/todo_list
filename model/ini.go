package model

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	//下面这个包是mysql的驱动，不导进来连接不上啊sql: unknown driver "mysql" (forgotten import?)
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//定义一个gorm全局操作对象
var DB *gorm.DB

//下面进行数据库的连接
func Database(connstring string) {
	fmt.Println("connstring", connstring)
	db, err := gorm.Open("mysql", connstring)
	if err != nil {
		fmt.Println(err) //连接不上看看是什么错误
		panic("Mysql数据库连接失败")
	}
	fmt.Println("数据库连接成功")
	db.LogMode(true) //打印gorm日志
	//如果gin的mode是发行版本，就不让其输出
	if gin.Mode() == "release" {
		db.LogMode(false)
	}
	db.SingularTable(true)       //创建表的时候表名不加s user user
	db.DB().SetMaxIdleConns(20)  //设置连接池
	db.DB().SetMaxOpenConns(100) //设置连接数
	DB = db                      //把DB赋值到全局的变量上，这样数据库就连接好了
	//在connfig里面连接好 model.Database(path)
	migration() //数据库连接的时候对数据进行迁移
}
