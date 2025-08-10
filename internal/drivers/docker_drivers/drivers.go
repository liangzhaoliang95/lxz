package docker_drivers

import (
	"context"
	"fmt"
	"github.com/liangzhaoliang95/lxz/internal/helper"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/image"
	"github.com/moby/moby/client"
	"io"
	"time"
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

func RemoveContainer(containerID string, force bool) error {
	cli, err := GetDockerClient()
	if err != nil {
		return fmt.Errorf("failed to get Docker client: %w", err)
	}
	if err = cli.ContainerRemove(context.Background(), containerID, container.RemoveOptions{
		Force: force,
	}); err != nil {
		return fmt.Errorf("failed to remove container %s: %w", containerID, err)
	}
	return nil
}

// StopContainer stops a running Docker container by its ID.
func StopContainer(containerID string, timeout *int) error {
	cli, err := GetDockerClient()
	if err != nil {
		return fmt.Errorf("failed to get Docker client: %w", err)
	}
	if err = cli.ContainerStop(context.Background(), containerID, container.StopOptions{
		Timeout: timeout,
	}); err != nil {
		return fmt.Errorf("failed to stop container %s: %w", containerID, err)
	}
	return nil
}

func InspectContainer(containerID string) (*container.InspectResponse, error) {
	cli, err := GetDockerClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get Docker client: %w", err)
	}
	containerInfo, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect container %s: %w", containerID, err)
	}
	return &containerInfo, nil
}

// InspectImage inspects a Docker image by its ID or name.
func InspectImage(imageID string) (*image.InspectResponse, error) {
	cli, err := GetDockerClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get Docker client: %w", err)
	}
	imageInfo, err := cli.ImageInspect(context.Background(), imageID)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect image %s: %w", imageID, err)
	}
	return &imageInfo, nil
}

func WaitContainerStopped(containerID string, timeout time.Duration) error {
	start := time.Now()
	for {
		inspect, err := InspectContainer(containerID)
		if err != nil {
			return err
		}

		// 检查状态不是 running
		if inspect.State != nil && inspect.State.Running == false {
			return nil
		}

		if time.Since(start) > timeout {
			return fmt.Errorf("timeout waiting for container %s to stop", containerID)
		}
		time.Sleep(500 * time.Millisecond) // 等待一小段时间再检查
	}
}
