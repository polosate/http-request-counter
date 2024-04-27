package main

import (
	"log"

	"simplesurance-test-task/internal/application"
)

func main() {
	app := application.New()

	if err := app.Init(); err != nil {
		log.Fatal(err)
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
