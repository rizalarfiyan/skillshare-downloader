package models

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/rizalarfiyan/skillshare-downloader/constants"
	"github.com/rizalarfiyan/skillshare-downloader/logger"
	"github.com/rizalarfiyan/skillshare-downloader/utils"
)

type Config struct {
	UrlOrId    string
	Cookies    string
	CookieFile string
	Lang       string
	Dir        string
	Worker     int
	IsVerbose  bool
}

type AppConfig struct {
	ID        int
	Cookies   string
	Lang      string
	Dir       string
	Worker    int
	IsVerbose bool
}

func (conf *AppConfig) parseID(config Config) error {
	logger.Debug("Parse ID from config")
	if config.UrlOrId == "" {
		return fmt.Errorf("class id or url is required")
	}

	logger.Debug("Checking url or id skillshare")
	isClassId, err := regexp.MatchString(constants.RegexSkillshareClassId, config.UrlOrId)
	if err != nil {
		return err
	}

	if isClassId {
		logger.Debug("Parse class id string to number")
		classId, err := strconv.Atoi(config.UrlOrId)
		if err != nil {
			return err
		}
		logger.Debug("Detected as skillshare class id")
		conf.ID = classId
		return nil
	}

	logger.Debug("Detected as skillshare class url")
	regex := regexp.MustCompile(constants.RegexSkillshareClassUrl)
	match := regex.FindStringSubmatch(config.UrlOrId)
	if len(match) > 0 {
		language := match[1]
		if language != "" {
			logger.Debug("Set language from url")
			conf.Lang = language
		}
		logger.Debug("Parse class id string to number")
		classId, err := strconv.Atoi(match[3])
		if err != nil {
			return err
		}
		conf.ID = classId
		return nil
	}

	return errors.New("invalid class id or url")
}

func (conf *AppConfig) parseCookies(config Config) error {
	if config.Cookies == "" && config.CookieFile == "" {
		return errors.New("cookies or cookie-file is required")
	}

	if config.Cookies != "" {
		logger.Debug("Set cookies with string cookies")
		conf.Cookies = config.Cookies
		logger.Info("Loaded raw text cookies")
		return nil
	}

	extension := filepath.Ext(config.CookieFile)
	switch extension {
	case ".txt":
		logger.Debug("Set cookies with file txt cookies")
		cookie, err := utils.GetCookieTxt(config.CookieFile)
		if err != nil {
			return err
		}
		conf.Cookies = cookie
		logger.Info("Loaded file txt cookies")
		return nil
	default:
		return errors.New("invalid cookie file extension")
	}
}

func (conf *AppConfig) parseLanguage(config Config) {
	if config.Lang == "" && conf.Lang == "" {
		logger.Debug("Set default language")
		conf.Lang = constants.DefaultLanguage
		return
	}

	if config.Lang != "" && conf.Lang == "" {
		logger.Debug("Set language from config")
		conf.Lang = config.Lang
	}
}

func (conf *AppConfig) parseDirectory(config Config) {
	if config.Dir == "" {
		logger.Debug("Set default directory")
		conf.Dir = constants.DefaultDir
		return
	}

	logger.Debug("Set directory from config")
	conf.Dir = config.Dir
}

func (conf *AppConfig) parseWorker(config Config) {
	if config.Worker == 0 {
		logger.Debug("Set default worker")
		conf.Worker = constants.DefaultWorker
		return
	}

	if config.Worker > constants.MaxWorker {
		logger.Warningf("Worker large than %s", constants.MaxWorker)
		logger.Info("Set default worker")
		conf.Worker = constants.DefaultWorker
		return
	}

	logger.Debug("Set worker from config")
	conf.Worker = config.Worker
}

func (conf *AppConfig) FromConfig(config Config) error {
	logger.Debug("Do parse class id")
	if err := conf.parseID(config); err != nil {
		return err
	}

	logger.Debug("Do parse cookies")
	if err := conf.parseCookies(config); err != nil {
		return err
	}

	logger.Debug("Clean Cookies")
	conf.Cookies = utils.CleanCookies(conf.Cookies)

	logger.Debug("Do language")
	conf.parseLanguage(config)

	logger.Debug("Do directory")
	conf.parseDirectory(config)

	logger.Debug("Do worker")
	conf.parseWorker(config)

	conf.IsVerbose = config.IsVerbose

	return nil
}
