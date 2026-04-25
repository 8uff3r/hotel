package main

import (
	"hotel/internal/models"
	"log"
	"os"
	"strings"

	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

const distFile = "wails/frontend/app/utils/auto-route-types.ts"

func main() {
	converter := typescriptify.New()
	converter.CreateConstructor = false
	converter.CreateInterface = true
	converter.BackupDir = ""

	for _, v := range models.AllForTypeGen() {
		converter.Add(v)
	}

	err := converter.ConvertToFile(distFile)
	if err != nil {
		panic(err)
	}
	data, err := os.ReadFile(distFile)
	if err != nil {
		log.Fatal(err)
	}

	out := strings.ReplaceAll(string(data), " id: number", " id?: number")

	err = os.WriteFile(distFile, []byte(out), 0644)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
