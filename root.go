package main

import (
	"github.com/spf13/cobra"
)
var rootCmd = &cobra.Command{
	Use: "dev-env-cli",
	Short: "dev-env-cli is a tool for managing your project",
	Long: `dev-env-cli is a tool for managing your project`,
}

func Execute() {
	rootCmd.Execute()
}