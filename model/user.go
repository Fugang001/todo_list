package model

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	UserName       string `gorm:"unique"`
	PasswordDigest string `gorm:"index;not null"` //存储的是密文，加密后的密码
}

const (
	PassWordCost = 12 //密码加密难度
)

//SetPassword 设置密码 加密处理
func (user *User) SetPassword(password string) error {
	//										传进去的是字节类型，第二参数是加密难度
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PassWordCost)
	if err != nil {
		return err
	}
	user.PasswordDigest = string(bytes) //存储的是密文，加密后的密码转成字符串
	return nil
}

//CheckPassword 校验密码
func (user *User) CheckPassword(password string) bool {
	//							第一个参数是加密后的密码，第二个是原始密码
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
	return err == nil //是的话返回true
}
