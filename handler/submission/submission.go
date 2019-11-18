package submission

import (
	"bytes"
	"context"
	"encoding/json"
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
	"strings"
	"time"
)

type submitReq struct {
	Code string `json:"code"`
}
func Submit(c *gin.Context) {
	var req submitReq
	c.BindJSON(&req)
	submissionId := uuid.New()

	base := filepath.Join("/Users/zhangcong/oj", submissionId.String())
	path := filepath.Join(base, "Main.java")
	if err := os.MkdirAll(base, os.ModePerm); err != nil {
		abortRequesrtWithError(c, err)
		return
	}

	if err := ioutil.WriteFile(path, []byte(req.Code), 0666); err != nil {
		abortRequesrtWithError(c, err)
		return
	}
	dockerCtx := context.Background()
	resp, err := docker.GetCli().ContainerCreate(dockerCtx, &container.Config{
		Image: "xsirfly/kylx:javademo",
		Cmd: []string{"runner", "-key", "java", "-expected", "hello world"},
	}, &container.HostConfig{
		//AutoRemove: true,
		Mounts: []mount.Mount{
			{
				Type: mount.TypeBind,
				Source: base,
				Target: "/usr/src/app",
			},
		},
	}, nil, "")
	if err != nil {
		abortRequesrtWithError(c, err)
		return
	}

	startTime := time.Now()
	if err := docker.GetCli().ContainerStart(dockerCtx, resp.ID, types.ContainerStartOptions{}); err != nil {
		abortRequesrtWithError(c ,err)
		return
	}

	statusCh, errCh := docker.GetCli().ContainerWait(dockerCtx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			abortRequesrtWithError(c, err)
			return
		}
	case <-statusCh:
	}

	out, err := docker.GetCli().ContainerLogs(dockerCtx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		abortRequesrtWithError(c, err)
		return
	}
	var result bytes.Buffer
	if _, err := stdcopy.StdCopy(&result, &result, out); err != nil {
		abortRequesrtWithError(c, err)
		return
	}

	var r map[string]interface{}
	if err := json.Unmarshal([]byte(trim(string(result.String()))), &r); err != nil {
		abortRequesrtWithError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"submit_id": submissionId.String(),
		"result": r,
	})
	fmt.Println(time.Now().Unix() - startTime.Unix())
}

func abortRequesrtWithError(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"message": err.Error(),
	})
}

func trim(s string) string {
	i := strings.Index(s, "{")
	n := s[i:]
	return n
}
