package cmd

import (
	"os"

	"github.com/MihaiBlebea/go-access-control/email"
	"github.com/MihaiBlebea/go-access-control/http"
	"github.com/MihaiBlebea/go-access-control/user"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the application server.",
	Long:  "Start the application server.",
	RunE: func(cmd *cobra.Command, args []string) error {

		l := logrus.New()

		l.SetFormatter(&logrus.JSONFormatter{})
		l.SetOutput(os.Stdout)
		l.SetLevel(logrus.InfoLevel)

		path, err := dbPathOrCreate()
		if err != nil {
			return err
		}

		conn, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
		if err != nil {
			return err
		}

		if err := conn.AutoMigrate(&user.User{}); err != nil {
			return err
		}

		es := email.New()

		ur := user.NewRepo(conn)
		us := user.NewService(ur, es)

		http.New(us, l)

		return nil
	},
}

func dbPathOrCreate() (string, error) {
	path := os.Getenv("DB_PATH")

	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	if _, err := os.Create(path); err != nil {
		return "", err
	}

	return path, nil
}
