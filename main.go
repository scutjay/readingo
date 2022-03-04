package main

import (
	"github.com/alecthomas/log4go"
	"github.com/gin-gonic/gin"
	"net/http"
	"readingo/conf"
	"readingo/controller"
	"readingo/model"
	"readingo/service"
)

func main() {
	log4go.LoadConfiguration("log4go.xml")

	r := gin.New()
	openAuthApi(r)
	openRedisApi(r)
	openStaticFile(r)
	redirectForEmptyPath(r)
	r.Run(conf.Server.Host + ":" + conf.Server.Port)
}

func LocalAuth(ctx *gin.Context) {
	if token, err := ctx.Cookie("token"); err == nil && service.CheckIfValidToken(token) {
		ctx.Next()
	} else {
		ctx.JSON(http.StatusOK, model.NewNoPermissionResponse())
		ctx.Abort()
	}
}

func redirectForEmptyPath(r *gin.Engine) {
	r.GET("/", func(ctx *gin.Context) {
		if token, err := ctx.Cookie("token"); err == nil && service.CheckIfValidToken(token) {
			ctx.Redirect(http.StatusFound, "/view/index")
		} else {
			ctx.Redirect(http.StatusFound, "/view/login")
		}
	})
}

func openAuthApi(r *gin.Engine) {
	authApi := r.Group("/auth")
	{
		authApi.POST("/login", controller.ServiceHandler(service.Login, model.LoginReq{}))
		authApi.POST("/logout", controller.ServiceHandler(service.Logout, nil))
	}
}

func openRedisApi(r *gin.Engine) {
	redisApi := r.Group("/rds", LocalAuth)
	{
		redisApi.GET("/db-tree", controller.ServiceHandler(service.GetDBTree, nil))
		redisApi.GET("/refresh-db-tree", controller.ServiceHandler(service.RefreshDBTree, nil))
		redisApi.POST("/operate", controller.ServiceHandler(service.OperateRedis, model.OperateRedisReq{}))
		redisApi.GET("/supportedActions", controller.ServiceHandler(service.GetSupportedActions, nil))
	}
}

func openStaticFile(r *gin.Engine) {
	r.StaticFile("/view/index", "./view/index.html")
	r.StaticFile("/view/login", "./view/login.html")
	// Open static directory
	r.Static("/static", "./static")
}
