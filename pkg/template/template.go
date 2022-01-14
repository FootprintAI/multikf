package template

import "io"

type TemplateExecutor interface {
	Filename() string
	Execute(io.Writer) error
	Populate(config *TemplateFileConfig) error
}

// TemplateFileConfig is a union template file config
type TemplateFileConfig struct {
	Name string
	// NOTE: only support virtualbox now

	CPUs   int // number of cpus allocated
	Memory int // number of bytes memory allocated

	// NOTE: GPUs are not supported now
	// GPUs string

	SSHPort     int
	KubeApiPort int
}
