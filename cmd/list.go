package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	models "todo-cli/models/tablelist"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	Register("list", List)
}

func List(ctx context.Context, logger *zap.Logger) *cobra.Command {

	command := &cobra.Command{
		Use:   "list",
		Short: "List all tasks",
		Run: func(cmd *cobra.Command, args []string) {

			tasksFileLocation := GetTaskFileLocation()

			taskFileContent, err := os.ReadFile(tasksFileLocation)

			if err != nil {
				logger.Error("Error reading tasks file", zap.Error(err))
				return
			}

			var presentTasks []Task

			if len(taskFileContent) > 0 {
				err = json.Unmarshal(taskFileContent, &presentTasks)
				if err != nil {
					logger.Error("Error unmarshalling tasks", zap.Error(err))
					return
				}
			}

			if len(presentTasks) == 0 {
				fmt.Println("No tasks found")
				return
			}

			tableRows := []table.Row{}

			for _, task := range presentTasks {
				tableRows = append(tableRows, table.Row{
					strconv.Itoa(task.ID),
					task.Description,
					task.Status,
					task.CreationDate,
					task.UpdateDate,
				})
			}

			tableModel := models.NewTableModel(tableRows, len(tableRows), func(i int) error { return nil })

			p := tea.NewProgram(tableModel)
			if _, err := p.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error starting program: %v", err)
				os.Exit(1)
			}
		},
	}

	return command
}
