package utils

import (
	"github.com/sirupsen/logrus"
	"runtime/debug"
	"exercise/database"
	"path/filepath"
	"exercise/conf"
	"strconv"
)

func SafeGo(f func()) {
	go func() {
		defer func() {
			e := recover()
			logrus.Errorf("panic: %v, stack: %s", e, string(debug.Stack()))
		}()
		f()
	}()
}

func GetCodeDir(course *database.Course) string {
	return filepath.Join(conf.GetConf().CodeBaseDir, course.Language, course.ProjectDir)
}

func GetSubmissionDir(userId int64, course *database.Course) string {
	return filepath.Join(conf.GetConf().SubmissionBaseDir, strconv.Itoa(int(userId)), course.Language, course.ProjectDir)
}

func GetChapterCodeDir(codeDir string, chapter *database.Chapter) string {
	return filepath.Join(codeDir, chapter.CodeLocation)
}
