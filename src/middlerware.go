package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var response UniverseResponse
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			response = ErrorResponse(fmt.Sprint("Do not get token, forbidden"))
			c.JSON(403, response)
			c.Abort()
			return
		}
		tokenString = tokenString[7:]
		claims, err := ParseToken(tokenString)
		if err != nil {
			response = ErrorResponse(fmt.Sprintf("Invalid token"))
			c.JSON(403, response)
			c.Abort()
			return
		}
		c.Set("uid", claims.Uid)
		c.Next()
	}
}
