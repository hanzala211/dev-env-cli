package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !errors.Is(err, os.ErrNotExist)
}

var initCmd = &cobra.Command{
	Use: "init",
	Short: "init is a tool for initializing tool for development",
	Long: `init is a tool for initializing tool for development`,
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()

		if fileExists(filepath.Join(home, "/dev-env-cli")) {
			log.Fatal("dev-env-cli already initialized")
		}
		if err != nil {
			log.Fatal(err)
		}
		os.MkdirAll(filepath.Join(home, "/dev-env-cli"), 0755)
		os.WriteFile(filepath.Join(home, "/dev-env-cli/projects.json"), []byte("[]"), 0644)
		os.WriteFile(filepath.Join(home, "/dev-env-cli/stats.json"), []byte("{}"), 0644)
		fmt.Printf("Initialized dev-env-cli in %s\n", filepath.Join(home, "/dev-env-cli"))
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}