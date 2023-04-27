package main

import (
	"context"

	"github.com/rizalarfiyan/skillshare-downloader/services"
)

func main() {
	ctx := context.Background()
	err := services.NewSkillshare(ctx).Run()
	if err != nil {
		panic(err)
	}
}
