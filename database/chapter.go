package database

type Chapter struct {
	ID           int64 `json:"id"`
	CourseID     int64 `json:"course_id"`
	Name         string `json:"name"`
	Tutorial     string `json:"tutorial"`
	CodeLocation string `json:"code_location"`
}

func GetChapter(id int64) (*Chapter, error) {
	var chapter Chapter
	err := db.Where("id = ?", id).First(&chapter).Error
	return &chapter, err
}

func GetChaptersByCourse(courseID int64) ([]*Chapter, error) {
	var chapters []*Chapter
	err := db.Where("course_id = ?", courseID).Find(&chapters).Error
	return chapters, err
}