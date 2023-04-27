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
		UrlOrId:   "9999999999",
		SessionId: "session_id",
	})
	if err != nil {
		log.Fatalln(err)
	}
}
