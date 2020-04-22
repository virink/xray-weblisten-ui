package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Resp - Web Response
type Resp struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func xrayWebhookHandler(c *gin.Context) {
	var (
		obj WebVul
		err error
	)
	err = c.ShouldBind(&obj)
	if err != nil {
		logger.Errorln(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, Resp{Code: 1, Msg: err.Error()})
	}
	// TODO: 处理 WebVul
	c.JSON(http.StatusOK, Resp{Code: 0, Msg: "success"})
}

func createProjectHandler(c *gin.Context) {}
func startProjectHandler(c *gin.Context)  {}
func stopProjectHandler(c *gin.Context)   {}
func getProjectsHandler(c *gin.Context)   {}
func getProjectHandler(c *gin.Context)    {}
