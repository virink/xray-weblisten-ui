package main

import (
	"bytes"
	"net/http"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/gin-gonic/gin"
)

func processExists(pid int) bool {
	if err := syscall.Kill(pid, 0); err == nil {
		return true
	}
	return false
}

func runXray(args []string) (pid int, err error) {
	cmd := exec.Command(xrayBin, args...)
	if err = cmd.Start(); err != nil {
		return 0, err
	}
	return cmd.Process.Pid, nil
}

// Resp - Web Response
type Resp struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func xrayWebhookHandler(c *gin.Context) {
	var (
		// obj WebVul
		err error
	)
	data, err := c.GetRawData()
	if bytes.Contains(data, []byte("num_found_urls")) {
		// TODO: Statistics
	} else {
		var obj WebVul
		err = c.ShouldBind(&obj)
		// TODO: 处理 WebVul
	}
	if err != nil {
		logger.Errorln(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, Resp{Code: 1, Msg: err.Error()})
	}
	c.JSON(http.StatusOK, Resp{Code: 0, Msg: "success"})
}

func createProjectHandler(c *gin.Context) {
	var (
		obj Project
		err error
	)
	err = c.ShouldBind(&obj)
	if err != nil {
		logger.Errorln(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, Resp{Code: 1, Msg: err.Error()})
		return
	}
	out, err := newProject(obj)
	if err != nil {
		logger.Errorln(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, Resp{Code: 2, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, Resp{Code: 0, Msg: "success", Data: out})
}

func getProjectsHandler(c *gin.Context) {
	var (
		objs []*Project
		err  error
	)
	limit, offset := pagination(c.Query("page"), c.Query("page_size"))
	if objs, err = findProjects(limit, offset); err != nil {
		logger.Errorln(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, Resp{Code: 1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, Resp{Code: 0, Msg: "success", Data: objs})
}

func getProjectHandler(c *gin.Context) {
	var (
		obj Project
		err error
	)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if obj, err = findProjectByID(uint(id)); err != nil {
		logger.Errorln(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, Resp{Code: 1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, Resp{Code: 0, Msg: "success", Data: obj})
}

func startProjectHandler(c *gin.Context) {
	var (
		obj Project
		err error
	)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if obj, err = findProjectByID(uint(id)); err != nil {
		logger.Errorln(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, Resp{Code: 1, Msg: err.Error()})
		return
	}
	if obj.ProcessID > 0 && processExists(obj.ProcessID) {
		c.JSON(http.StatusOK, Resp{Code: 0, Msg: "Already Started!"})
		return
	}
	xrayArgs := []string{
		"--config", obj.Config,
		"webscan", "--listen", obj.Listen,
		"--webhook-output", "http://127.0.0.1:8088" + webhook,
		"--plugins", obj.Plugins,
	}
	logger.Debugln(xrayArgs)
	var pid int
	if pid, err = runXray(xrayArgs); err != nil {
		logger.Errorln(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, Resp{Code: 1, Msg: err.Error()})
		return
	}
	if obj, err = updateProjectPID(obj.ID, pid); err != nil {
		logger.Errorln(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, Resp{Code: 1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, Resp{Code: 0, Msg: "success", Data: obj})
}

func stopProjectHandler(c *gin.Context) {
	var (
		obj Project
		err error
	)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if obj, err = findProjectByID(uint(id)); err != nil {
		logger.Errorln(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, Resp{Code: 1, Msg: err.Error()})
		return
	}
	if obj.ProcessID == 0 || !processExists(obj.ProcessID) {
		c.JSON(http.StatusOK, Resp{Code: 0, Msg: "It's not running"})
		return
	}
	if err = syscall.Kill(obj.ProcessID, 15); err != nil {
		logger.Errorln(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, Resp{Code: 1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, Resp{Code: 0, Msg: "success", Data: obj})
}
