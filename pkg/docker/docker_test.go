package docker

import (
	"context"
	"testing"

	"github.com/romanchechyotkin/orchestrator/internal/task"

	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
)

func Test_NewClient(t *testing.T) {
	client, err := NewClient(
		client.WithHost("unix:///Users/rchechetkin/.colima/ozon/docker.sock"),
		client.WithAPIVersionNegotiation(),
	)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func Test_Run(t *testing.T) {
	client, err := NewClient(
		client.WithHost("unix:///Users/rchechetkin/.colima/ozon/docker.sock"),
		client.WithAPIVersionNegotiation(),
	)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	cfg := &task.Config{
		Name:  "test-container",
		Image: "alpine",
		Cmd:   []string{"echo", "hello world"},
	}

	res := client.Run(context.Background(), cfg)
	assert.NoError(t, res.Error)

	res = client.Stop(res.ContainerID)
	assert.NoError(t, res.Error)
}
