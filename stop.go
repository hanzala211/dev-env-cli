package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop [project-name]",
	Short: "stop is a tool for stopping a project",
	Long:  `stop is a tool for stopping a project`,
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		if len(args) == 0 {
			log.Fatal("Project name is required")
		}
		projectName := args[0]
		if projectName == "" {
			log.Fatal("Project name is required")
		}
		if err != nil {
			log.Fatal(err)
		}
		if !fileExists(filepath.Join(home, "/dev-env")) {
			log.Fatal("dev-env not initialized, run 'dev-env init' to initialize")
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
		projectsJson, err := os.ReadFile(filepath.Join(home, "/dev-env/projects.json"))
		if err != nil {
			log.Fatal(err)
		}
		projects := []Project{}
		err = json.Unmarshal(projectsJson, &projects)
		if err != nil {
			log.Fatal(err)
		}
		for _, project := range projects {
			if project.Name == projectName {
				if _, ok := stats[projectName]; !ok {
					log.Fatal("Project not running")
				}
				if runtime.GOOS == "windows" {
					cmd := exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(stats[projectName]))
					cmd.Run()
					if err != nil {
						log.Fatal(err)
					}
				} else {
					process, err := os.FindProcess(stats[projectName])
					if err != nil {
						log.Fatal(err)
					}
					err = process.Kill()
					if err != nil {
						log.Fatal(err)
					}
				}
				delete(stats, projectName)
				statsBytes, err := json.Marshal(stats)
				if err != nil {
					log.Fatal(err)
				}
				os.WriteFile(filepath.Join(home, "/dev-env/stats.json"), statsBytes, 0644)
				fmt.Printf("Successfully stopped '%s'\n", projectName)
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}