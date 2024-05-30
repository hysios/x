package container

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/hysios/go-dexec"
)

type Container struct {
	dock *dexec.Docker
}

type Cmd struct {
}

type Config struct {
	Image string
	Cli   *docker.Client
}

// NewContainer
// Create a new container
func NewContainer(cfg *Config) (*Container, error) {
	if cfg.Cli == nil {
		cl, err := docker.NewClientFromEnv()
		if err != nil {
			return nil, err
		}
	}

	m, err := dexec.ByCreatingContainer(docker.CreateContainerOptions{
		Config: &docker.Config{Image: cfg.Image}})

	if err != nil {
		return nil, err
	}

	return &Container{dock: m}, nil
}

// Exec
// Execute a command inside a container
func (c *Container) Exec(cmd string, args ...string) error {
	panic("nonimplement")
}
