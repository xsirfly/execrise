package chapter

import (
	"exercise/database"
	"exercise/utils"
	"gopkg.in/src-d/go-git.v4"
	"github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

type ChapterDto struct {
	*database.Chapter
	Files map[string]string `json:"files"`
}

func GetChapter(chapterId int64) (code int, resp interface{}) {
	var (
		chapter *database.Chapter
		course *database.Course
		err error
	)

	if chapter, err = database.GetChapter(chapterId); err != nil {
		return utils.ErrorResp(err)
	}

	if course, err = database.GetCourse(chapter.CourseID); err != nil {
		return utils.ErrorResp(err)
	}

	codeDir := utils.GetCodeDir(course)
	ok, err := utils.PathExists(codeDir)
	if err != nil {
		return utils.ErrorResp(err)
	}
	if !ok {
		_, err = git.PlainClone(codeDir, false, &git.CloneOptions{
			Auth: &http.BasicAuth{
				Username: "xsirfly",
				Password: "zc3245985wq",
			},
			URL: course.GitRepo,
		})
		if err != nil {
			logrus.WithError(err).Error("git clone failed")
			return utils.ErrorResp(err)
		}
	}

	sourceCode, err := utils.ReadAll(utils.GetChapterCodeDir(codeDir, chapter))
	if err != nil {
		return utils.ErrorResp(err)
	}

	return utils.SuccessResp(&ChapterDto{
		Chapter: chapter,
		Files: sourceCode,
	})






}