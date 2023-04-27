package services

import "github.com/rizalarfiyan/skillshare-downloader/models"

type Skillshare interface {
	Run(conf models.Config) error
}
