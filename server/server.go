package server

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)


type ServerProject struct {
	Name string `json:"name"`
	Cmd string `json:"cmd"`
	Path string `json:"path"`
	Running bool `json:"running"`
}

type Project struct {
	Name string `json:"name"`
	Cmd string `json:"cmd"`
	Path string `json:"path"`
}

func getFileSystem(files embed.FS) http.FileSystem {
	fsys, err := fs.Sub(files, "web/dist")
	if err != nil {
		panic(err)
	}
	return http.FS(fsys)
}

func Server(files embed.FS) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/api" , func(r chi.Router) {
		r.Get("/projects", func(w http.ResponseWriter, r *http.Request) {
			home, err := os.UserHomeDir()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			projectsBytes, err := os.ReadFile(filepath.Join(home, "dev-env-cli", "projects.json"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			projects := []Project{}
			json.Unmarshal(projectsBytes, &projects)
			statsBytes, err := os.ReadFile(filepath.Join(home, "dev-env-cli", "stats.json"))
			stats := map[string]int{}
			json.Unmarshal(statsBytes, &stats)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			serverProjects := []ServerProject{}
			for _, project := range projects {
				serverProjects = append(serverProjects, ServerProject{
					Name: project.Name,
					Cmd: project.Cmd,
					Path: project.Path,
					Running: stats[project.Name] != 0,
				})
			}
			json.NewEncoder(w).Encode(map[string]interface{}{
				"projects": serverProjects,
			},
			)
		},)
		r.Post("/projects/start", func(w http.ResponseWriter, r *http.Request) {
			nameObj := struct {
				Name string `json:"name"`
			}{}
			json.NewDecoder(r.Body).Decode(&nameObj)
			if nameObj.Name == "" {
				http.Error(w, "Name is required", http.StatusBadRequest)
				return
			}
			home, err := os.UserHomeDir()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			statsBytes, err := os.ReadFile(filepath.Join(home, "dev-env-cli", "stats.json"))
			stats := map[string]int{}
			json.Unmarshal(statsBytes, &stats)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if _, ok := stats[nameObj.Name]; ok {
				http.Error(w, "Project already running", http.StatusBadRequest)
				return
			}
			exec.Command("dev-env-cli", "start", nameObj.Name).Run()
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Project started",
			},
			)
		},)
		r.Post("/projects/stop", func(w http.ResponseWriter, r *http.Request) {
			nameObj := struct {
				Name string `json:"name"`
			}{}
			json.NewDecoder(r.Body).Decode(&nameObj)
			if nameObj.Name == "" {
				http.Error(w, "Name is required", http.StatusBadRequest)
				return
			}
			home, err := os.UserHomeDir()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			statsBytes, err := os.ReadFile(filepath.Join(home, "dev-env-cli", "stats.json"))
			stats := map[string]int{}
			json.Unmarshal(statsBytes, &stats)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if _, ok := stats[nameObj.Name]; !ok {
				http.Error(w, "Project not running", http.StatusBadRequest)
				return
			}
			exec.Command("dev-env-cli", "stop", nameObj.Name).Run()
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Project stopped",
			},
			)
	},)
	})

	// Serve embedded frontend for all non-API paths (including /assets/*)
	r.Handle("/*", http.FileServer(getFileSystem(files)))

	return r 
}