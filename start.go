package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	syscall "golang.org/x/sys/windows"
)

var startCmd = &cobra.Command{
	Use:   "start [project-name]",
	Short: "start is a tool for starting a project",
	Long:  `start is a tool for starting a project`,
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		if len(args) == 0 {
			log.Fatal("project name is required")
		}
		projectName := args[0]
		if projectName == "" {
			log.Fatal("project name is required")
		}
		if err != nil {
			log.Fatal(err)
		}
		if !fileExists(filepath.Join(home, "/dev-env-cli")) {
			log.Fatal("dev-env-cli not initialized, run 'dev-env-cli init' to initialize")
		}
		projectsJson, err := os.ReadFile(filepath.Join(home, "/dev-env-cli/projects.json"))
		if err != nil {
			log.Fatal(err)
		}
		projects := []Project{}
		err = json.Unmarshal(projectsJson, &projects)
		if err != nil {
			log.Fatal(err)
		}
		statsJson, err := os.ReadFile(filepath.Join(home, "/dev-env-cli/stats.json"))
		if err != nil {
			log.Fatal(err)
		}
		stats := map[string]int{}
		err = json.Unmarshal(statsJson, &stats)
		if err != nil {
			log.Fatal(err)
		}
		for _, project := range projects {
			if project.Name == projectName {
				parts := strings.Fields(project.Cmd)
				if stats[projectName] != 0 {
					log.Fatal("Project already running")
				}
				if err != nil {
					log.Fatal(err)
				}
				if len(parts) == 0 {
					log.Fatal("project command is empty")
				}
				c := exec.Command(parts[0], parts[1:]...)
				c.Dir = project.Path
				c.SysProcAttr = &syscall.SysProcAttr{
					CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP | syscall.DETACHED_PROCESS,
				}
				if err := c.Start(); err != nil {
					log.Fatal(err)
				}
				stats[projectName] = c.Process.Pid
				statsBytes, err := json.Marshal(stats)
				if err != nil {
					log.Fatal(err)
				}
				os.WriteFile(filepath.Join(home, "/dev-env-cli/stats.json"), statsBytes, 0644)
				fmt.Printf("Successfully started '%s'\n", project.Name)
				return
			}
		}
		log.Fatal("Project not found")
	},
}

func init(){
	rootCmd.AddCommand(startCmd)
}