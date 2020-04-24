package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/util/json"
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
	if err = c.Bind(&obj); err != nil {
		logger.Errorln(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, Resp{Code: 1, Msg: err.Error()})
		return
	}
	logger.Debugln(obj)
	if obj.CreateTime > 0 {
		params, _ := json.Marshal(&obj.Detail.Param)
		vul := Vul{
			CreateTime: obj.CreateTime,
			Domain:     obj.Detail.Host,
			URL:        obj.Detail.URL,
			Title:      obj.Detail.Title,
			Type:       obj.Detail.Type,
			VulnClass:  obj.VulnClass,
			Plugin:     obj.Plugin,
			Params:     string(params),
		}
		if _, err = newVul(vul); err != nil {
			logger.Errorln(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, Resp{Code: 1, Msg: err.Error()})
			return
		}
		// TODO: Push Message
	} else {
		// Stat
		logger.Debugln(obj)
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
	obj.Name = strings.ReplaceAll(strings.ToLower(obj.Name), " ", "_")
	obj.Config = filepath.Join(xrayConfigPath, fmt.Sprintf("xray_config_%s.yaml", obj.Name))
	obj.Listen = -1
	out, err := newProject(obj)
	if err != nil {
		logger.Errorln(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, Resp{Code: 2, Msg: err.Error()})
		return
	}
	err = writeXrayConfg(out)
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
	if obj.Listen <= 0 || portInUse(obj.Listen) {
		obj.Listen = randomPort()
	}
	xrayArgs := []string{
		"--config", obj.Config,
		"webscan", "--listen", fmt.Sprintf("0.0.0.0:%d", obj.Listen),
		"--webhook-output", fmt.Sprintf("http://%s:%d%s", conf.Server.Host, conf.Server.Port, webhook),
		"--plugins", obj.Plugins,
	}
	logger.Debugln(xrayArgs)
	var pid int
	if pid, err = runXray(xrayArgs); err != nil {
		logger.Errorln(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, Resp{Code: 1, Msg: err.Error()})
		return
	}
	if obj, err = updateProjectPidAndListenPort(obj.ID, pid, obj.Listen); err != nil {
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
	out, err := updateProjectPidAndListenPort(uint(id), 0, obj.Listen)
	if err != nil {
		logger.Errorln(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, Resp{Code: 1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, Resp{Code: 0, Msg: "success", Data: out})
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

func getVulsHandler(c *gin.Context) {
	var (
		objs    []*Vul
		project Project
		err     error
		domain  string
	)
	// Find Vuls By Project ID -> Domain
	// Find Vuls By Domain
	limit, offset := pagination(c.Query("page"), c.Query("page_size"))
	domain = c.Query("domain")
	id, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if id > 0 {
		if project, err = findProjectByID(uint(id)); err != nil {
			logger.Errorln(err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, Resp{Code: 1, Msg: err.Error()})
			return
		}
		domain = project.Domain
	}
	domain = strings.TrimSpace(domain)
	if objs, err = findVulsByDomains(domain, limit, offset); err != nil {
		logger.Errorln(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, Resp{Code: 1, Msg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, Resp{Code: 0, Msg: "success", Data: objs})
}
