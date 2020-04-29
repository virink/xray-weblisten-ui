package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

// FeishuData -
type FeishuData struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text"`
}

func pushFeishuMessage(title, text string) {
	var (
		req  *http.Request
		resp *http.Response
	)
	d := FeishuData{
		Text: text,
	}
	if len(title) > 0 {
		d.Title = title
	}
	data, err := json.Marshal(&d)
	if err != nil {
		logger.Errorln(err)
		return
	}
	if req, err = http.NewRequest("POST", conf.Server.Pusher, bytes.NewBuffer(data)); err != nil {
		logger.Errorln(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Transport: tr, Timeout: 5 * time.Second}
	if resp, err = client.Do(req); err != nil {
		logger.Errorln(err)
		return
	}
	defer resp.Body.Close()
}
