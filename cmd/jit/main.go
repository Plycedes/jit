package main

import (
	"fmt"
	"os"

	"github.com/plycedes/jit/internal/repo"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "jit",
		Short: "Jit - a Git implementation in Go",
		Long:  "Jit is a naive attempt made of a junior dev of implementing core Git features using Go",
	}

	rootCmd.AddCommand(initCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a new, empty Jit repository",
		RunE: func(cmd *cobra.Command, args []string) error {
			return repo.InitRepo(".")
		},
	}
}
