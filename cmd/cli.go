package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type HookFunc func(context.Context, *zap.Logger) *cobra.Command

var Registered map[string]HookFunc

type Task struct {
	ID           int    `json:"id"`
	Description  string `json:"description"`
	Status       string `json:"status"`
	CreationDate string `json:"creation_date"`
	UpdateDate   string `json:"update_date"`
}

func Register(name string, f HookFunc) {
	if Registered == nil {
		Registered = make(map[string]HookFunc)
	}
	Registered[name] = f
}

func Root(ctx context.Context, logger *zap.Logger) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "task",
		Short: "Task is a CLI task manager",
	}

	for _, f := range Registered {
		rootCmd.AddCommand(f(ctx, logger))
	}

	return rootCmd
}
