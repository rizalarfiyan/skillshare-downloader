package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/rizalarfiyan/skillshare-downloader/constants"
	"github.com/rizalarfiyan/skillshare-downloader/models"
	"github.com/rizalarfiyan/skillshare-downloader/utils"
)

type skillshare struct {
	ctx  context.Context
	conf models.AppConfig

	dir struct {
		base     string
		json     string
		video    string
		subtitle string
	}
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

	if err := s.initDir(); err != nil {
		return err
	}

	_, err := s.loadClassData()
	if err != nil {
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

func (s *skillshare) initDir() error {
	err := utils.CreateDir(s.conf.Dir)
	if err != nil {
		return err
	}

	dirs, err := utils.ReadDir(s.conf.Dir)
	if err != nil {
		return err
	}

	dirs = utils.Filter(dirs, func(dir string) bool {
		return strings.HasPrefix(dir, fmt.Sprintf("[%s]", s.conf.ID))
	})

	if len(dirs) != 1 {
		return nil
	}

	s.dir.base = path.Join(s.conf.Dir, dirs[0])
	err = s.loadDir()
	if err != nil {
		return err
	}
	return nil
}

func (s *skillshare) loadDir() error {
	if s.dir.base == "" {
		return nil
	}

	s.dir.json = path.Join(s.dir.base, "json")
	err := utils.CreateDir(s.dir.json)
	if err != nil {
		return err
	}

	s.dir.video = path.Join(s.dir.base, "video")
	err = utils.CreateDir(s.dir.video)
	if err != nil {
		return err
	}

	s.dir.subtitle = path.Join(s.dir.base, "subtitle")
	err = utils.CreateDir(s.dir.subtitle)
	if err != nil {
		return err
	}

	return nil
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

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("Skillshare class not found")
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

func (s *skillshare) createJsonClass(data models.ClassApiResponse) error {
	value, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	fileJson := path.Join(s.dir.json, constants.FilenameClassData)
	err = os.WriteFile(fileJson, value, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *skillshare) loadClassDataCache() (*models.ClassApiResponse, error) {
	if s.dir.base == "" || s.dir.json == "" {
		return nil, nil
	}

	fileJson := path.Join(s.dir.json, constants.FilenameClassData)
	if !utils.IsExistPath(fileJson) {
		return nil, nil
	}

	data, err := os.ReadFile(fileJson)
	if err != nil {
		return nil, err
	}

	dest := &models.ClassApiResponse{}
	err = json.Unmarshal(data, dest)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

func (s *skillshare) loadClassData() (*models.ClassApiResponse, error) {
	getCache, err := s.loadClassDataCache()
	if err != nil {
		return nil, err
	}

	if getCache != nil {
		return getCache, nil
	}

	getData, err := s.fetchClassApi()
	if err != nil {
		return nil, err
	}

	safeTitle := utils.SafeFolderName(getData.Title)
	folderName := fmt.Sprintf(constants.FolderName, s.conf.ID, safeTitle)
	s.dir.base = path.Join(s.conf.Dir, folderName)
	err = s.loadDir()
	if err != nil {
		return nil, err
	}

	err = s.createJsonClass(*getData)
	if err != nil {
		return nil, err
	}

	return getData, nil
}
