package main

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
)

// WebVul - For Xray Webhook
type WebVul struct {
	gorm.Model
	CreateTime int64 `json:"create_time"`

	NumFoundUrls            int64   `json:"num_found_urls"`             // 发现的 url 数
	NumScannedUrls          int64   `json:"num_scanned_urls"`           // 扫描完成的 url 数
	NumSentHTTPRequests     int64   `json:"num_sent_http_requests"`     // 已发送的 http 请求数
	AverageResponseTime     float64 `json:"average_response_time"`      // 最近 30s 平均响应时间
	RatioFailedHTTPRequests float64 `json:"ratio_failed_http_requests"` // 最近 30s 请求失败率

	Detail struct {
		Host  string `json:"host"`
		Param struct {
			Key      string `json:"key"`
			Position string `json:"position"`
			Value    string `json:"value"`
		} `json:"param,omitempty"`
		Payload            string `json:"payload"`
		Port               int64  `json:"port"`
		Request            string `json:"request"`
		Response           string `json:"response"`
		Request1           string `json:"request1,omitempty"`
		Response1          string `json:"response1,omitempty"`
		Request2           string `json:"request2,omitempty"`
		Response2          string `json:"response2,omitempty"`
		Request3           string `json:"request3,omitempty"`
		Response3          string `json:"response3,omitempty"`
		Title              string `json:"title"`
		Type               string `json:"type"`
		URL                string `json:"url"`
		ExpectedValue      string `json:"expected_value,omitempty"`
		HeaderName         string `json:"header_name,omitempty"`
		HeaderValue        string `json:"header_value,omitempty"`
		ConfirmRetry       string `json:"confirm_retry,omitempty"`
		ConfirmRetryResult string `json:"confirm_retry_result,omitempty"`
		Filename           string `json:"filename,omitempty"`
	} `json:"detail"`
	Plugin string `json:"plugin"`
	Target struct {
		Params []struct {
			Path     []string `json:"path"`
			Position string   `json:"position"`
		} `json:"params"`
		URL string `json:"url"`
	} `json:"target,omitempty"`
	Type      string `json:"type"`
	VulnClass string `json:"vuln_class"`
}

// Vul - 被动扫描项目
type Vul struct {
	gorm.Model
	URL        string `gorm:"type:varchar(200);unique_index" json:"url"`
	Domain     string `json:"domain"` // xxx,xxx,xxx
	Title      string `json:"title"`
	Type       string `json:"type"`
	Payload    string `gorm:"type:text" json:"payload"`
	Params     string `json:"params"`
	Plugin     string `json:"plugin"`
	VulnClass  string `json:"vuln_class"`
	CreateTime int64  `json:"create_time"`
	Raw        string `gorm:"type:text" json:"-"`
}

func newVul(p Vul) (out Vul, err error) {
	if !conn.First(&out, Vul{URL: p.URL}).RecordNotFound() {
		return out, errors.New("record is exists")
	}
	if err = conn.Create(&p).Error; err != nil {
		return p, err
	}
	return p, nil
}

func findVuls(limit, offset int) (outs []*Vul, err error) {
	if conn.Find(&outs).Limit(limit).Offset(offset).RecordNotFound() {
		return outs, errors.New("record not found")
	}
	return outs, nil
}

func findVulByID(id uint) (out Vul, err error) {
	if conn.First(&out, Vul{Model: gorm.Model{ID: id}}).RecordNotFound() {
		return out, errors.New("record not found")
	}
	return out, nil
}

func findVulsByDomains(domain string, limit, offset int) (outs []*Vul, err error) {
	likes := strings.Split(strings.ReplaceAll(domain, "*", "%"), ",")
	stmp := conn.Limit(limit).Offset(offset)
	for _, like := range likes {
		stmp = stmp.Or("domain LIKE ?", like)
	}
	if stmp.Find(&outs).RecordNotFound() {
		return outs, errors.New("record not found")
	}
	return outs, nil
}
