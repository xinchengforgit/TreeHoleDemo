package main

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// TODO, learn git log

var db *gorm.DB

type User struct {
	Id       int `gorm:"primaryKey"`
	Name     string
	Mail     string
	Password string
}

type Comment struct {
	Cid       int `gorm:"primaryKey`
	Pid       int `gorm:"index"`
	Name      string
	Text      string
	ReplyText string // 回复的帖子的内容,即楼中楼
	Uid       int    // name 通过 uid + cid 的hash
	UpdatedAt time.Time
	CreatedAt time.Time
}

type Post struct {
	Pid       int    `gorm:"primaryKey"`
	Text      string // Text
	Title     string // Title
	Type      string
	Uid       int
	UpdatedAt time.Time
	CreatedAt time.Time
	Replys    int
	Likes     int
}

func InitDb() {
	dsn := fmt.Sprintf("%s?charset=utf8mb4&parseTime=True&loc=Local", viper.Get("mysql_resource"))
	fmt.Println(dsn)
	db1, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Can not open the mysql ")
	}
	db = db1
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Post{})
}

func SaveUser(username, password string) error {

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

// 设计接口
// 首先post 一个帖子的时候
// 关键在于如何保证同一个用户在同一个帖子下面的name是一样的
func PostOne(text string, title string, ty string, uid int) error {
	post := Post{
		Text:      text,
		Type:      ty,
		Uid:       uid,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
		Likes:     0,
		Replys:    0,
	}
	result := db.Create(&post)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func PostComment(text string, replyText string, uid int, pid int) error {
	comment := Comment{
		Pid:       pid,
		Text:      text,
		ReplyText: replyText,
		Uid:       uid,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}
	result := db.Create(&comment)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetPostByPid(pid int) (*Post, error) {
	var post Post
	result := db.Where(&Post{Pid: pid}).First(&post)
	if result.Error != nil {
		return nil, result.Error
	}
	return &post, nil
}

func GetAllPosts() (*[]Post, error) {
	var posts []Post
	result := db.Find(&posts)
	if result.Error != nil {
		return nil, result.Error
	}
	return &posts, nil
}

func GetPostsByPidLists(pids []int) (*[]Post, error) {
	var posts []Post
	for _, pid := range pids {
		post, err := GetPostByPid(pid)
		if err != nil {
			return nil, err
		}
		posts = append(posts, *post)
	}
	return &posts, nil
}
