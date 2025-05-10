package util

import (
	"github.com/joho/godotenv"
)

func LoadEnv(path string) error {
	err := godotenv.Load(path)
	return err
}
