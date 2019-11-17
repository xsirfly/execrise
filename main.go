package main

import (
	"github.com/docker/docker/client"
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types"
	"os"
	"github.com/docker/docker/pkg/stdcopy"
	"time"
	"fmt"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	cli.NegotiateAPIVersion(ctx)

	resp, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image: "kylx:javademo",
		Cmd: []string{"../runner", "-key", "java", "-expected", "hello world"},
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	startTime := time.Now()
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	fmt.Println(time.Now().Unix() - startTime.Unix())
}
