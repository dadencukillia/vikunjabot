package main

import (
	"log"
	"vikunjabot/internal"
)

func main() {
	_, err := internal.GetConfig()
	if err != nil {
		log.Panic(err)
	}
}
