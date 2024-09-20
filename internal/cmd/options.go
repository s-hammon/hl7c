package cmd

import (
	"io"

	"github.com/s-hammon/hl7c/internal/config"
)

type Options struct {
	Stderr  io.Writer
	Tags    []string
	Against string
}

func (o *Options) ReadConfig(dir, filename string) (string, *config.Config, error) {
	return readConfig(o.Stderr, dir, filename)
}
