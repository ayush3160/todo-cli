package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
	models "todo-cli/models/tablelist"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	Register("complete", Complete)
}

func Complete(ctx context.Context, logger *zap.Logger) *cobra.Command {

	command := &cobra.Command{
		Use:   "complete",
		Short: "Marks a Task as completed",
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

			tableModel := models.NewTableModel(tableRows, len(tableRows), completeTask(logger, presentTasks))

			p := tea.NewProgram(tableModel)
			if _, err := p.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error starting program: %v", err)
				os.Exit(1)
			}
		},
	}

	return command
}

func completeTask(logger *zap.Logger, presentTasks []Task) func(taskID int) error {
	return func(taskID int) error {
		for i := range presentTasks {
			if i == taskID {
				presentTasks[i].Status = "Completed"
				presentTasks[i].UpdateDate = time.Now().Format("2006-01-02")
			}
		}

		newTaskFileContent, err := json.Marshal(presentTasks)
		if err != nil {
			logger.Error("Error marshalling tasks", zap.Error(err))
			return err
		}

		tasksFileLocation := GetTaskFileLocation()

		err = os.WriteFile(tasksFileLocation, newTaskFileContent, 0644)
		if err != nil {
			logger.Error("Error writing tasks file", zap.Error(err))
			return err
		}

		fmt.Println("Task marked as completed successfully!")

		return nil
	}
}
