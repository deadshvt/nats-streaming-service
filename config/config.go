package config

import (
	"github.com/joho/godotenv"
)

func Load(path string) {
	if err := godotenv.Load(path); err != nil {
		panic(err)
	}
}
