package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestProjectJSON(t *testing.T) {
	p := Project{
		Name:    "test_project",
		Domain:  "*.sandbox.ctfhub.com",
		Config:  "test_project",
		Plugins: "sqldet,cmd_injection",
		Listen:  19999,
	}
	data, err := json.Marshal(&p)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(data))
}
