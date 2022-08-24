package component

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func TestXxx(t *testing.T) {
	logger := GetLogger()
	fmt.Println(logger)
}

func TestFoo(t *testing.T) {
	ids, _ := os.Getwd()
	fmt.Println(ids)
}

func TestDockerCli(t *testing.T) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
}


func TestQuery(t *testing.T) {
	SingleQuery()
}