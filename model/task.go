package model

import (
	"github.com/jinzhu/gorm"
)

//任务模型  是备忘录模型
//Uid不能为null,主要是有个外键关联，在migration.go
type Task struct {
	gorm.Model
	//哪个人的备忘录，外键关联Uid
	User      User   `gorm:"ForeignKey:Uid"`
	Uid       uint   `gorm:"not null"`
	Title     string `gorm:"index;not null"`
	Status    int    `gorm:"default:0"` //备忘录的状态，0是未完成，1是已完成
	Content   string `gorm:"type:longtext"`
	StartTime int64  //备忘录开始时间
	EndTime   int64  `gorm:"default:0"` //备忘录完成时间
}
