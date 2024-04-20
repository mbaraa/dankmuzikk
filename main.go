package main

import (
	"dankmuzikk/services/youtube"
	"fmt"
)

func main() {
	vids, err := youtube.Search("lana del rey raise me up")
	if err != nil {
		panic(err)
	}

	for _, vid := range vids {
		fmt.Printf("%+v\n", vid)
	}
}
