package database

import "time"

type Course struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Logo        string     `json:"logo"`
	Description string     `json:"description"`
	Language    string     `json:"lanuage"`
	BuildTool   string     `json:"build_tool"`
	ProjectDir  string     `json:"project_dir"`
	Type        string     `json:"type"`
	CreateAt    *time.Time `json:"create_at"`
	UpdateAt    *time.Time `json:"update_at"`
	DeleteAt    *time.Time `json:"delete_at"`
	CreateUser  int64      `json:"create_user"`
	GitRepo     string     `json:"git_repo"`
}

func GetCourse(id int64) (*Course, error) {
	var course Course
	err := db.Where("id = ?", id).First(&course).Error
	return &course, err
}

func GetCourses() ([]*Course, error) {
	var courses []*Course
	err := db.Find(&courses).Error
	return courses, err
}
