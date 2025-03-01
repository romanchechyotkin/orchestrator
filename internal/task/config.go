package task

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

type Config struct {
	Name string

	AttachStdin  bool
	AttachStdout bool
	AttachStderr bool

	Image         string
	Cmd           []string
	ExposedPorts  nat.PortSet
	Env           []string
	RestartPolicy container.RestartPolicyMode

	Cpu    float64
	Memory int64
	Disk   int64
}
