package main

import (
	"log"
	"os"

	"github.com/ashwanthkumar/gotlb/providers"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC | log.Lshortfile)
	log.SetOutput(os.Stdout)

	log.Println("Starting gotlb ...")
	marathonHost := os.Args[1]

	provider := providers.NewMarathonProvider(marathonHost)
	NewManager().Start(provider)
}
