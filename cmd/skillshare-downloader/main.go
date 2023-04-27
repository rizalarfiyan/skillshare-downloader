package main

import (
	"context"

	"github.com/rizalarfiyan/skillshare-downloader/logger"
	"github.com/rizalarfiyan/skillshare-downloader/models"
	"github.com/rizalarfiyan/skillshare-downloader/services"
	"github.com/sirupsen/logrus"
)

func init() {
	logger.Init()
}

func main() {
	log := logger.Get()
	defer func() {
		if rec := recover(); rec != nil {
			log.Fatalln("Panic: ", rec)
		}
	}()

	// update with cli
	logger.SetLevel(logrus.DebugLevel)

	ctx := context.Background()
	err := services.NewSkillshare(ctx).Run(models.Config{
		UrlOrId:   "1088693386",
		SessionId: "92073dcb0cd0524343d7ee19c3971912",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
