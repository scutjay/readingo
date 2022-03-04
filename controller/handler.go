package controller

import (
	"context"
	"github.com/alecthomas/log4go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"readingo/model"
	"reflect"
)

type ServiceFunc func(context.Context, interface{}) (interface{}, int, error)

func ServiceHandler(serviceFunc ServiceFunc, reqVal interface{}) func(*gin.Context) {
	var reqType reflect.Type = nil
	if reqVal != nil {
		value := reflect.Indirect(reflect.ValueOf(reqVal))
		reqType = value.Type()
	}

	return func(c *gin.Context) {
		tag := c.Request.RequestURI
		log4go.Info(tag + " enter")

		var req interface{} = nil
		if reqType != nil {
			req = reflect.New(reqType).Interface()

			//使用http 200 ok 响应code
			if err := c.ShouldBindBodyWith(req, binding.JSON); err != nil {
				log4go.Error("Bind json failed. error: %v", err)
				c.JSON(http.StatusOK, model.NewBindFailedResponse(tag))
				return
			}
		}

		resp, code, err := serviceFunc(c, req)
		if err != nil {
			log4go.Error(tag+" error: %v", err)
			c.JSON(http.StatusOK, model.NewErrorResponse(tag, code, err))
			return
		}

		c.JSON(http.StatusOK, model.NewSuccessResponse(tag, resp))
	}
}
