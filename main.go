package main

import "github.com/utilyre/golang-backend-template/application"

func main() {
	app := application.New()
	app.Init()
	app.Start()
}
