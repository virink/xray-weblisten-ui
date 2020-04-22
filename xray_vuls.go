package main

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// Statistics -
type Statistics struct {
	AverageResponseTime     float64 `json:"average_response_time"`
	NumFoundUrls            int64   `json:"num_found_urls"`
	NumScannedUrls          int64   `json:"num_scanned_urls"`
	NumSentHTTPRequests     int64   `json:"num_sent_http_requests"`
	RatioFailedHTTPRequests float64 `json:"ratio_failed_http_requests"`
}

// WebVul -
type WebVul struct {
	gorm.Model
	CreateTime int64 `json:"create_time"`
	Detail     struct {
		Host  string `json:"host"`
		Param struct {
			Key      string `json:"key"`
			Position string `json:"position"`
			Value    string `json:"value"`
		} `json:"param"`
		Payload       string `json:"payload"`
		Port          int64  `json:"port"`
		Request       string `json:"request"`
		Request1      string `json:"request1"`
		Request2      string `json:"request2"`
		Response      string `json:"response"`
		Response1     string `json:"response1"`
		Response2     string `json:"response2"`
		Title         string `json:"title"`
		Type          string `json:"type"`
		URL           string `json:"url"`
		ExpectedValue string `json:"expected_value"`
		HeaderName    string `json:"header_name"`
		HeaderValue   string `json:"header_value"`
	} `json:"detail"`
	Plugin string `json:"plugin"`
	Target struct {
		Params []struct {
			Path     []string `json:"path"`
			Position string   `json:"position"`
		} `json:"params"`
		URL string `json:"url"`
	} `json:"target"`
	Type      string `json:"type"`
	VulnClass string `json:"vuln_class"`
}

// Vul - 被动扫描项目
type Vul struct {
	gorm.Model
	Domain    string `json:"domain"` // xxx,xxx,xxx
	Title     string `json:"title"`
	Type      string `json:"type"`
	URL       string `json:"url"`
	Payload   string `json:"payload"`
	Plugin    string `json:"plugin"`
	VulnClass string `json:"vuln_class"`
	Raw       string `gorm:"type:text" json:"-"`
}

// vuln_name = models.TextField(blank=True)
// vuln_type = models.CharField(max_length=256, blank=True)
// vuln_url = models.CharField(max_length=256)
// vuln_payload = models.TextField(blank=True)

func newVul(p *Vul) (out *Vul, err error) {
	if !conn.First(&out, Vul{Domain: p.Domain}).RecordNotFound() {
		return out, errors.New("record is exists")
	}
	if err = conn.Create(p).Error; err != nil {
		return p, err
	}
	return p, nil
}

func findVulByID(id uint) (out Vul, err error) {
	if conn.First(&out, Vul{Model: gorm.Model{ID: id}}).RecordNotFound() {
		return out, errors.New("record not found")
	}
	return out, nil
}

func findVulByName(domain string) (out Vul, err error) {
	if conn.First(&out, Vul{Domain: domain}).RecordNotFound() {
		return out, errors.New("record not found")
	}
	return out, nil
}
