package docker

import "fmt"

func NewContainerName(clustername string) ContainerName {
	return ContainerName{clustername: clustername}
}

// ContainerName is a place holder for clustername and its underlying container name created by kind
type ContainerName struct {
	clustername string
}

func (c ContainerName) Name() string {
	return fmt.Sprintf("/%s-control-plane", c.clustername)

}
