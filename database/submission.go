package database

import "time"

type Submission struct {
	ID              int64      `json:"id"`
	UserID          int64      `json:"user_id"`
	ChapterID       int64      `json:"chapter_id"`
	SubmissionToken string     `json:"submission_toke"`
	Code            string     `json:"code"`
	CreateAt        *time.Time `json:"create_at"`
	Status          int64      `json:"status"`
}

func CreateSubmission(submission *Submission) error {
	return db.Create(submission).Error
}

func GetSubmissions(userID, chapterID int64) ([]*Submission, error) {
	var submissions []*Submission
	err := db.Where("user_id = ? && chapter_id = ?", userID, chapterID).Order("create_at desc").Find(&submissions).Error
	return submissions, err
}
