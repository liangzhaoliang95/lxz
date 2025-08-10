package docker_drivers

/*
import (
	"context"
	"fmt"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("ID: %s, Image: %s, Names: %v, Status: %s\n",
			container.ID[:12], container.Image, container.Names, container.Status)
	}
}
*/
