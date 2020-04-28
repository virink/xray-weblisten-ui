package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	signalChan chan os.Signal
	httpServer *http.Server
)

func newRouter() *gin.Engine {
	if conf.Server.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
		gin.DisableConsoleColor()
	}
	engine := gin.Default()

	engine.POST(webhook, xrayWebhookHandler)

	engine.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, DELETE")
		c.Header("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Accept")
		c.Header("Access-Control-Max-Age", "1800")
		if strings.ToUpper(c.Request.Method) == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})
	// TODO: AddRouter
	engine.POST("/project", createProjectHandler)
	// engine.PUT("/project/:id", updateProjectHandler)
	engine.GET("/start/:id", startProjectHandler)
	engine.GET("/stop/:id", stopProjectHandler)
	engine.GET("/projects", getProjectsHandler)
	engine.GET("/project/:id", getProjectHandler)
	engine.GET("/vuls", getVulsHandler)

	return engine
}

func main() {

	if _, err := getDefaultXrayConfig(); err != nil {
		logger.Errorln(err)
		return
	}

	router := newRouter()
	httpServer := &http.Server{
		Addr:           fmt.Sprintf(":%d", conf.Server.Port),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("HTTP server listen: %s\n", err)
		}
	}()

	logger.Debugf("HTTP Server Listening: http://127.0.0.1:%d", conf.Server.Port)

	// 等待中断信号以优雅地关闭服务器
	signalChan = make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt)
	sig := <-signalChan
	logger.Println("Get Signal:", sig)

	logger.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	logger.Println("Server exiting")
}
