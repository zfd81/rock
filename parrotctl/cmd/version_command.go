package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

const (
	Version = "1.0"
)

func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints the version of parrotctl",
		Run:   versionCommandFunc,
	}
}

func versionCommandFunc(cmd *cobra.Command, args []string) {
	fmt.Println("parrotctl version:", Version)
}
