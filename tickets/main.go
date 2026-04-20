package main

import (
	"tickets/app/routes"
	"tickets/core"
)

func main() {
	app := core.New()
	defer app.Close()

	routes.Init(app)
	app.Run()
}