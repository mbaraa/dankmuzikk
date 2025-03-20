package main

import (
	"dankmuzikk/log"
	"dankmuzikk/mariadb"
)

func main() {
	err := mariadb.Migrate()
	if err != nil {
		log.Fatalln(err)
	}
}
