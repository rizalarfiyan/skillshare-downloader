package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/rizalarfiyan/skillshare-downloader/constants"
	"github.com/rizalarfiyan/skillshare-downloader/logger"
	"github.com/rizalarfiyan/skillshare-downloader/models"
	"github.com/rizalarfiyan/skillshare-downloader/utils"
	"github.com/sirupsen/logrus"
)

type skillshare struct {
	ctx  context.Context
	conf models.AppConfig
	log  *logrus.Logger

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
		log: logger.Get(),
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

	ssClass, err := s.loadClassData()
	if err != nil {
		return err
	}

	ssData, err := s.loadVideoData(*ssClass)
	if err != nil {
		return err
	}

	fmt.Println(ssData)

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
		return strings.HasPrefix(dir, fmt.Sprintf("[%d]", s.conf.ID))
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

func (s *skillshare) fetchClassApi() (*models.ClassData, error) {
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

	dest := &models.ClassData{}
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

func (s *skillshare) fetchVideoApi(videoID int) (*models.VideoData, error) {
	fmt.Printf("fetching video api %d\n", videoID)
	client := &http.Client{}
	url := fmt.Sprintf(constants.APIVideo, constants.BrightcoveAccountId, videoID)
	req, err := http.NewRequestWithContext(s.ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Accept":     {fmt.Sprintf("application/json;pk=%s", constants.PolicyKey)},
		"User-Agent": {"Mozilla/5.0 (X11; Linux x86_64; rv:52.0) Gecko/20100101 Firefox/52.0"},
		"Origin":     {"https://www.skillshare.com/"},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("Skillshare video not found")
	}

	dest := &models.VideoData{}
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

func (s *skillshare) createJsonClass(classData models.ClassData) error {
	value, err := json.MarshalIndent(classData, "", "    ")
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

func (s *skillshare) createJsonVideo(idx int, videoData models.SkillshareVideo, sourceData models.VideoData) error {
	value, err := json.MarshalIndent(sourceData, "", "    ")
	if err != nil {
		return err
	}

	filename := fmt.Sprintf(constants.FilenameVideoData, idx, utils.ToSnakeCase(videoData.Title))
	fileJson := path.Join(s.dir.json, filename)
	err = os.WriteFile(fileJson, value, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *skillshare) loadClassDataCache() (*models.ClassData, error) {
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

	dest := &models.ClassData{}
	err = json.Unmarshal(data, dest)
	if err != nil {
		return nil, err
	}

	if !dest.IsValidVideoId() {
		fmt.Println("Invalid video id, fetch new api again")
		return nil, nil
	}

	return dest, nil
}

func (s *skillshare) loadClassData() (*models.ClassData, error) {
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

	safeTitle := utils.SafeName(getData.Title)
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

	if !getData.IsValidVideoId() {
		return nil, errors.New("invalid video id, please use php session id with premium account")
	}

	return getData, nil
}

func (s *skillshare) loadVideoData(ssClass models.ClassData) (*models.SkillshareClass, error) {
	ss := ssClass.Mapper()

	for idx, val := range ss.Videos[:1] {
		fmt.Printf("%03d. %s\n", idx+1, val.Title)
		video, err := s.fetchVideoApi(val.ID)
		if err != nil {
			return nil, err
		}

		err = s.createJsonVideo(idx+1, val, *video)
		if err != nil {
			return nil, err
		}

		ss.Videos[idx].AddSourceSubtitle(*video)
	}

	return &ss, nil
}
