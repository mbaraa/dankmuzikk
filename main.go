package main

import (
	"dankmuzikk/cmd/migration"
	"dankmuzikk/cmd/server"
	"dankmuzikk/log"
	"embed"
	"os"
)

//go:embed static/*
var static embed.FS

//go:generate npx tailwindcss build -i static/css/style.css -o static/css/tailwind.css -m

func main() {
	var err error
	switch os.Args[1] {
	case "serve", "server":
		err = server.StartServer(static)
	case "migrate", "migration", "theotherthing":
		err = migration.Migrate()
	}
	if err != nil {
		log.Fatalln(log.ErrorLevel, err)
	}
}
