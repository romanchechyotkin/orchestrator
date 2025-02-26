package docker

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type Client struct {
	cli *client.Client
}

func NewClient(opts ...client.Opt) (*Client, error) {
	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		cli: cli,
	}, nil
}

func (cl *Client) Run(ctx context.Context, imageName string, imageRef string) error {
	reader, err := cl.cli.ImagePull(ctx, imageRef, image.PullOptions{}) // cli.ImagePull is asynchronous.
	if err != nil {
		return err
	}
	defer reader.Close()

	io.Copy(os.Stdout, reader)

	resp, err := cl.cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Cmd:   []string{"echo", "hello world"},
		Tty:   false,
	}, nil, nil, nil, "")
	if err != nil {
		return err
	}

	log.Printf("container created %s", resp.ID)

	if err := cl.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return err
	}

	out, err := cl.cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	return nil
}

func (cl *Client) Close() error {
	return cl.cli.Close()
}
