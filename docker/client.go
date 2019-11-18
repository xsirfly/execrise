package docker

import (
	"github.com/docker/docker/client"
	"context"
)

var dockerCli *client.Client

func Init() error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	cli.NegotiateAPIVersion(ctx)
	dockerCli = cli
	return nil
}

func GetCli() *client.Client {
	return dockerCli
}


