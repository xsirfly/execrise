package submission

import (
	"context"
	"exercise/command"
	"exercise/database"
	"exercise/docker"
	"exercise/utils"
	"io/ioutil"
	"net/http"
	"os"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/docker/docker/pkg/stdcopy"
	"fmt"
	"time"
	"exercise/cache"
	"exercise/push"
)

type submitReq struct {
	CourseId  int64             `json:"course_id" uri:"course" binding:"required"`
	ChapterId int64             `json:"chapter_id" uri:"course" binding:"required"`
	Code      map[string]string `json:"code" form:"code" binding:"required"`
	ConnId    string            `json:"connId" binding:"required"`
}

const (
	SubmitStatusSuccess = "success"
	SubmitStatusFailed = "failed"
)

type RunLogOut struct {
	Conn *push.Connection
}

func (c *RunLogOut) Write(p []byte) (n int, err error) {
	buf := make([]byte, len(p))
	c.Conn.SendMessage(push.CmdRunLog, string(buf))
	return len(p), nil
}

func Submit(req *submitReq) (code int, response interface{}){
	submissionId := uuid.New()
	userID := int64(1)
	logrus.Infof("receview submit %v", req)

	var (
		course         *database.Course
		chapter        *database.Chapter
		commandContext command.Context
		err            error
		ok             bool
	)

	if course, err = database.GetCourse(req.CourseId); err != nil {
		return utils.ErrorResp(err)
	}

	if chapter, err = database.GetChapter(req.ChapterId); err != nil {
		return utils.ErrorResp(err)
	}

	if commandContext, ok = command.GetRunContextByBuild(course.BuildTool); !ok {
		return utils.ErrorResp(err)
	}

	codeDir := utils.GetCodeDir(course)
	submissionDir := utils.GetSubmissionDir(userID, course)

	exists, err := utils.PathExists(submissionDir)
	if err != nil {
		return utils.ErrorResp(err)
	}
	if !exists {
		if err := os.MkdirAll(submissionDir, os.ModePerm); err != nil {
			return utils.ErrorResp(err)
		}
		if err := utils.CopyDir(codeDir, submissionDir); err != nil {
			return utils.ErrorResp(err)
		}
	}

	for file, code := range req.Code {
		if err := ioutil.WriteFile(file, []byte(code), 0666); err != nil {
			return utils.ErrorResp(err)
		}
	}

	dockerCtx := context.Background()
	resp, err := docker.GetCli().ContainerCreate(dockerCtx, &container.Config{
		Image:      commandContext.GetImage(),
		WorkingDir: commandContext.GetWorkDir(),
		Cmd:        commandContext.GenCmd(chapter),
	}, &container.HostConfig{
		Mounts: commandContext.GenMounts(submissionDir),
	}, nil, "")
	if err != nil {
		return utils.ErrorResp(err)
	}

	if err := docker.GetCli().ContainerStart(dockerCtx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return utils.ErrorResp(err)
	}

	out, err := docker.GetCli().ContainerLogs(dockerCtx, resp.ID, types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Follow:     true,
		Tail:       "all",
		Details:    true,
	})
	if err != nil {
		return utils.ErrorResp(err)
	}

	go func() {
		wsconn, ok := push.GetPushClient(req.ConnId)
		if !ok {
			logrus.WithField("connId", req.ConnId).Error("ws connection not found")
			return
		}
		logOut := &RunLogOut{
			Conn: wsconn,
		}
		stdcopy.StdCopy(logOut, logOut, out)
	}()

	go func() {
		statusCh, errCh := docker.GetCli().ContainerWait(dockerCtx, resp.ID, container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				logrus.WithError(err).Error("container wait failed")
				cache.Client.RPush(genSubmitStatusToken(submissionId.String()), SubmitStatusFailed)
			}
		case <-statusCh:
			cache.Client.RPush(genSubmitStatusToken(submissionId.String()), SubmitStatusSuccess)
		}
	}()

	var r Result
	r.Success = true
	return http.StatusOK, gin.H{
		"submit_id": submissionId.String(),
	}
}

func GetSubmitStatus(token string) (code int, resp interface{}) {
	res, err := cache.Client.BLPop(30 * time.Second, genSubmitStatusToken(token)).Result()
	if err != nil {
		logrus.Error(err)
		return utils.ErrorResp(err)
	}
	return utils.SuccessResp(gin.H{
		"status": res[1],
	})
}

func genSubmitStatusToken(submitId string) string {
	return fmt.Sprintf("submit_status:%s", submitId)
}
