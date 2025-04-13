package server

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/vzahanych/vm-api/core"
	"gopkg.in/yaml.v3"
)

var logger *slog.Logger
var config core.Config
var configPath string

func init() {
	// Add the server command to root
	rootCmd.AddCommand(serverCmd)

	// Add the --config flag for the config file path
	serverCmd.Flags().StringVarP(&configPath, "config", "c", "config.yaml", "Path to the configuration file")
}

var serverCmd = &cobra.Command{
	Use: "server",
	Run: func(cmd *cobra.Command, args []string) {

		if err := initConfig(configPath); err != nil {
			log.Fatalf("Error initializing configuration: %v", err)
		}

		// Initialize logger
		// Create an empty HandlerOptions
		handlerOptions := &slog.HandlerOptions{}
		level, err := parseLogLevel(config.Logging.LogLevel)
		if err != nil {
			log.Fatalf("Invalid log level: %v", err)
		}
		handlerOptions.Level = level
		logger = slog.New(slog.NewJSONHandler(os.Stdout, handlerOptions)) // Use standard library's slog

		r := gin.Default()

		// Enable CORS for all origins
		r.Use(cors.Default())

		// Custom logger middleware to use slog
		r.Use(core.CustomGinLogger(logger))

		r.POST("/vms", core.CreateVMHandler)       // Create VM
		r.DELETE("/vms/:id", core.DeleteVMHandler) // Delete VM
		r.GET("/vms/:id/status", core.GetVMStatus) // Get VM Status

		// Start the server
		srv := &http.Server{
			Addr:    config.Server.Address,
			Handler: r.Handler(),
		}

		go func() {
			// Start service connections
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Error("Failed to start server", "error", err) // Using logger here
			}
		}()

		// Wait for shutdown signal
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		logger.Info("Shutdown Server ...")

		ctx, cancel := context.WithTimeout(context.Background(), config.Server.GracefulShutdownDelay)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("Server Shutdown:", "error", err) // Using logger here
		}
		<-ctx.Done()
		logger.Info("Server exiting")
	},
}

func initConfig(path string) error {
	// Open the config file
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading config file, %s", err)
	}

	// Unmarshal the config into the Config struct
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("unable to decode into config struct, %s", err)
	}

	// Override server address with environment variable if set
	if serverAddress := os.Getenv("SERVER_ADDRESS"); serverAddress != "" {
		config.Server.Address = serverAddress
	}

	return nil
}

// parseLogLevel converts a string (e.g., "info", "debug") to a slog.Level
func parseLogLevel(level string) (slog.Level, error) {
	switch level {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelInfo, fmt.Errorf("unknown log level: %s", level)
	}
}
