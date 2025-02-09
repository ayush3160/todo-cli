package cmd

import (
	"os"
	"path"
)

func GetTaskFileLocation() string {

	tempDir := os.TempDir()

	taskFileLocation := path.Join(tempDir, "/todo-cli/tasks.json")

	return taskFileLocation
}
