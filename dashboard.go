package main

import (
	"embed"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/hanzala211/dev-env-cli/server"
	"github.com/spf13/cobra"
)

//go:embed web/dist
var files embed.FS

var dashboardCmd = &cobra.Command{
	Use: "dashboard",
	Short: "dashboard is a tool for starting a dashboard",
	Long: `dashboard is a tool for starting a dashboard`,
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		if !fileExists(filepath.Join(home, "/dev-env-cli")) {
			log.Fatal("dev-env-cli not initialized, run 'dev-env-cli init' to initialize")
		}

		listener, listenErr := net.Listen("tcp", ":8080")
		if listenErr != nil {
			log.Println("Dashboard is already running at http://localhost:8080")
			_ = openBrowser("http://localhost:8080")
			return
		}

        s := server.Server(files)
		log.Println("Dashboard started on port 8080")
		_ = openBrowser("http://localhost:8080")
		log.Fatal(http.Serve(listener, s))
	},
}

func init() { 
	rootCmd.AddCommand(dashboardCmd)
}

func openBrowser(url string) error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return exec.Command("xdg-open", url).Start()
	}
}

