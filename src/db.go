package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

type User struct {
	Id       int `gorm:"primaryKey"`
	Name     string
	Password string
}

// 首先得完成系统设计
// 考虑树洞的验证
// 每个树洞的匿名性怎么办???
// 这有一个设计原则在于同一个树洞下的帖子得是同一个Id, THU Hole是怎么维护的呢????
// 然后一个树洞需要维护的信息又哪些
// 考虑先设计几个接口

// 如何保证同一个人在同一个树洞里面id是一致的呢????
//
//
// 可以考虑利用hash pid + user ====> 生成一个值,
func InitDb() {
	dsn := fmt.Sprintf("%s?charset=utf8mb4&parseTime=True&loc=Local", viper.Get("mysql_resource"))
	fmt.Println(dsn)
	db1, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Can not open the mysql ")
	}
	db = db1
	db.AutoMigrate(&User{})
}

// 密码考虑加盐认证
func SaveUser(username, password string) error {
	// 存入表中
	passwordHash := HashStr(password)
	user := User{
		Name:     username,
		Password: passwordHash,
	}
	result := db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

//
func FindUserByUsername(username, password string) (*User, error) {
	var user User
	result := db.Where(&User{Name: username, Password: password}).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func FindUserByUid(uid int) (*User, error) {
	var user User
	result := db.Where(&User{Id: uid}).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
