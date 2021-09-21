package cmd

import (
	"errors"
	"fmt"
	"os"

	proj "github.com/MihaiBlebea/go-access-control/project"
	"github.com/spf13/cobra"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var projectName string
var projectHost string

func init() {
	rootCmd.AddCommand(createProjectCmd)

	createProjectCmd.Flags().StringVarP(&projectName, "name", "n", "", "Project name")
	createProjectCmd.Flags().StringVarP(&projectHost, "host", "h", "", "Project host")
}

var createProjectCmd = &cobra.Command{
	Use:   "create-project",
	Short: "Create a new project.",
	Long:  "Create a new project.",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := os.Getenv("DB_PATH")
		if path == "" {
			return errors.New("DB_PATH env is not set")
		}

		if projectName == "" {
			return errors.New("project name cannot be empty")
		}

		conn, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
		if err != nil {
			return err
		}

		pr := proj.NewRepo(conn)
		ps := proj.NewService(pr)

		project, err := ps.CreateProject(projectName, projectHost)
		if err != nil {
			return err
		}

		fmt.Printf("Project slug: %s", project.Slug)

		return nil
	},
}
