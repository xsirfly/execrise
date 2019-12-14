package docker

import (
	"context"

	"github.com/docker/docker/client"
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
