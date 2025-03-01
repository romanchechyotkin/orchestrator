package docker

import (
	"context"
	"io"
	"log"
	"math"
	"os"

	"github.com/romanchechyotkin/orchestrator/internal/task"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type Result struct {
	Error       error
	Action      string
	ContainerID string
	Result      string
}

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

func (cl *Client) Run(ctx context.Context, cfg *task.Config) *Result {
	reader, err := cl.cli.ImagePull(ctx, cfg.Image, image.PullOptions{}) // cli.ImagePull is asynchronous.
	if err != nil {
		return &Result{
			Error: err,
		}
	}
	defer reader.Close()

	io.Copy(os.Stdout, reader)

	rp := container.RestartPolicy{
		Name: cfg.RestartPolicy,
	}

	r := container.Resources{
		Memory:   cfg.Memory,
		NanoCPUs: int64(cfg.Cpu * math.Pow(10, 9)),
	}

	cc := container.Config{
		Image:        cfg.Image,
		Tty:          false,
		Env:          cfg.Env,
		Cmd:          cfg.Cmd,
		ExposedPorts: cfg.ExposedPorts,
	}

	hc := container.HostConfig{
		RestartPolicy:   rp,
		Resources:       r,
		PublishAllPorts: true,
	}

	resp, err := cl.cli.ContainerCreate(ctx, &cc, &hc, nil, nil, cfg.Name)
	if err != nil {
		log.Println("failed to create container", err)
		return &Result{
			Error: err,
		}
	}

	log.Printf("container created %s", resp.ID)

	if err := cl.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return &Result{
			Error: err,
		}
	}

	out, err := cl.cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		log.Printf("failed to get logs for container %s: %v\n", resp.ID, err)
		return &Result{
			Error: err,
		}
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	return &Result{
		ContainerID: resp.ID,
		Action:      "start",
		Result:      "success",
	}
}

func (cl *Client) Stop(containerID string) *Result {
	log.Printf("attempting to stop container %s\n", containerID)
	ctx := context.Background()

	if err := cl.cli.ContainerStop(ctx, containerID, container.StopOptions{}); err != nil {
		log.Printf("failed to stop container %s: %v\n", containerID, err)
		return &Result{
			Error: err,
		}
	}

	if err := cl.cli.ContainerRemove(ctx, containerID, container.RemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         false,
	}); err != nil {
		log.Printf("failed to remove container %s: %v\n", containerID, err)
		return &Result{
			Error: err,
		}
	}

	return &Result{
		Action: "stop",
		Result: "success",
	}
}

func (cl *Client) Close() error {
	return cl.cli.Close()
}
