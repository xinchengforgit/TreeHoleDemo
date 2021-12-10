package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func HashStr(str string) string {
	str = str + salt
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// func GenToken(uid int) (string, time.Time)
// func VerifyToken()

// 之后采用jwt鉴权认证
// 如何保证洞主还是洞主呢???
// 如果记录把
// 数据库里面存一个userId
type Claims struct {
	Uid int `json:"uid"`
	jwt.StandardClaims
}

func GenToken(uid int) (string, time.Time, error) {
	mySigningKey := []byte(JwtKey)
	expireTime := time.Now().Add(time.Duration(ExpireTime) * time.Minute)

	// Create the Claims
	claims := Claims{
		uid,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "xincheng",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", time.Now(), err
	}
	return ss, expireTime, nil
}
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("Invalid Token")
}
