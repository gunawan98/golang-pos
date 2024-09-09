package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

var projectName string = "golang-pos"

func LoadEnvVars() {
	cwd, _ := os.Getwd()
	dirString := strings.Split(cwd, projectName)
	dir := strings.Join([]string{dirString[0], projectName}, "")
	AppPath := dir
	err := godotenv.Load(filepath.Join(AppPath, "/.env"))
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
