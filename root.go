package main

import (
	"github.com/spf13/cobra"
)
var rootCmd = &cobra.Command{
	Use: "dev-env",
	Short: "dev-env is a tool for managing your project",
	Long: `dev-env is a tool for managing your project`,
}

func Execute() {
	rootCmd.Execute()
}