package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	models "todo-cli/models/textarea"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	Register("add", Add)
}

func Add(ctx context.Context, logger *zap.Logger) *cobra.Command {

	command := &cobra.Command{
		Use:   "add",
		Short: "Add a new task",
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

			p := tea.NewProgram(models.NewTextAreaModel(addTask(logger, presentTasks)))
			if _, err := p.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v", err)
				os.Exit(1)
			}
		},
	}

	return command
}

func addTask(logger *zap.Logger, presentTasks []Task) func(string) error {
	return func(taskDescription string) error {

		fmt.Println("Adding task: ", taskDescription)

		tasksFileLocation := GetTaskFileLocation()

		task := Task{
			ID:           len(presentTasks) + 1,
			Description:  taskDescription,
			Status:       "Pending",
			CreationDate: time.Now().Format("2006-01-02"),
			UpdateDate:   time.Now().Format("2006-01-02"),
		}

		presentTasks = append(presentTasks, task)

		newTaskFileContent, err := json.Marshal(presentTasks)
		if err != nil {
			logger.Error("Error marshalling tasks", zap.Error(err))
			return err
		}

		err = os.WriteFile(tasksFileLocation, newTaskFileContent, 0644)
		if err != nil {
			logger.Error("Error writing tasks file", zap.Error(err))
			return err
		}

		return nil
	}
}
