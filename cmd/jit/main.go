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
	rootCmd.AddCommand(hashObjectCmd())
	rootCmd.AddCommand(catFileCmd())
	rootCmd.AddCommand(addCmd())

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

func addCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <file>",
		Short: "Add file contents to the index (stage them for commit)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return repo.Add("./.jit", args[0])
		},
	}
}

func hashObjectCmd() *cobra.Command {
	var write bool
	cmd := &cobra.Command{
		Use:   "hash-object [flags] <file>",
		Short: "Compute object id (sha1) of a file. Use -w to write to object store.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			sha, err := repo.HashObject("./.jit", data, "blob", write)
			if err != nil {
				return err
			}
			fmt.Println(sha)
			return nil
		},
	}
	cmd.Flags().BoolVarP(&write, "write", "w", false, "Write the object into the object database")
	return cmd
}

func catFileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cat-file -p <sha1>",
		Short: "Pretty-print object contents (only -p supported)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sha := args[0]
			typ, content, err := repo.ReadObject("./.jit", sha)
			if err != nil {
				return err
			}

			fmt.Printf("type: %s\n", typ)
			os.Stdout.Write(content)
			return nil
		},
	}
	return cmd
}
