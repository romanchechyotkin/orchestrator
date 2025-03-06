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

func NewConfig(t *Task) *Config {
	return &Config{
		Name:          t.Name,
		Image:         t.Image,
		Memory:        t.Memory,
		Cpu:           t.Cpu,
		Disk:          t.Disk,
		ExposedPorts:  t.ExposedPorts,
		RestartPolicy: t.RestartPolicy,
	}
}
