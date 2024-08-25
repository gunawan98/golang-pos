package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

func LoadEnvVars() {
	cwd, _ := os.Getwd()
	dirString := strings.Split(cwd, "golang-pos")
	dir := strings.Join([]string{dirString[0], "golang-pos"}, "")
	AppPath := dir
	err := godotenv.Load(filepath.Join(AppPath, "/.env"))
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
