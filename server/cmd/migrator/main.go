package main

import (
	"dankmuzikk/log"
	"dankmuzikk/mariadb"
)

func main() {
	return
	err := mariadb.Migrate2()
	if err != nil {
		log.Fatalln(err)
	}
}
