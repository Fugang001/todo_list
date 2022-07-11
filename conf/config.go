package conf

import (
	"fmt"
	"gopkg.in/ini.v1"
	"strings"
	"todo_list/model"
)

var (
	AppMode     string
	HttpPort    string
	RedisDb     string
	RedisAddr   string
	RedisPw     string
	RedisDbName string
	Db          string
	DbHost      string
	DbPort      string
	DbUser      string
	DbPassWord  string
	DbName      string
)

func Init() {
	//是把文件加载出来了，具体配置下面才读取出来
	file, err := ini.Load("./conf/config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误，请检查文件路径")
	}
	LoadServer(file) //传service参数
	LoadMysql(file)  //传mysql参数
	//连接mysql库
	//拼接mysql连接字符串
	path := strings.Join([]string{DbUser, ":", DbPassWord, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8&parseTime=true"}, "")
	//path := strings.Join([]string{DbUser, ":", DbPassWord, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8mb4&parseTime=true"}, "")
	//上面只是对配置信息进行读写，下面是连接
	model.Database(path) //mysql是在model里面连接的
	//cache.Redis() //redis初始化，放到cache里面
}

//写个方法把配置信息加载出来
func LoadServer(file *ini.File) {
	//下面的方法是选中配置文件中的[service],找里面的配置
	AppMode = file.Section("service").Key("AppMode").String()
	HttpPort = file.Section("service").Key("HttpPort").String()
}
func LoadMysql(file *ini.File) {
	Db = file.Section("mysql").Key("Db").String()
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbPort = file.Section("mysql").Key("DbPort").String()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassWord = file.Section("mysql").Key("DbPassWord").String()
	DbName = file.Section("mysql").Key("DbName").String()
}
