package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list is a tool for listing all projects",
	Long:  `list is a tool for listing all projects`,
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		nameFlag, _ := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}
		if !fileExists(filepath.Join(home, "/dev-env")) {
			log.Fatal("dev-env not initialized, run 'dev-env init' to initialize")
		}
		projectsJson, err := os.ReadFile(filepath.Join(home, "/dev-env/projects.json"))
		if err != nil {
			log.Fatal(err)
		}
		projects := []Project{}
		err = json.Unmarshal(projectsJson, &projects)
		if err != nil {
			log.Fatal(err)
		}
		statsJson, err := os.ReadFile(filepath.Join(home, "/dev-env/stats.json"))
		if err != nil {
			log.Fatal(err)
		}
		stats := map[string]int{}
	
		err = json.Unmarshal(statsJson, &stats)
		if err != nil {
			log.Fatal(err)
		}
		if nameFlag != "" {
			state := ""
			if _, ok := stats[nameFlag]; !ok {
				state = "[STOPPED]"
			}else {
				state = "[RUNNING]"
			}
			for _, project := range projects {
				if project.Name == nameFlag {
					fmt.Printf("%s - %s - %s\n", state, nameFlag, project.Path)
					return
				}
			}
			log.Fatal("Project not found")
		}
		for _, project := range projects {
			state := ""
			if _, ok := stats[project.Name]; !ok {
				state = "[STOPPED]"
			}else {
				state = "[RUNNING]"
			}
			fmt.Printf("%s - %s\n", state, project.Name)
		}
	},
}

func init() {
	listCmd.Flags().String("name", "", "name of the project to list")
	rootCmd.AddCommand(listCmd)
}