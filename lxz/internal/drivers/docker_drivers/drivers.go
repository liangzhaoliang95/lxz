package docker_drivers

import (
	"context"
	"fmt"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"io"
	"lxz/internal/helper"
)

var dockerClient *client.Client

func InitDockerClient() error {
	var err error
	dockerClient, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}
	return nil
}

func GetDockerClient() (*client.Client, error) {
	if dockerClient == nil {
		if err := InitDockerClient(); err != nil {
			return nil, err
		}
	}
	return dockerClient, nil
}

func ListContainers() ([]*container.Summary, error) {
	cli, err := GetDockerClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get Docker client: %w", err)
	}
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}
	list := make([]*container.Summary, 0, len(containers))
	for _, item := range containers {
		fmt.Println(fmt.Sprintf("%s", helper.Prettify(item)))
		list = append(list, &item)
	}
	return list, nil
}

func RestartContainer(containerID string, timeout *int) error {
	cli, err := GetDockerClient()
	if err != nil {
		return fmt.Errorf("failed to get Docker client: %w", err)
	}
	if err = cli.ContainerRestart(context.Background(), containerID, container.StopOptions{
		Timeout: timeout,
	}); err != nil {
		return fmt.Errorf("failed to restart container %s: %w", containerID, err)
	}
	return nil
}

func ContainerLogs(containerID string) (io.ReadCloser, error) {
	cli, err := GetDockerClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get Docker client: %w", err)
	}
	reader, err := cli.ContainerLogs(context.TODO(), containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       "100",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get logs for container %s: %w", containerID, err)
	}
	return reader, nil
}
