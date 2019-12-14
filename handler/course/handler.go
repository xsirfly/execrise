package course

import (
	"exercise/database"
	"exercise/utils"
	"github.com/sirupsen/logrus"
)

func GetCourses() (code int, resp interface{}) {
	var (
		courses []*database.Course
		err error
	)
	if courses, err = database.GetCourses(); err != nil {
		logrus.WithError(err).Error("get course failed")
		return utils.ErrorResp(err)
	}

	return utils.SuccessResp(courses)
}

func GetChaptersByCourse(courseId int64) (code int, resp interface{}) {
	var (
		chapters []*database.Chapter
		err error
	)

	if chapters, err = database.GetChaptersByCourse(courseId); err != nil {
		return utils.ErrorResp(err)
	}

	return utils.SuccessResp(chapters)

}
