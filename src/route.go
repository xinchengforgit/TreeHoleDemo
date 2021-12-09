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
	}
	password, ok := c.GetPostForm("password")
	if !ok || password == "" {
		response = ErrorResponse(fmt.Sprintf("check your password"))
		c.JSON(400, response)
	}
	err := SaveUser(username, password)
	if err != nil {
		response = ErrorResponse(fmt.Sprintf("Register failed, please contact admin"))
		c.JSON(400, response)
	}
	response = SuccessResponse(nil)
	c.JSON(200, response)
}

func Login(c *gin.Context) {
	username, ok := c.GetPostForm("username")
	var response UniverseResponse
	if !ok || username == "" {
		response = ErrorResponse(fmt.Sprint("Bad Username"))
		c.JSON(400, response)
	}
	password, ok := c.GetPostForm("password")
	if !ok || password == "" {
		response = ErrorResponse(fmt.Sprintf("check your password"))
		c.JSON(400, response)
	}
	uid, hint, err := FindUser(username, password)
	if err != nil {
		response = ErrorResponse(fmt.Sprintln(hint))
		c.JSON(400, response)
	}
	token, exprire_at := GenToken(uid)
	loginResponseData := LoginResponseData{
		Username:   username,
		Uid:        uid,
		Token:      token,
		ExpireTime: exprire_at,
	}
	c.JSON(200, loginResponseData)
}
