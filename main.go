package main

import (
	"log"
	"os"

	"github.com/ashwanthkumar/gotlb/providers"
)

func main() {
	log.Println("Starting gotlb ...")
	marathonHost := os.Args[1]

	provider := providers.NewMarathonProvider(marathonHost)
	NewManager().Start(provider)
}
