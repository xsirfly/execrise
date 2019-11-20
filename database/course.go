package database

import "time"

type Course struct {
	ID          int64 `json:"id"`
	Name        string `json:"name"`
	Logo        string `json:"logo"`
	Description string `json:"description"`
	CreateAt    *time.Time `json:"create_at"`
	UpdateAt    *time.Time `json:"update_at"`
	CreateUser  int64 `json:"create_user"`
}

func GetCourses() ([]*Course, error) {
	var courses []*Course
	err := db.Table("course").Order("create_at desc").Where("is_del = 0").Find(&courses).Error
	if err != nil {
		return nil, err
	}
	return courses, nil
}
