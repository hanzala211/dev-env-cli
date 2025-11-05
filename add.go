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
)

type Project struct {
	Name string `json:"name"`
	Cmd string `json:"cmd"`
	Path string `json:"path"`
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add is a tool for adding a new project",
	Long:  `add is a tool for adding a new project`,
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		name, _ := cmd.Flags().GetString("name")
		cmdFlag, _ := cmd.Flags().GetString("cmd")
		pathFlag, _ := cmd.Flags().GetString("path")
		if name == "" {
			log.Fatal("name is required")
		}

		if len(args) > 0 {
			if cmdFlag == "" {
				cmdFlag = strings.Join(args, " ")
			} else {
				cmdFlag = strings.TrimSpace(cmdFlag + " " + strings.Join(args, " "))
			}
		}
		if cmdFlag == "" {
			log.Fatal("cmd is required, e.g. --cmd \"npm run dev\" or use: add --name X -- npm run dev")
		}
		if pathFlag == "" {
			pathFlag, _ = os.Getwd()
		}

		project := Project{
			Name: name,
			Cmd: cmdFlag,
			Path: pathFlag,
		}
		if !fileExists(filepath.Join(home, "/dev-env-cli")) {
			if err := exec.Command("go", "run", ".", "init").Run(); err != nil {
				log.Fatal(err)
			}
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

		for _, val := range projects {
			if val.Name == name || val.Path == pathFlag {
				log.Fatal("Project already exists")
			}
		}
		projects = append(projects, project)
		projectsBytes, err := json.Marshal(projects)
		if err != nil {
			log.Fatal(err)
		}
		os.WriteFile(filepath.Join(home, "/dev-env-cli/projects.json"), projectsBytes, 0644)
		fmt.Printf("Project %s added successfully\n", name)
	},
}

func init() {
	addCmd.Flags().String("name", "", "name of the project (required)")
	addCmd.Flags().String("cmd", "", "command to run the project (required) e.g. 'npm run dev'")
	addCmd.Flags().String("path", "", "path to the project (optional, will use the current directory if not provided)")
	rootCmd.AddCommand(addCmd)
}