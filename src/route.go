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
	ReplyText string `json: "reply_text,omitempty"`
}

type RenderPost struct {
	Content string `json: "content"`
	Title   string `json: "title`
	Pid     int    `json: "pid"`
	Likes   int    `json: "likes"`
	Replys  int    `json: "replys"`
}

type RenderComment struct {
	Name      string `json: "name"` // 表示回复者的名字
	Content   string `json: "content"`
	ReplyText string `json: "reply_text,omitempty"` // 楼中楼
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
	err2 := PostOne(newpost.Content, newpost.Title, newpost.Title, uid.(int))
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
	pid2, err3 := strconv.Atoi(pid)
	if err3 != nil {
		response = ErrorResponse(fmt.Sprintf("invalid pid"))
		c.JSON(400, response)
		return
	}
	err2 := PostComment(comment.Content, comment.ReplyText, uid.(int), pid2)
	if err2 != nil {
		response = ErrorResponse(fmt.Sprint("store comment error, please contact the admin"))
		c.JSON(400, response)
		return
	}
	err4 := UpdatePost(pid2)
	if err4 != nil {
		response = ErrorResponse(fmt.Sprintf("update replys wrong, please contact the admin"))
		c.JSON(400, response)
		return
	}
	response = SuccessResponse(nil)
	c.JSON(200, response)
}

// Posts
func GetPosts(c *gin.Context) {
	// GetPosts;
	// 这个东西不需要鉴权;
	var datas []RenderPost
	var response UniverseResponse
	// 获取全部的posts
	posts, err := GetAllPosts()
	if err != nil {
		response = ErrorResponse(fmt.Sprint("get posts error, please contact the admin"))
		c.JSON(400, response)
	}
	for _, v := range *posts {
		data := RenderPost{
			Title:   v.Title,
			Content: v.Text,
			Pid:     v.Pid,
			Replys:  v.Replys,
			Likes:   v.Likes,
		}
		datas = append(datas, data)
	}
	response = SuccessResponse(datas)
	c.JSON(200, response)
}

func GetOnePostInfo(c *gin.Context) {
	pid := c.Param("pid")
	pid2, err := strconv.Atoi(pid)
	var response UniverseResponse
	if err != nil {
		response = ErrorResponse(fmt.Sprintf("invalid pid"))
		c.JSON(400, response)
	}
	var post RenderPost
	p, err2 := GetPostByPid(pid2)
	if err2 != nil {
		response = ErrorResponse(fmt.Sprintf("can not find the post"))
		c.JSON(400, response)
	}
	post = RenderPost{
		Pid:     pid2,
		Title:   p.Title,
		Content: p.Text,
		Likes:   p.Likes,
		Replys:  p.Replys,
	}
	response = SuccessResponse(post)
	c.JSON(200, response)
}

func GetOnePostComment(c *gin.Context) {
	//获取参数
	pid := c.Param("pid") // 获取pid
	pid2, err := strconv.Atoi(pid)
	var response UniverseResponse
	if err != nil {
		response = ErrorResponse(fmt.Sprintf("invalid pid"))
		c.JSON(400, response)
	}
	var comments []RenderComment
	datas, err2 := GetCommentsByPid(pid2)
	if err2 != nil {
		response = ErrorResponse(fmt.Sprintf("find data wrong, please contact the admin"))
		c.JSON(400, response)
	}
	for _, data := range *datas {
		// 根据pid 和 uid生成一个名字
		name := GenName(data.Pid, data.Uid)
		comment := RenderComment{
			Name:      name,
			Content:   data.Text,
			ReplyText: data.ReplyText,
		}
		comments = append(comments, comment)
	}
	response = SuccessResponse(comments)
	c.JSON(200, response)
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
	// 获取的是下面的所有的评论
	// 关键在于保证每一个人的
	r1.GET("/post/:pid", JWTAuth(), GetOnePostInfo)
	r1.GET("/post/:pid/replys", JWTAuth(), GetOnePostComment)
	r.Run(":8080")
}
