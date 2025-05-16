package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase/tools/inflector"
)

func goBlankTemplate() string {
	return `package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// add up queries...

		return nil
	}, func(app core.App) error {
		// add down queries...

		return nil
	})
}
`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("file name not found")
	}

	name := inflector.Snakecase(strings.TrimSpace(os.Args[1]))
	if name == "" {
		log.Fatalln("file name empty")
	}

	dirPath := "migrations"
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		log.Fatalf("failed to create directory: %s", err.Error())
	}

	filePath := fmt.Sprintf("%s/%d_%s.go", dirPath, time.Now().Unix(), inflector.Snakecase(name))
	if err := os.WriteFile(filePath, []byte(goBlankTemplate()), 0644); err != nil {
		log.Fatalf("%s", err.Error())
	}

	fmt.Println("Created new file:", filePath)
	os.Exit(0)
}
