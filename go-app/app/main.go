package main

import (
	"log"

	"github.com/Danila331/hach-evroasia/app/servers"
)

func main() {
	err := servers.StartServer()
	if err != nil {
		log.Fatal(err)
	}
}
