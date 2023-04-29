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
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/cheggaaa/pb/v3"
	"github.com/melbahja/got"
	"github.com/rizalarfiyan/skillshare-downloader/constants"
	"github.com/rizalarfiyan/skillshare-downloader/logger"
	"github.com/rizalarfiyan/skillshare-downloader/models"
	"github.com/rizalarfiyan/skillshare-downloader/utils"
)

type skillshare struct {
	ctx  context.Context
	conf models.AppConfig
	spin *spinner.Spinner
	dir  struct {
		base  string
		json  string
		video string
	}
}

func NewSkillshare(ctx context.Context) Skillshare {
	spin := spinner.New(spinner.CharSets[78], 100*time.Millisecond, func(s *spinner.Spinner) {
		s.Color("green")
	})

	return &skillshare{
		ctx:  ctx,
		spin: spin,
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
	ssData, err := s.workerVideoData(*ssClass)
	if err != nil {
		return err
	}

	err = s.workerDownloadVideo(*ssData)
	if err != nil {
		return err
	}

	err = s.workerDownloadSubtitle(*ssData)
	if err != nil {
		return err
	}

	return nil
}

// ` + "`" + `
func (s *skillshare) splash() {
	text := `
    +=----------------------------------------------------------------------------------------------------------=+
    
    .d88b. 8    w 8 8      8                        888b.                        8               8            
    YPwww. 8.dP w 8 8 d88b 8d8b. .d88 8d8b .d88b    8   8 .d8b. Yb  db  dP 8d8b. 8 .d8b. .d88 .d88 .d88b 8d8b 
        d8 88b  8 8 8 ` + "`" + `Yb. 8P Y8 8  8 8P   8.dP'    8   8 8' .8  YbdPYbdP  8P Y8 8 8' .8 8  8 8  8 8.dP' 8P   
    ` + "`" + `Y88P' 8 Yb 8 8 8 Y88P 8   8 ` + "`" + `Y88 8    ` + "`" + `Y88P    888P' ` + "`" + `Y8P'   YP  YP   8   8 8 ` + "`" + `Y8P' ` + "`" + `Y88 ` + "`" + `Y88 ` + "`" + `Y88P 8

    +=-------------------------------------------- By Rizal Arfiyan --------------------------------------------=+
    `

	fmt.Printf("\n%s\n\n", text)
	panic("test")
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

	logger.Debugf("[%d] Success get class data from api", s.conf.ID)
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

	logger.Debugf("[%d] Success get video data from api", videoID)
	return dest, nil
}

func (s *skillshare) fetchSubtitle(sub models.SubtitleWorker) ([]byte, error) {
	client := &http.Client{}
	logger.Debugf("Prepare request subtitle video id: %d (%s)", sub.VideoId, sub.Label)
	req, err := http.NewRequestWithContext(s.ctx, "GET", sub.Src, nil)
	if err != nil {
		return nil, err
	}

	logger.Debugf("[%d](%s) Prepare request header", sub.VideoId, sub.Label)
	req.Header = http.Header{
		"Accept":     {"application/vnd.skillshare.class+json;,version=0.8"},
		"User-Agent": {"Skillshare/5.3.0; Android 9.0.1"},
		"Host":       {"skillshare.com"},
		"Referer":    {"https://www.skillshare.com/"},
	}

	logger.Debugf("[%d](%s) Send request to API", sub.VideoId, sub.Label)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	logger.Debugf("[%d](%s) Has status code: %d", sub.VideoId, sub.Label, resp.StatusCode)
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("subtitle not found")
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return nil, err
	}

	logger.Debugf("[%d](%s) Read response body", sub.VideoId, sub.Label)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	logger.Debugf("[%d](%s) Success get subtitle data from api", sub.VideoId, sub.Label)
	return body, nil
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

	logger.Debugf("Succes create json class id: %d", classData.ID)
	return nil
}

