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
var projectSlug string
var projectHost string

func init() {
	rootCmd.AddCommand(projectCreateCmd)
	rootCmd.AddCommand(projectRemoveCmd)

	projectCreateCmd.Flags().StringVarP(&projectName, "name", "n", "", "Project name")
	projectCreateCmd.Flags().StringVarP(&projectHost, "host", "o", "", "Project host")

	projectRemoveCmd.Flags().StringVarP(&projectSlug, "slug", "s", "", "Project slug")
}

var projectCreateCmd = &cobra.Command{
	Use:   "project-create",
	Short: "Create or remove a project.",
	Long:  "Create or remove a project.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if projectName == "" {
			return errors.New("project name cannot be empty")
		}

		ps, err := getProjectService()
		if err != nil {
			return err
		}

		project, err := ps.CreateProject(projectName, projectHost)
		if err != nil {
			return err
		}

		fmt.Printf(
			"{\"project_slug\": \"%s\", \"api_key\": \"%s\"}",
			project.Slug,
			project.ApiKey,
		)

		return nil
	},
}

var projectRemoveCmd = &cobra.Command{
	Use:   "project-remove",
	Short: "Remove a project.",
	Long:  "Remove a project.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if projectSlug == "" {
			return errors.New("project slug cannot be empty")
		}

		ps, err := getProjectService()
		if err != nil {
			return err
		}

		if err := ps.RemoveProject(projectSlug); err != nil {
			return err
		}

		return nil
	},
}

func getProjectService() (proj.Service, error) {
	path := os.Getenv("DB_PATH")
	if path == "" {
		return nil, errors.New("DB_PATH env is not set")
	}

	conn, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	pr := proj.NewRepo(conn)
	ps := proj.NewService(pr)

	return ps, nil
}
