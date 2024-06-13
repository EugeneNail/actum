package main

import (
	"fmt"
	"github.com/EugeneNail/actum/internal/service/env"
)

func main() {
	env.Load()

	placeholders := "?,"
	fmt.Println("(" + placeholders[:len(placeholders)-1] + ")")

	return

	//photosDirectory := filepath.Join(env.Get("APP_PATH"), "storage", "photos")
	//files, err := os.ReadDir(photosDirectory)
	//fmt.Println(files)
	//tests.Check(err)
	//
	//for _, file := range files {
	//	tests.Check(os.Remove(filepath.Join(photosDirectory, file.Name())))
	//}
}
