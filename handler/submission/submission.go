package submission

import (
	"bytes"
	"context"
	"exercise/docker"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"exercise/utils"
	"exercise/command"
	"errors"
	"exercise/handler/submission/model"
	"github.com/sirupsen/logrus"
)

type submitReq struct {
	Code string `json:"code"`
}


func Submit(c *gin.Context) {
	var req submitReq
	c.BindJSON(&req)
	submissionId := uuid.New()

	logrus.Infof("receview submit %v", req)

	command.Init(":chapter01")
	runContext, _ := command.GetRunContextByProjectKey("gradleJava")
	var cmd []string
	cmd = append(cmd, runContext.RunCommand)
	for _, arg := range runContext.RunArgs {
		cmd = append(cmd, arg)
	}

	base := "/Users/xsir/ideaProjects/javademo/"
	path := filepath.Join(base, "chapter01/src/main/java/Add.java")
	if err := os.MkdirAll(base, os.ModePerm); err != nil {
		utils.AbortRequesrtWithError(c, err)
		return
	}

	if err := ioutil.WriteFile(path, []byte(req.Code), 0666); err != nil {
		utils.AbortRequesrtWithError(c, err)
		return
	}
	dockerCtx := context.Background()
	resp, err := docker.GetCli().ContainerCreate(dockerCtx, &container.Config{
		Image:      "kylx:gradle",
		WorkingDir: "/usr/src/app",
		Cmd:        cmd,
	}, &container.HostConfig{
		//AutoRemove: true,
		Mounts: []mount.Mount{
			{
				Type: mount.TypeBind,
				Source: base,
				Target: "/usr/src/app",
			},
			{
				Type: mount.TypeBind,
				Source: "/Users/xsir/.gradle",
				Target: "/home/gradle/.gradle",
			},
		},
	}, nil, "")
	if err != nil {
		utils.AbortRequesrtWithError(c, err)
		return
	}

	startTime := time.Now()
	if err := docker.GetCli().ContainerStart(dockerCtx, resp.ID, types.ContainerStartOptions{}); err != nil {
		utils.AbortRequesrtWithError(c ,err)
		return
	}

	statusCh, errCh := docker.GetCli().ContainerWait(dockerCtx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			utils.AbortRequesrtWithError(c, err)
			return
		}
	case <-statusCh:
	}

	out, err := docker.GetCli().ContainerLogs(dockerCtx, resp.ID, types.ContainerLogsOptions{ShowStderr: true})
	if err != nil {
		utils.AbortRequesrtWithError(c, err)
		return
	}
	var result bytes.Buffer
	if _, err := stdcopy.StdCopy(&result, &result, out); err != nil {
		utils.AbortRequesrtWithError(c, err)
		return
	}

	if len(result.String()) > 0 {
		utils.AbortRequesrtWithError(c, errors.New(result.String()))
		return
	}

	var r Result
	r.Success = true
	var testSuite model.TestSuite
	if err := testSuite.UnmarshalFromFile(filepath.Join(base, "chapter01/build/test-results/test/TEST-TestAdd.xml")); err != nil {
		utils.AbortRequesrtWithError(c, err)
		return
	}
	r.TestSuite = &testSuite
	c.JSON(http.StatusOK, gin.H{
		"submit_id": submissionId.String(),
		"result": r,
	})
	fmt.Println(time.Now().Unix() - startTime.Unix())
}
