package boot

import (
	"errors"
	"github.com/EugeneNail/actum/internal/service/log"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
)

func LoadEnv() {
	envPath := filepath.Join(GetRootDirectory(), ".env")
	if err := godotenv.Load(envPath); err != nil {
		panic(err)
	}
}

func GetRootDirectory() string {
	directory, err := os.Getwd()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	for {
		_, err := os.Stat(filepath.Join(directory, "go.mod"))
		if err == nil {
			return directory
		}

		parentDirectory := filepath.Dir(directory)

		if parentDirectory == directory {
			err := errors.New("can't find the project root")
			log.Error(err)
			os.Exit(1)
		}

		directory = parentDirectory
	}
}
