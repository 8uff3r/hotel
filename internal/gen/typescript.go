package main

import (
	"hotel/backend/internal/models"
	"os"

	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

func main() {
	converter := typescriptify.New()
	converter.CreateConstructor = false
	converter.CreateInterface = true
	converter.BackupDir = ""

	for _, v := range models.All() {
		converter.Add(v)
	}

	err := converter.ConvertToFile("web/app/utils/auto-route-types.ts")
	if err != nil {
		panic(err)
	}

	os.Exit(0)
}
