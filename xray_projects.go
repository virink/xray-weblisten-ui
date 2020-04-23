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
	Listen    string // Listen Port
	ProcessID int
	// Worked    bool
}

func updateProjectPID(id uint, pid int) (out Project, err error) {
	if !conn.First(&out, Project{Model: gorm.Model{ID: id}}).RecordNotFound() {
		return out, errors.New("record is exists")
	}
	if err := conn.Model(&out).Updates(map[string]interface{}{"process_id": pid}).Error; err != nil {
		return out, err
	}
	return out, nil
}

func newProject(p Project) (out Project, err error) {
	if !conn.First(&out, Project{Name: p.Name}).RecordNotFound() {
		return out, errors.New("record is exists")
	}
	if err = conn.Create(&p).Error; err != nil {
		return p, err
	}
	return p, nil
}

func findProjects(limit, offset int) (outs []*Project, err error) {
	if conn.Find(&outs).Limit(limit).Offset(offset).RecordNotFound() {
		return outs, errors.New("record not found")
	}
	return outs, nil
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
