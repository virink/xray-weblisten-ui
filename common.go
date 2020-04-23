package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	configFileName = "config.yaml"
	loggerFilename = "weblisten.log"
	xrayConfigPath = "/tmp/"
	webhook        = "/vul_webhook"
)

// Config - Config
type Config struct {
	MySQL struct {
		Host    string `yaml:"host"`
		User    string `yaml:"user"`
		Pass    string `yaml:"pass"`
		Name    string `yaml:"name"`
		Charset string `yaml:"charset"`
	} `yaml:"mysql"`
	Xray struct {
		Path string `yaml:"path"`
		Bin  string `yaml:"bin"`
	} `yaml:"xray"`
	Server struct {
		Debug bool   `yaml:"debug"`
		Port  int    `yaml:"port"`
		Host  string `yaml:"host"`
	} `yaml:"server"`
}

var (
	tr      *http.Transport
	client  *http.Client
	conn    *gorm.DB
	logger  *logrus.Logger
	conf    Config
	xrayBin string
)

func templateConfig() []byte {
	conf := &Config{}
	conf.MySQL.Charset = "utf8mb4"
	conf.MySQL.Host = "127.0.0.1"
	conf.MySQL.User = "root"
	conf.MySQL.Pass = "123456"
	conf.MySQL.Name = "xray_weblisten"
	data, err := yaml.Marshal(conf)
	if err != nil {
		logger.Errorln(err.Error())
		return nil
	}
	return data
}

func loadConfig() (err error) {
	var yamlFile []byte
	_, err = os.Stat(configFileName)
	if err != nil && os.IsNotExist(err) {
		if data := templateConfig(); data != nil {
			if err = ioutil.WriteFile(configFileName, data, 0666); err != nil {
				return err
			}
		}
	}
	if yamlFile, err = ioutil.ReadFile(configFileName); err != nil {
		logger.Errorln(err.Error())
		return err
	}
	if err = yaml.Unmarshal(yamlFile, &conf); err != nil {
		logger.Errorln(err.Error())
		return err
	}
	return nil
}

func initLogger(filename string, level logrus.Level) *logrus.Logger {
	logger = logrus.New()
	logger.SetLevel(level)
	if level == logrus.DebugLevel || level == logrus.InfoLevel {
		logger.SetFormatter(&logrus.TextFormatter{
			ForceColors:            true,
			DisableLevelTruncation: true,
			TimestampFormat:        "2006-01-02 15:04:05",
		})
		logger.SetOutput(os.Stdout)
	} else {
		logger.SetFormatter(&logrus.JSONFormatter{})
		logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logger.SetOutput(logFile)
		}
	}
	return logger
}

func initConnect() (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local",
		conf.MySQL.User, conf.MySQL.Pass, conf.MySQL.Host, conf.MySQL.Name, conf.MySQL.Charset,
	)
	if db, err = gorm.Open("mysql", dsn); err != nil {
		logger.Errorln(err.Error())
		return nil, err
	}
	db.LogMode(conf.Server.Debug)
	// db.Debug()
	db.DB().SetConnMaxLifetime(100 * time.Second) // 最大连接周期，超过时间的连接就close
	db.DB().SetMaxOpenConns(100)                  // 设置最大连接数
	db.DB().SetMaxIdleConns(16)                   // 设置闲置连接数
	return
}

func init() {
	logger = initLogger(loggerFilename, logrus.DebugLevel)
	logger.AddHook(NewLogHook())
	err := loadConfig()
	if err != nil {
		logger.Fatalln(err.Error())
		return
	}
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{
		Transport: tr,
		Timeout:   5 * time.Second,
	}
	conn, _ = initConnect()

	// Model
	//conn.CreateTable(&Project{}, &Vul{})

	xrayBin = filepath.Join(conf.Xray.Path, conf.Xray.Bin)
	rand.Seed(time.Now().UnixNano())
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
