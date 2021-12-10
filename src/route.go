package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// 首先user表
// user_name password  uid
// Register

type LoginResponseData struct {
	Username   string    `json: "username"`
	Uid        int       `json: "uid"`
	Token      string    `json: "token"`
	ExpireTime time.Time `json: "expire_time"`
}

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

func ListenHttp() {
	r := gin.Default()
	r1 := r.Group("/api/v0")
	r1.POST("/register", Register)
	r1.POST("/login", Login)
	r1.GET("/userInfo", JWTAuth(), GetUserInfo)
	r.Run(":8080")
}
