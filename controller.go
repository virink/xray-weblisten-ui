package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
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
		obj WebVul
		err error
	)
	if err = c.Bind(&obj); err != nil {
		logger.Errorln(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, Resp{Code: 1, Msg: err.Error()})
		return
	}
	logger.Debugln(obj)
	if obj.Type == "web_vuln" {
		// TODO: Vul

		// 创建空的字段
		var vul = Vul{
			Model:      gorm.Model{},
			URL:        "", // 漏洞路径 *
			Domain:     "", // Host头 *
			Title:      "",
			Type:       "",
			Payload:    "", // 利用方法 *
			Params:     "",
			Plugin:     "", // 插件 *
			VulnClass:  "", // 漏洞类型 *
			CreateTime: 0,
			Raw:        "", // 特别大，可以考虑去掉
		}

		switch {
		case obj.Plugin == "dirscan":
			// 路径扫描格式
			raw, _ := json.Marshal(&obj.Detail)
			vul = Vul{
				CreateTime: obj.CreateTime,
				Domain:     obj.Detail.Host,
				URL:        obj.Target.URL,
				Title:      obj.Plugin,
				Type:       obj.Plugin,
				VulnClass:  obj.VulnClass,
				Plugin:     obj.Plugin,
				Params:     obj.Detail.Filename,
				Payload:    obj.Target.URL,
				Raw:        string(raw),
			}

		case strings.HasPrefix(obj.Plugin, "poc-"):
			// 自定义格式poc格式
			params, _ := json.Marshal(&obj.Detail.Param)
			raw, _ := json.Marshal(&obj.Detail)
			vul = Vul{
				CreateTime: obj.CreateTime,
				Domain:     obj.Detail.Host,
				URL:        obj.Target.URL,
				Title:      obj.Plugin,
				Type:       obj.Plugin,
				VulnClass:  obj.VulnClass,
				Plugin:     obj.Plugin,
				Params:     string(params),
				Payload:    obj.Detail.Request,
				Raw:        string(raw),
			}

		default:
			// 默认格式
			params, _ := json.Marshal(&obj.Detail.Param)
			raw, _ := json.Marshal(&obj.Detail)
			vul = Vul{
				CreateTime: obj.CreateTime,
				Domain:     obj.Detail.Host,
				URL:        obj.Detail.URL,
				Title:      obj.Detail.Title,
				Type:       obj.Detail.Type,
				VulnClass:  obj.VulnClass,
				Plugin:     obj.Plugin,
				Params:     string(params),
				Payload:    obj.Detail.Payload,
				Raw:        string(raw),
			}
		}

		if _, err = newVul(&vul); err != nil {
			logger.Errorln(err)
			c.AbortWithStatusJSON(http.StatusBadRequest, Resp{Code: 1, Msg: err.Error()})
			return
		}

	} else if obj.Type == "web_statistic" {
		// Statistic
		// num_found_urls - num_scanned_urls == 0 可以认为扫描结束了
		logger.Debugln("State")
		num_of_remain_scans := obj.NumFoundUrls - obj.NumScannedUrls
		if num_of_remain_scans == 0 {
			// 扫描完成
		}

	}
	c.JSON(http.StatusOK, Resp{Code: 0, Msg: "success"})
}

func portInUse(port int) bool {
	output, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("lsof -i:%d", port)).CombinedOutput()
	return err == nil && len(output) > 0
}

func randomPort() (p int) {
	for {
		p = 30000 + rand.Intn(10000)
		if !portInUse(p) {
			break
		}
	}
	return
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

func writeXrayConfg(obj Project) (err error) {
	xrayConf, err := getDefaultXrayConfig()
	if err != nil {
		logger.Errorln(err)
		return
	}
	domains := strings.Split(obj.Domain, ",")
	xrayConf.Mitm.Restriction.Includes = domains
	// TODO: BasicCrawler
	// xrayConf.BasicCrawler.Restriction.Includes = domains
	data, err := yaml.Marshal(&xrayConf)
	if err != nil {
		logger.Errorln(err)
		return
	}
	err = ioutil.WriteFile(obj.Config, data, 0666)
	if err != nil {
		logger.Errorln(err)
		return
	}
	return nil
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
