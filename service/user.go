package service

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"readingo/conf"
	"readingo/constant"
	"readingo/model"
	"strconv"
	"strings"
	"time"
)

var userConfigs = make([]conf.UserConf, 0)

func init() {
	for _, ur := range conf.Users {
		if ur.Role == "" {
			ur.Role = "readonly"
		}
		userConfigs = append(userConfigs, conf.UserConf{
			Username: fmt.Sprintf("%X", md5.Sum([]byte(ur.Username))),
			Password: fmt.Sprintf("%X", md5.Sum([]byte(ur.Password))),
			Role:     ur.Role,
		})
	}
}

func checkUsernameAndPassword(usernameInMD5, passwordInMD5 string) string {
	for _, user := range userConfigs {
		if user.Username == usernameInMD5 && user.Password == passwordInMD5 {
			return user.Role
		}
	}
	return ""
}

func Login(ctx context.Context, req interface{}) (resp interface{}, code int, err error) {
	request := req.(*model.LoginReq)
	usernameInMD5 := strings.ToUpper(request.Username)
	passwordInMD5 := strings.ToUpper(request.Password)

	if role := checkUsernameAndPassword(usernameInMD5, passwordInMD5); role != "" {
		token := strconv.FormatInt(rand.New(rand.NewSource(time.Now().Unix())).Int63(), 10)
		addTokenToCache(token, role)
		setTokenToCookie(ctx, token)
		resp = model.LoginResp{Token: token}
	} else {
		return nil, model.NoPermission, errors.New("no permission")
	}
	return
}

func Logout(ctx context.Context, req interface{}) (resp interface{}, code int, err error) {
	ginCtx := ctx.(*gin.Context)
	token := GetTokenFromCookie(ginCtx)
	removeTokenFromCache(token)
	return
}

func GetTokenFromCookie(ctx context.Context) string {
	ginCtx := ctx.(*gin.Context)
	token, err := ginCtx.Cookie("token")
	if err != nil {
		return ""
	} else {
		return token
	}
}

func setTokenToCookie(ctx context.Context, token string) {
	ginCtx := ctx.(*gin.Context)
	ginCtx.SetCookie("token", token, constant.TokenDurationInCookie, "", "", false, false)
}
