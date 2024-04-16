package env

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func Load() {
	directory := getRootDirectory()
	err := os.Setenv("APP_PATH", directory)
	check(err)
	file, err := os.Open(filepath.Join(directory, ".env"))
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	regex := regexp.MustCompile("^[a-zA-Z0-9_]+=")

	for scanner.Scan() {
		if !regex.MatchString(scanner.Text()) {
			continue
		}
		tuple := strings.SplitN(scanner.Text(), "=", 2)
		err := os.Setenv(tuple[0], tuple[1])
		check(err)
	}
}

func getRootDirectory() string {
	directory, err := os.Getwd()
	check(err)

	for {
		_, err := os.Stat(filepath.Join(directory, "go.mod"))
		if err == nil {
			return directory
		}
		parentDirectory := filepath.Dir(directory)

		if directory == parentDirectory {
			panic("could not find the root directory")
		}
		directory = parentDirectory
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
