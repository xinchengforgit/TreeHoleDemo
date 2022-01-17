package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginResponseData struct {
	Username   string    `json: "username"`
	Uid        int       `json: "uid"`
	Token      string    `json: "token"`
	ExpireTime time.Time `json: "expire_time"`
}

type DeliveryPost struct {
	Content string `json: "content"`
	Title   string `json: "title"`
}

type DeliveryComment struct {
	Content   string `json: "content"`
	ReplyText string `json: "reply_text"`
}

// 注册和登陆采用表单

func Register(c *gin.Context) {
	username, ok := c.GetPostForm("username")
	var response UniverseResponse
	if !ok || username == "" {
		response = ErrorResponse(fmt.Sprint("Bad Username"))
		c.JSON(400, response)
		return
	}
	password, ok := c.GetPostForm("password")
	if !ok || password == "" {
		response = ErrorResponse(fmt.Sprintf("check your password"))
		c.JSON(400, response)
		return
	}
	err := SaveUser(username, password)
	if err != nil {
		response = ErrorResponse(fmt.Sprintf("Register failed, please contact admin"))
		c.JSON(400, response)
		return
	}
	response = SuccessResponse(nil)
	c.JSON(200, response)
	return
}

func Login(c *gin.Context) {
	username, ok := c.GetPostForm("username")
	var response UniverseResponse
	if !ok || username == "" {
		response = ErrorResponse(fmt.Sprint("Bad Username"))
		c.JSON(400, response)
		return
	}
	password, ok := c.GetPostForm("password")
	if !ok || password == "" {
		response = ErrorResponse(fmt.Sprintf("check your password"))
		c.JSON(400, response)
		return
	}
	passwordHash := HashStr(password)
	user, err := FindUserByUsername(username, passwordHash)
	if err != nil {
		response = ErrorResponse(fmt.Sprintln(err))
		c.JSON(400, response)
		return
	}
	uid := user.Id
	token, expireTime, err2 := GenToken(uid)
	if err2 != nil {
		response = ErrorResponse(fmt.Sprintf("Generate token failed, please contact admin"))
		c.JSON(400, response)
		return
	}
	//token, exprire_at := GenToken(uid)
	loginResponseData := LoginResponseData{
		Username:   username,
		Uid:        uid,
		Token:      token,
		ExpireTime: expireTime,
	}
	response = SuccessResponse(loginResponseData)
	c.JSON(200, response)
	return
}

func GetUserInfo(c *gin.Context) {
	uid, ok := c.Get("uid")
	var response UniverseResponse
	if !ok {
		response = ErrorResponse(fmt.Sprintf("can not get the info by token"))
		c.JSON(400, response)
		return
	}
	user, err := FindUserByUid(uid.(int))
	if err != nil {
		response = ErrorResponse(fmt.Sprintf("can not find user"))
		c.JSON(400, response)
		return
	}
	username := user.Name
	response = SuccessResponse(gin.H{
		"username": username,
		"uid":      user.Id,
	})
	c.JSON(200, response)
	return
}

func PostNewPost(c *gin.Context) {
	uid, ok := c.Get("uid")
	var response UniverseResponse
	if !ok {
		response = ErrorResponse(fmt.Sprintf("can not get the info by token"))
		c.JSON(400, response)
		return
	}
	var newpost DeliveryPost
	err := c.BindJSON(&newpost)
	if err != nil {
		response = ErrorResponse(fmt.Sprintf("invalid json format"))
		c.JSON(400, response)
		return
	}
	err2 := PostOne(newpost.Content, "text", newpost.Title, uid.(int))
	if err2 != nil {
		response = ErrorResponse(fmt.Sprintf("store post error, please contact the admin"))
		c.JSON(400, response)
		return
	}
	response = SuccessResponse(nil)
	c.JSON(200, response)
}

// 首先是获取参数
func PostNewComment(c *gin.Context) {
	pid := c.Param("pid") // 获取pid
	uid, ok := c.Get("uid")
	var response UniverseResponse
	if !ok {
		response = ErrorResponse(fmt.Sprintf("can not get the info by token"))
		c.JSON(400, response)
		return
	}
	var comment DeliveryComment
	err := c.BindJSON(&comment)
	if err != nil {
		response = ErrorResponse(fmt.Sprint("invalid json format"))
		c.JSON(400, response)
		return
	}
	pid2, _ := strconv.Atoi(pid)
	err2 := PostComment(comment.Content, comment.ReplyText, uid.(int), pid2)
	if err2 != nil {
		response = ErrorResponse(fmt.Sprint("store comment error, please contact the admin"))
	}
	response = SuccessResponse(nil)
	c.JSON(200, response)
}

func GetPosts(c *gin.Context) {
	// 就是这个地方, 才是展示给name的真正地方

}

func ListenHttp() {
	r := gin.Default()
	r1 := r.Group("/api/v0")
	r1.POST("/register", Register)
	r1.POST("/login", Login)
	r1.GET("/userInfo", JWTAuth(), GetUserInfo)
	r1.POST("/post", JWTAuth(), PostNewPost)
	r1.POST("/post/:pid", JWTAuth(), PostNewComment)
	r1.GET("/post", JWTAuth(), GetPosts)
	r.Run(":8080")
}