func (s *skillshare) createJsonVideo(idx int, videoData models.SkillshareVideo, sourceData models.VideoData) error {
	logger.Debugf("[%d] Pretty json class data", videoData.ID)
	value, err := json.MarshalIndent(sourceData, "", "    ")
	if err != nil {
		return err
	}

	filename := fmt.Sprintf(constants.FilenameVideoData, idx+1, utils.ToSnakeCase(videoData.Title))
	fileJson := path.Join(s.dir.json, filename)
	logger.Debugf("[%d] Write json class data to file: %s", videoData.ID, fileJson)
	err = os.WriteFile(fileJson, value, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	logger.Debugf("[%d] Succes create json video id", videoData.ID)
	return nil
}

func (s *skillshare) createSubtitle(sub models.SubtitleWorker, data []byte) error {
	extension := utils.MatchExtenstion(sub.Src, ".vtt")
	filename := fmt.Sprintf(constants.FilenameSubtitle, sub.Idx+1, utils.ToSnakeCase(sub.Title), extension)
	fileSubtitle := path.Join(s.dir.video, filename)
	logger.Debugf("[%d](%s) Write json class data to file: %s", sub.VideoId, sub.Label, fileSubtitle)
	err := os.WriteFile(fileSubtitle, data, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	logger.Debugf("[%d](%s) Succes create subtitle", sub.VideoId, sub.Label)
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

	if !s.conf.IsVerbose {
		s.spin.Suffix = fmt.Sprintf(" Fetching skillshare class data with id %d\n", s.conf.ID)
		s.spin.Start()
	}

	logger.Debug("Do load fetch data to api")
	getData, err := s.fetchClassApi()
	if err != nil {
		return nil, err
	}

	if !s.conf.IsVerbose {
		s.spin.Suffix = " Create skillshare json data\n"
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

	if !s.conf.IsVerbose {
		s.spin.Suffix = " Fetching skillshare done\n"
		s.spin.Stop()
	}

	logger.Info("Skillshare class data is ready")
	logger.Debug("Check valid video id")
	if !getData.IsValidVideoId() {
		return nil, errors.New("invalid video id, please use cookies with premium account")
	}

	logger.Info("All video id is valid")
	return getData, nil
}

func (s *skillshare) createWorkerVideo(ss models.SkillshareClass) <-chan models.VideoWorker {
	chanWorker := make(chan models.VideoWorker)

	go func() {
		for idx, val := range ss.Videos {
			chanWorker <- models.VideoWorker{
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

func (s *skillshare) actionWorkerVideo(chanIn <-chan models.VideoWorker) <-chan models.VideoWorker {
	chanWorker := make(chan models.VideoWorker)
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

					logger.Debugf("[%d] Success get video", val.VideoId)
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

	if !s.conf.IsVerbose {
		s.spin.Suffix = fmt.Sprintf(" \x1b[36m[%d/%d]\x1b[0m Fetching skillshare video data with id\n", 0, len(ss.Videos))
		s.spin.Start()
	}

	chanIn := s.createWorkerVideo(ss)
	chanOut := s.actionWorkerVideo(chanIn)

	countError := 0
	countSuccess := 0
	for worker := range chanOut {
		if worker.Error != nil {
			logger.Warningf("Error get video %s", worker.Error.Error())
			countError++
			continue
		}

		countSuccess++
		if !s.conf.IsVerbose {
			s.spin.Suffix = fmt.Sprintf(" \x1b[36m[%d/%d]\x1b[0m Fetching skillshare video data with id\n", countSuccess, len(ss.Videos))
		}

		logger.Debugf("[%d] Mapping data source to subtitle", worker.VideoId)
		ss.Videos[worker.Idx].AddSourceSubtitle(*worker.Video)
		logger.Debugf("[%d] Success fetch video: %03d. %s", worker.VideoId, worker.Idx+1, ss.Videos[worker.Idx].Title)
	}

	if !s.conf.IsVerbose {
		s.spin.Suffix = " Fetching skillshare video done\n"
		s.spin.Stop()
	}

	logger.Info("All video data is ready")

	return &ss, nil
}

func (s *skillshare) cleanVideoDir() error {
	files, err := utils.SearchFiles(s.dir.video, "*.part*")
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := os.Remove(file); err != nil {
			return err
		}
	}

	return nil
}

func (s *skillshare) workerDownloadVideo(ssData models.SkillshareClass) error {
	defer func() {
		if rec := recover(); rec != nil {
			logger.Debug("Do clean video dir")
			err := s.cleanVideoDir()
			if err != nil {
				panic(err)
			}
		}
	}()

	for idx, val := range ssData.Videos {
		title := utils.SafeName(val.Title)
		if len(val.Sources) < 1 {
			logger.Warningf("[%d] Video %s has no source", val.ID, title)
			logger.Infof("[%d] Skipping download", val.ID)
			continue
		}

		logger.Debugf("[%d] Preapare download video", val.ID)
		source := val.Sources[0]

		extension := utils.MatchExtenstion(source.Src, fmt.Sprintf(".%s", strings.ToLower(source.Container)))
		fileName := fmt.Sprintf(constants.FilenameVideo, idx+1, utils.ToSnakeCase(title), extension)
		filePath := filepath.Join(s.dir.video, fileName)

		logger.Infof("\x1b[36m\x1b[36m[%d/%d]\x1b[0m\x1b[0m %s", idx+1, len(ssData.Videos), val.Title)
		var bar *pb.ProgressBar
		dl := got.New()
		dl.ProgressFunc = func(download *got.Download) {
			download.Concurrency = uint(runtime.NumCPU())
			if bar == nil {
				bar = pb.ProgressBarTemplate(constants.ProgressBarTemplate).Start64(int64(download.TotalSize()))
				bar.Set(pb.Bytes, true)
			}
			bar.SetCurrent(int64(download.Size()))
		}

		logger.Debugf("[%d] Do download video: %s", val.ID, val.Title)
		err := dl.Download(source.Src, filePath)
		if err != nil {
			return err
		}

		bar.SetCurrent(bar.Total())
		bar.Finish()
	}

	logger.Info("Download video done")

	return nil
}

func (s *skillshare) createWorkerSubtitle(ss models.SkillshareClass) <-chan models.SubtitleWorker {
	chanWorker := make(chan models.SubtitleWorker)

	lang := s.checkLanguage(ss)
	s.conf.Lang = lang.Lang

	go func() {
		for idx, val := range ss.Videos {
			for _, sub := range val.Subtitles {
				if !strings.EqualFold(lang.Lang, sub.Lang) {
					continue
				}

				chanWorker <- models.SubtitleWorker{
					SkillshareVideoSubtitle: sub,
					Title:                   val.Title,
					Idx:                     idx,
					VideoId:                 val.ID,
				}
			}
		}

		close(chanWorker)
	}()

	return chanWorker
}

func (s *skillshare) actionWorkerSubtitle(chanIn <-chan models.SubtitleWorker) <-chan models.SubtitleWorker {
	chanWorker := make(chan models.SubtitleWorker)
	wg := new(sync.WaitGroup)
	wg.Add(s.conf.Worker)

	logger.Debug("Do Loop for subtitles")
	go func() {
		for workerIdx := 0; workerIdx < s.conf.Worker; workerIdx++ {
			go func(workerIdx int) {
				for val := range chanIn {
					logger.Debugf("[%d](%s) Do run download sutitle", val.VideoId, val.Label)
					data, err := s.fetchSubtitle(val)
					if err != nil {
						val.Error = err
						chanWorker <- val
						continue
					}

					logger.Debugf("[%d](%s) Do create subtitle", val.VideoId, val.Label)
					err = s.createSubtitle(val, data)
					if err != nil {
						val.Error = err
						chanWorker <- val
						continue
					}

					logger.Debugf("[%d](%s) Success download subtitle", val.VideoId, val.Label)
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

func (s *skillshare) checkLanguage(ssData models.SkillshareClass) models.Checklang {
	sub := models.Checklang{
		IsFalid: true,
		Lang:    s.conf.Lang,
	}

	if strings.EqualFold(s.conf.Lang, constants.DefaultLanguage) {
		return sub
	}

	for idx, val := range ssData.Videos {
		if idx == 0 {
			mapKey := make(map[string]string)
			for _, item := range val.Subtitles {
				key := strings.ToLower(strings.Split(item.Lang, "-")[0])
				mapKey[key] = item.Lang
			}

			getKey := strings.ToLower(strings.Split(sub.Lang, "-")[0])
			if subOri, ok := mapKey[getKey]; ok {
				sub.Lang = strings.ToLower(subOri)
			} else {
				logger.Infof("[%d] language %s not found, change to default lang", val.ID, sub.Lang)
				sub.Lang = strings.ToLower(constants.DefaultLanguage)
			}

			continue
		}

		mapKey := make(map[string]bool)
		for _, item := range val.Subtitles {
			key := strings.ToLower(item.Lang)
			mapKey[key] = true
		}

		if _, ok := mapKey[sub.Lang]; !ok {
			sub.IsFalid = false
			break
		}
	}

	return sub
}

func (s *skillshare) workerDownloadSubtitle(ss models.SkillshareClass) error {
	if !s.conf.IsVerbose {
		s.spin.Suffix = fmt.Sprintf(" \x1b[36m[%d/%d]\x1b[0m Fetching skillshare video data with id\n", 0, len(ss.Videos))
		s.spin.Start()
	}

	chanIn := s.createWorkerSubtitle(ss)
	chanOut := s.actionWorkerSubtitle(chanIn)

	countError := 0
	countSuccess := 0
	for worker := range chanOut {
		if worker.Error != nil {
			logger.Warningf("Error get subtitle %s", worker.Error.Error())
			countError++
			continue
		}

		countSuccess++
		if !s.conf.IsVerbose {
			s.spin.Suffix = fmt.Sprintf(" \x1b[36m[%d/%d]\x1b[0m Download skillshare subtitle data with language %s\n", countSuccess, len(ss.Videos), worker.Label)
		}
	}

	if !s.conf.IsVerbose {
		s.spin.Suffix = " Download skillshare subtitle done\n"
		s.spin.Stop()
	}

	logger.Info("Download subtitle done")

	return nil
}
