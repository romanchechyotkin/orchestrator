package docker

import (
	"context"
	"testing"

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

	err = client.Run(context.Background(), "alpine", "docker.io/library/alpine")
	assert.NoError(t, err)
}
