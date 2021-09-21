package main

import (
	"github.com/MihaiBlebea/go-access-control/cmd"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("./.env")
}

func main() {
	cmd.Execute()
}
