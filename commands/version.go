package commands

import (
	"fmt"
	"io"
)

type Version struct {
	version []byte
	output  io.Writer
}

func NewVersion(version string, output io.Writer) *Version {
	return &Version{
		version: []byte(version + "\n"),
		output:  output,
	}
}

func (v Version) Execute([]string) error {
	_, err := v.output.Write(v.version)
	if err != nil {
		return fmt.Errorf("could not print version: %s", err)
	}

	return nil
}
