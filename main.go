package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"todo-cli/cmd"

	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewDevelopment())
	defer logger.Sync()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tasksFileLocation := cmd.GetTaskFileLocation()

	if _, err := os.Stat(tasksFileLocation); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(tasksFileLocation), os.ModePerm); err != nil {
			logger.Error("Error creating tasks directory", zap.Error(err))
			return
		}
		_, err := os.Create(tasksFileLocation)
		if err != nil {
			logger.Error("Error creating tasks file", zap.Error(err))
			return
		}
	}

	rootCmd := cmd.Root(ctx, logger)

	if err := rootCmd.Execute(); err != nil {
		if strings.HasPrefix(err.Error(), "unknown command") || strings.HasPrefix(err.Error(), "unknown shorthand") {
			fmt.Println("Unknown command. Please run `task --help` for usage.")
		}
	}
}
