package main

import (
	"context"
	"log"

	"github.com/rizalarfiyan/skillshare-downloader/models"
	"github.com/rizalarfiyan/skillshare-downloader/services"
)

func main() {
	ctx := context.Background()
	err := services.NewSkillshare(ctx).Run(models.Config{
		UrlOrId:   "1088693386",
		SessionId: "92073dcb0cd0524343d7ee19c3971912",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
