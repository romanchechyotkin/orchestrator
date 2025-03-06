package task

import (
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID
	ContainerID string
	Name        string
	State       State

	Image         string
	Memory        int64
	Disk          int64
	Cpu           float64
	ExposedPorts  nat.PortSet
	PortBindings  map[string]string
	RestartPolicy container.RestartPolicyMode

	StartTime  time.Time
	FinishTime time.Time
}
