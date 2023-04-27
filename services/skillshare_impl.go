package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rizalarfiyan/skillshare-downloader/constants"
	"github.com/rizalarfiyan/skillshare-downloader/models"
)

type skillshare struct {
	ctx  context.Context
	conf models.AppConfig
}

func NewSkillshare(ctx context.Context) Skillshare {
	return &skillshare{
		ctx: ctx,
	}
}

func (s *skillshare) Run(conf models.Config) error {
	s.splash()
	if err := s.conf.FromConfig(conf); err != nil {
		return err
	}
	if err := s.loadClassData(); err != nil {
		return err
	}
	return nil
}

func (s *skillshare) splash() {
	text := `     ____    _      _   _   _         _                                 ____    _     
    / ___|  | | __ (_) | | | |  ___  | |__     __ _   _ __    ___      |  _ \  | |    
    \___ \  | |/ / | | | | | | / __| | '_ \   / _` + "`" + ` | | '__|  / _ \     | | | | | |    
     ___) | |   <  | | | | | | \__ \ | | | | | (_| | | |    |  __/     | |_| | | |___ 
    |____/  |_|\_\ |_| |_| |_| |___/ |_| |_|  \__,_| |_|     \___|     |____/  |_____|
    
                                        By Rizal Arfiyan                                     
    `

	fmt.Printf("\n%s\n\n", text)
}

func (s *skillshare) fetchClassApi() (*models.ClassApiResponse, error) {
	client := &http.Client{}
	url := fmt.Sprintf(constants.APIClass, s.conf.ID)
	req, err := http.NewRequestWithContext(s.ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Accept":     {"application/vnd.skillshare.class+json;,version=0.8"},
		"User-Agent": {"Skillshare/5.3.0; Android 9.0.1"},
		"Host":       {"api.skillshare.com"},
		"Referer":    {"https://www.skillshare.com/"},
		"cookie":     {fmt.Sprintf("PHPSESSID=%s", s.conf.SessionId)},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	dest := &models.ClassApiResponse{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, dest)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

func (s *skillshare) loadClassData() error {
	dest, err := s.fetchClassApi()
	if err != nil {
		return err
	}

	fmt.Println(dest)

	return nil
}
