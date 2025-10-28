package main

import (
	"github.com/rishirishhh/vought/src/cmd/api/config"
	"github.com/rishirishhh/vought/src/cmd/encoder/eventhandler"
	"github.com/rishirishhh/vought/src/pkg/clients"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Starting Vought Encoder")

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Cannot parse Env var ", err)
	}

	if cfg.DevMode {
		log.SetLevel(log.DebugLevel)
	}

	// S3 client to access the videos
	s3Client, err := clients.NewS3Client(cfg.S3Host, cfg.S3Region, cfg.S3Bucket, cfg.S3AuthKey, cfg.S3AuthPwd)
	if err != nil {
		log.Fatal("Fail to create S3Client ", err)
	}
	amqpClientVideoUpload, _ := clients.NewAmqpClient(cfg.RabbitmqUser, cfg.RabbitmqPwd, cfg.RabbitmqAddr)

	// Listen, consume and publish on amqpClientVideoUpload
	eventhandler.ConsumeEvents(amqpClientVideoUpload, s3Client)

}
