package main

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// Project - 被动扫描项目
type Project struct {
	gorm.Model
	Name      string `gorm:"type:varchar(100);unique_index" json:"name"`
	Domain    string // xxx,xxx,xxx
	Config    string // config_{name}.yaml
	Plugins   string
	Reverse   string
	Listen    string // 0.0.0.0:nnnnn
	ProcessID int
	// Worked    bool
}

func newProject(p *Project) (out *Project, err error) {
	if !conn.First(&out, Project{Name: p.Name}).RecordNotFound() {
		return out, errors.New("record is exists")
	}
	if err = conn.Create(p).Error; err != nil {
		return p, err
	}
	return p, nil
}

func findProjectByID(id uint) (out Project, err error) {
	if conn.First(&out, Project{Model: gorm.Model{ID: id}}).RecordNotFound() {
		return out, errors.New("record not found")
	}
	return out, nil
}

func findProjectByName(name string) (out Project, err error) {
	if conn.First(&out, Project{Name: name}).RecordNotFound() {
		return out, errors.New("record not found")
	}
	return out, nil
}
