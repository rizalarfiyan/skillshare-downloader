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
	"sync"

	"github.com/rizalarfiyan/skillshare-downloader/constants"
	"github.com/rizalarfiyan/skillshare-downloader/logger"
	"github.com/rizalarfiyan/skillshare-downloader/models"
	"github.com/rizalarfiyan/skillshare-downloader/utils"
)

type skillshare struct {
	ctx  context.Context
	conf models.AppConfig
	dir  struct {
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
	logger.Debug("Load the config")
	if err := s.conf.FromConfig(conf); err != nil {
		return err
	}

	logger.Info("Success load config")
	logger.Debug("Initial directory")
	if err := s.initDir(); err != nil {
		return err
	}

	logger.Info("Success create directory")
	logger.Debug("Load class data")
	ssClass, err := s.loadClassData()
	if err != nil {
		return err
	}

	logger.Debug("Load video data")
	_, err = s.workerVideoData(*ssClass)
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
	logger.Debugf("Create directory: %s", s.conf.Dir)
	err := utils.CreateDir(s.conf.Dir)
	if err != nil {
		return err
	}

	logger.Debugf("Read directory: %s", s.conf.Dir)
	dirs, err := utils.ReadDir(s.conf.Dir)
	if err != nil {
		return err
	}

	logger.Debugf("Search downloaded directory: %s", s.conf.Dir)
	dirs = utils.Filter(dirs, func(dir string) bool {
		return strings.HasPrefix(dir, fmt.Sprintf("[%d]", s.conf.ID))
	})

	if len(dirs) != 1 {
		logger.Debug("Skip cache directory")
		return nil
	}

	s.dir.base = path.Join(s.conf.Dir, dirs[0])
	logger.Debugf("Directory found: %s", s.dir.base)
	err = s.loadDir()
	if err != nil {
		return err
	}
	return nil
}

func (s *skillshare) loadDir() error {
	if s.dir.base == "" {
		logger.Debug("Skip load directory")
		return nil
	}

	s.dir.json = path.Join(s.dir.base, "json")
	logger.Debugf("Create directory: %s", s.dir.json)
	err := utils.CreateDir(s.dir.json)
	if err != nil {
		return err
	}

	s.dir.video = path.Join(s.dir.base, "video")
	logger.Debugf("Create directory: %s", s.dir.video)
	err = utils.CreateDir(s.dir.video)
	if err != nil {
		return err
	}

	s.dir.subtitle = path.Join(s.dir.base, "subtitle")
	logger.Debugf("Create directory: %s", s.dir.subtitle)
	err = utils.CreateDir(s.dir.subtitle)
	if err != nil {
		return err
	}

	return nil
}

func (s *skillshare) fetchClassApi() (*models.ClassData, error) {
	client := &http.Client{}
	url := fmt.Sprintf(constants.APIClass, s.conf.ID)
	logger.Debugf("Prepare request API with class id: %d", s.conf.ID)
	req, err := http.NewRequestWithContext(s.ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	logger.Debugf("[%d] Prepare request header", s.conf.ID)
	req.Header = http.Header{
		"Accept":     {"application/vnd.skillshare.class+json;,version=0.8"},
		"User-Agent": {"Skillshare/5.3.0; Android 9.0.1"},
		"Host":       {"api.skillshare.com"},
		"Referer":    {"https://www.skillshare.com/"},
		"cookie":     {s.conf.Cookies},
	}

	logger.Debugf("[%d] Send request to API", s.conf.ID)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	logger.Debugf("[%d] Has status code: %d", s.conf.ID, resp.StatusCode)
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("Skillshare class not found")
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return nil, fmt.Errorf("invalid Skillshare cookies")
	}

	logger.Debugf("[%d] Read response body", s.conf.ID)
	dest := &models.ClassData{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	logger.Debugf("[%d] Parse json response body to struct", s.conf.ID)
	err = json.Unmarshal(body, dest)
	if err != nil {
		return nil, err
	}

	logger.Infof("[%d] Success get class data from api", s.conf.ID)
	return dest, nil
}

func (s *skillshare) fetchVideoApi(videoID int) (*models.VideoData, error) {
	logger.Debugf("[%d] Prepare request API with video id", videoID)
	client := &http.Client{}
	url := fmt.Sprintf(constants.APIVideo, constants.BrightcoveAccountId, videoID)
	req, err := http.NewRequestWithContext(s.ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	logger.Debugf("[%d] Prepare request header", videoID)
	req.Header = http.Header{
		"Accept":     {fmt.Sprintf("application/json;pk=%s", constants.PolicyKey)},
		"User-Agent": {"Mozilla/5.0 (X11; Linux x86_64; rv:52.0) Gecko/20100101 Firefox/52.0"},
		"Origin":     {"https://www.skillshare.com/"},
	}

	logger.Debugf("[%d] Send request to API", videoID)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	logger.Debugf("[%d] Has status code: %d", videoID, resp.StatusCode)
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("Skillshare video id %d not found", videoID)
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return nil, fmt.Errorf("Skillshare video id %d internal server error", videoID)
	}

	dest := &models.VideoData{}
	logger.Debugf("[%d] Read response body", videoID)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	logger.Debugf("[%d] Parse json response body to struct", videoID)
	err = json.Unmarshal(body, dest)
	if err != nil {
		return nil, err
	}

	logger.Infof("[%d] Success get class data from api", videoID)
	return dest, nil
}

func (s *skillshare) createJsonClass(classData models.ClassData) error {
	logger.Debug("Pretty json class data")
	value, err := json.MarshalIndent(classData, "", "    ")
	if err != nil {
		return err
	}

	fileJson := path.Join(s.dir.json, constants.FilenameClassData)
	logger.Debugf("Write json class data to file: %s", fileJson)
	err = os.WriteFile(fileJson, value, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	logger.Infof("Succes create json class id: %s", classData.ID)
	return nil
}

func (s *skillshare) createJsonVideo(idx int, videoData models.SkillshareVideo, sourceData models.VideoData) error {
	logger.Debugf("[%d] Pretty json class data", videoData.ID)
	value, err := json.MarshalIndent(sourceData, "", "    ")
	if err != nil {
		return err
	}

	filename := fmt.Sprintf(constants.FilenameVideoData, idx, utils.ToSnakeCase(videoData.Title))
	fileJson := path.Join(s.dir.json, filename)
	logger.Debugf("[%d] Write json class data to file: %s", videoData.ID, fileJson)
	err = os.WriteFile(fileJson, value, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	logger.Infof("[%d] Succes create json video id", videoData.ID)
	return nil
}

func (s *skillshare) loadClassDataCache() (*models.ClassData, error) {
	if s.dir.base == "" || s.dir.json == "" {
		logger.Info("No cache in directory")
		return nil, nil
	}

	fileJson := path.Join(s.dir.json, constants.FilenameClassData)
	logger.Debugf("Check cache in file: %s", fileJson)
	if !utils.IsExistPath(fileJson) {
		logger.Info("No cache in directory")
		return nil, nil
	}

	logger.Debugf("Load cache from directory: %s", fileJson)
	data, err := os.ReadFile(fileJson)
	if err != nil {
		return nil, err
	}

	dest := &models.ClassData{}
	logger.Debug("Parse json from cache")
	err = json.Unmarshal(data, dest)
	if err != nil {
		return nil, err
	}

	logger.Debug("Check valid video id")
	if !dest.IsValidVideoId() {
		logger.Warning("Invalid video id, fetch new api again")
		return nil, nil
	}

	logger.Info("All video id in cache is valid")
	return dest, nil
}

func (s *skillshare) loadClassData() (*models.ClassData, error) {
	logger.Debug("Do load class data from cache")
	getCache, err := s.loadClassDataCache()
	if err != nil {
		return nil, err
	}

	if getCache != nil {
		logger.Info("Load class data from cache")
		return getCache, nil
	}

	logger.Debug("Do load fetch data to api")
	getData, err := s.fetchClassApi()
	if err != nil {
		return nil, err
	}

	safeTitle := utils.SafeName(getData.Title)
	folderName := fmt.Sprintf(constants.FolderName, s.conf.ID, safeTitle)
	logger.Debugf("Prepare folder name: %s", folderName)
	s.dir.base = path.Join(s.conf.Dir, folderName)
	logger.Debugf("Create directory: %s", s.dir.base)
	err = s.loadDir()
	if err != nil {
		return nil, err
	}

	logger.Debug("Do create json for data class")
	err = s.createJsonClass(*getData)
	if err != nil {
		return nil, err
	}

	logger.Debug("Check valid video id")
	if !getData.IsValidVideoId() {
		return nil, errors.New("invalid video id, please use cookies with premium account")
	}

	logger.Info("All video id is valid")
	return getData, nil
}

type VideoWorker struct {
	Idx           int
	VideoId       int
	Name          string
	Video         *models.VideoData
	OriginalVideo models.SkillshareVideo
	Error         error
}

func (s *skillshare) createWorkerVideo(ss models.SkillshareClass) <-chan VideoWorker {
	chanWorker := make(chan VideoWorker)

	go func() {
		for idx, val := range ss.Videos {
			chanWorker <- VideoWorker{
				Idx:           idx,
				OriginalVideo: val,
				VideoId:       val.ID,
				Name:          fmt.Sprintf("%03d. %s", idx+1, val.Title),
			}
		}

		close(chanWorker)
	}()

	return chanWorker
}

func (s *skillshare) actionWorkerVideo(chanIn <-chan VideoWorker) <-chan VideoWorker {
	chanWorker := make(chan VideoWorker)
	wg := new(sync.WaitGroup)
	wg.Add(s.conf.Worker)

	logger.Debug("Do Loop for videos")
	go func() {
		for workerIdx := 0; workerIdx < s.conf.Worker; workerIdx++ {
			go func(workerIdx int) {
				for val := range chanIn {
					logger.Debugf("[%d] Do run video: %s", val.VideoId, val.Name)
					video, err := s.fetchVideoApi(val.VideoId)
					if err != nil {
						val.Error = err
						chanWorker <- val
						continue
					}

					logger.Debugf("[%d] Do create json", val.VideoId)
					err = s.createJsonVideo(val.Idx, val.OriginalVideo, *video)
					if err != nil {
						val.Error = err
						chanWorker <- val
						continue
					}

					logger.Infof("[%d] Success get video", val.VideoId)
					val.Video = video
					chanWorker <- val
				}
				wg.Done()
			}(workerIdx)
		}
	}()

	go func() {
		wg.Wait()
		close(chanWorker)
	}()

	return chanWorker
}

func (s *skillshare) workerVideoData(ssClass models.ClassData) (*models.SkillshareClass, error) {
	logger.Debug("Mapping response api to new struct")
	ss := ssClass.Mapper()

	chanIn := s.createWorkerVideo(ss)
	chanOut := s.actionWorkerVideo(chanIn)

	countError := 0
	for worker := range chanOut {
		if worker.Error != nil {
			logger.Warningf("Error get video %s", worker.Error.Error())
			countError++
			continue
		}

		logger.Debugf("[%d] Mapping data source to subtitle", worker.VideoId)
		ss.Videos[worker.Idx].AddSourceSubtitle(*worker.Video)
		logger.Infof("[%d] Success fetch video: %03d. %s", worker.VideoId, worker.Idx+1, ss.Videos[worker.Idx].Title)
	}

	return &ss, nil
}
