package main

import (
	"log"
	"os"
	"time"

	"github.com/ashwanthkumar/gotlb/providers"
	"github.com/rcrowley/go-metrics"
)

// DebugMetricsRegistry is used for pushing debug level metrics by rest of the app
var MetricsRegistry metrics.Registry

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC | log.Lshortfile)
	log.SetOutput(os.Stdout)

	MetricsRegistry = metrics.NewPrefixedRegistry("gotlb-")
	go metrics.Log(MetricsRegistry, 60*time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))

	log.Println("Starting gotlb ...")
	marathonHost := os.Args[1]

	provider := providers.NewMarathonProvider(marathonHost)
	NewManager().Start(provider)
}
