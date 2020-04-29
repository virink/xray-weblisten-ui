package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"gopkg.in/yaml.v2"
)

func processExists(pid int) bool {
	if err := syscall.Kill(pid, 0); err == nil {
		return true
	}
	return false
}

func runXray(args []string) (pid int, err error) {
	cmd := exec.Command(xrayBin, args...)
	cmd.Dir = conf.Xray.Path
	if err = cmd.Start(); err != nil {
		return 0, err
	}
	return cmd.Process.Pid, nil
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

// pagination - 分页生成
// @return limit,offset
func pagination(page, pageSize string) (int, int) {
	pageSizeInt, _ := strconv.Atoi(pageSize)
	pageInt, _ := strconv.Atoi(page)
	if pageSizeInt < 20 || pageSizeInt > 100 {
		pageSizeInt = 20
	}
	return pageSizeInt, (pageInt - 1) * pageSizeInt
}

// MD5 -
func MD5(text string) string {
	ctx := md5.New()
	_, _ = ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}
