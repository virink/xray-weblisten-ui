package main

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// Project - 被动扫描项目
type Project struct {
	gorm.Model
	Name      string `gorm:"type:varchar(100);unique_index" json:"name"`
	Domain    string `json:"domain"` // xxx,xxx,xxx
	Config    string `json:"config"` // config_{name}.yaml
	Plugins   string `json:"plugins"`
	Reverse   string `json:"reverse"` // ???
	Listen    int    `json:"listen"`  // Listen Port
	ProcessID int    `json:"process_id"`
	// Worked    bool
}

func updateProjectPidAndListenPort(id uint, pid, port int) (out Project, err error) {
	if conn.First(&out, Project{Model: gorm.Model{ID: id}}).RecordNotFound() {
		return out, errors.New("record is not found")
	}
	if err := conn.Model(&out).Updates(map[string]interface{}{
		"process_id": pid, "listen": port,
	}).Error; err != nil {
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
