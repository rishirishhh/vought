package main

import (
	"time"

	"github.com/rishirishhh/vought/src/cmd/api/config"
	log "github.com/sirupsen/logrus"
)

const GORILLA_MUX_SHUTDOWN_TIMEOUT time.Duration = time.Second * 2
const GOROUTINE_FLUSH_TIMEOUT time.Duration = time.Millisecond * 100

func main() {
	log.Info("Starting Vought API")

	// Retrieve environment variables
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Failed to parse environment variables : ", err)
	}

	if cfg.DevMode {
		log.SetReportCaller(true)
		log.SetLevel(log.DebugLevel)
	}

}
