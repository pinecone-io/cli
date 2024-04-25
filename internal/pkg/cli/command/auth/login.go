package auth

import (
	"context"
	"embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/pinecone-io/cli/internal/pkg/utils/config"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

func NewLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Fetch your Pinecone API key via the browser",
		Run: func(cmd *cobra.Command, args []string) {
			startServer()
		},
	}

	return cmd
}

// Define a struct for the nested "api_key" object
type ConnectionData struct {
	Success bool                       `json:"success"`
	ApiKey  ConnectionDataApiKeyConfig `json:"key"`
}

// Define the main struct to match the JSON structure
type ConnectionDataApiKeyConfig struct {
	Id            string `json:"id"`
	Value         string `json:"value"`
	IntegrationId string `json:"integration_id"`
}

//go:embed callback.html
var content embed.FS

const ConnectUrl string = "https://app.pinecone.io/connect/cli"

func openBrowser(url string) error {
	var cmd string
	var args []string

	fmt.Println("Opening browser to authenticate with Pinecone...")

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}

	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func startServer() {
	srv := &http.Server{Addr: ":59049"}

	tmpl, err := template.ParseFS(content, "callback.html")
	if err != nil {
		panic(err) // Handle the error appropriately in production code
	}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		connectionData := r.URL.Query().Get("connectionData")
		decodedBytes, err := base64.StdEncoding.DecodeString(connectionData)
		if err != nil {
			exit.Error(err)
		}
		decodedJsonString := string(decodedBytes)

		var connectionDataJson ConnectionData
		jsonErr := json.Unmarshal([]byte(decodedJsonString), &connectionDataJson)
		if jsonErr != nil {
			exit.Error(fmt.Errorf("error parsing json: %e", err))
			return
		}

		apiKey := connectionDataJson.ApiKey.Value
		config.ApiKey.Set(apiKey)
		config.SaveConfig()
		fmt.Printf("✅ Successfully authenticated. Config property %s updated in %s\n", style.Emphasis("api_key"), style.Emphasis(config.NewConfigLocations().ConfigPath))

		// Data to pass to the template
		data := struct {
			ApiKey string
		}{
			ApiKey: apiKey,
		}

		// Serve the parsed HTML template with data
		w.Header().Set("Content-Type", "text/html")
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
		}

		// Optionally shutdown the server after handling the request
		go func() {
			if err := srv.Shutdown(context.Background()); err != nil {
				fmt.Printf("HTTP server Shutdown: %v", err)
			}
		}()
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			fmt.Println("Server forced to shutdown:", err)
		}
		fmt.Println("Server gracefully stopped.")
	}()

	fmt.Println("Starting temporary server process at http://localhost:59049 and waiting to receive authentication callback.")
	fmt.Printf("Please complete authentication flow in the browser at %s or press Ctrl-C to exit.\n", style.Emphasis(ConnectUrl))
	openBrowser(ConnectUrl)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Println("Error starting server:", err)
	}
}
