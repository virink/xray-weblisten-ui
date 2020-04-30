package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

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
	if err = c.Bind(&obj); err != nil {
		logger.Errorln(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, Resp{Code: 1, Msg: err.Error()})
		return
	}
	if obj.Type == "web_vuln" {
		raw, _ := json.Marshal(&obj.Detail)
		vul := Vul{
			CreateTime: obj.CreateTime,
			Domain:     obj.Detail.Host,
			URL:        obj.Detail.URL,
			Raw:        string(raw),
			VulnClass:  obj.VulnClass,
			Plugin:     obj.Plugin,
		}
		switch {
		case obj.Plugin == "dirscan":
			// 路径扫描格式
			vul.Title = obj.Plugin
			vul.Type = obj.Plugin
			vul.Params = obj.Detail.Filename
			vul.Payload = obj.Target.URL
		case obj.Plugin == "sqldet":
			// SQL注入
			params, _ := json.Marshal(&obj.Detail.Param)
			vul.Title = obj.Detail.Title
			vul.Type = obj.Detail.Type
			vul.Params = string(params)
			vul.Payload = obj.Detail.Payload
		case strings.HasPrefix(obj.Plugin, "poc-"):
			// 自定义格式poc格式
			params, _ := json.Marshal(&obj.Detail.Param)
			vul.Title = obj.Plugin
			vul.Type = obj.Plugin
			vul.Params = string(params)
			vul.Payload = obj.Detail.Request
		default:
			// 默认格式
			params, _ := json.Marshal(&obj.Detail.Param)
			vul.Title = obj.Plugin
			vul.Type = obj.Detail.Type
			vul.Params = string(params)
			vul.Payload = obj.Detail.Payload
		}
		vul.Hash = MD5(vul.URL + vul.Plugin)
		if _, err = newVul(vul); err != nil {
			logger.Errorln(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, Resp{Code: 1, Msg: err.Error()})
			return
		}
		go pushFeishuMessage("发现漏洞", fmt.Sprintf(`Title: %s
Plugin: %s
URL:   %s`, vul.Title, vul.Plugin, vul.URL))
	} else if obj.Type == "web_statistic" {
		// Statistic
		project:= c.Param("project")
		if project == ""{
			return
		}
		// num_found_urls - num_scanned_urls == 0 可以认为扫描结束了
		if obj.NumFoundUrls == obj.NumScannedUrls {
			// 扫描完成
			logger.Debugln("Finish")
			logger.Debug(project)
			if _, flag := statistic[project]; !flag {
				// 创建pid时间戳
				statistic[project] = time.Now().Unix()
			}else{
				if time.Now().Unix() - statistic[project] > 3600 {

					var (
						obj Project
						err error
					)

					// 空闲超过1小时，主动结束
					defer delete(statistic, project)

					if obj ,err = findProjectByName(project); err !=nil{
						logger.Errorln(err.Error())
						c.AbortWithStatusJSON(http.StatusInternalServerError, Resp{Code: 1, Msg: err.Error()})
						return
					}
					//缺少根据ProcessID超着项目ID，能够停止，但无法更新项目状态！
					if obj.ProcessID == 0 || !processExists(obj.ProcessID) {
						c.JSON(http.StatusOK, Resp{Code: 0, Msg: "It's not running"})
						return
					}
					if err = syscall.Kill(obj.ProcessID, 15); err != nil {
						logger.Errorln(err)
						c.AbortWithStatusJSON(http.StatusInternalServerError, Resp{Code: 1, Msg: err.Error()})
						return
					}
					_, err = updateProjectPidAndListenPort(uint(obj.ID), 0, obj.Listen)
					if err != nil {
						logger.Errorln(err.Error())
						c.AbortWithStatusJSON(http.StatusInternalServerError, Resp{Code: 1, Msg: err.Error()})
						return
					}

				}
			}

		}else {
			if _,flag := statistic[project]; flag{
				//更新当前的时间戳
				statistic[project] = time.Now().Unix()
			}
		}
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
		objs  []*Project
		count int
		err   error
	)
	limit, offset := pagination(c.Query("page"), c.Query("page_size"))
	if objs, err = findProjects(limit, offset); err != nil {
		logger.Errorln(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, Resp{Code: 1, Msg: err.Error()})
		return
	}
	if count = getProjectsCount(); count == -1 {
		msg := "Cannot get count"
		logger.Errorln(msg)
		c.AbortWithStatusJSON(http.StatusInternalServerError, Resp{Code: 1, Msg: msg})
		return
	}
	// isRunning
	// for _,obj := range objs {
	// 	// obj.Listen
	// 	// obj.ProcessID
	// }
	c.JSON(http.StatusOK, Resp{Code: 0, Msg: "success", Data: map[string]interface{}{
		"count": count, "objs": objs,
	}})
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

	// 删除statics的pid值
	delete(statistic, obj.Name)
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
