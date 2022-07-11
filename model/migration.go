package model

//执行数据迁移
func migration() {
	//自动迁移模式  DB是全局的，在model包ini.go里面所以不需要写Model.DB,本包下可以直接使用
	err := DB.Set("gorm:table_options", "charset=utf8").
		AutoMigrate(&User{}, &Task{}) //迁移User和Task
	if err != nil {
		return
	}
	//添加外键 uid关联到User表里面的id
	DB.Model(&Task{}).AddForeignKey("uid", "User(id)", "CASCADE", "CASCADE")
}
