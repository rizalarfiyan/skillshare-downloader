package models

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/rizalarfiyan/skillshare-downloader/constants"
	"github.com/rizalarfiyan/skillshare-downloader/logger"
)

type Config struct {
	UrlOrId   string
	SessionId string
	Lang      string
	Dir       string
}

type AppConfig struct {
	ID        int
	SessionId string
	Lang      string
	Dir       string
}

func (conf *AppConfig) parseID(config Config) error {
	logger.Debug("Parse ID from config")
	if config.UrlOrId == "" {
		return errors.New("class id or url is required")
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
		if language == "" {
			logger.Debug("Set default language")
			language = constants.DefaultLanguage
		}
		if conf.Lang == "" {
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

func (conf *AppConfig) parseSessionID(config Config) error {
	if config.SessionId == "" {
		return errors.New("session id is required")
	}

	logger.Debug("Set session id from config")
	conf.SessionId = config.SessionId
	return nil
}

func (conf *AppConfig) parseLanguage(config Config) {
	if config.Lang == "" && conf.Lang != "" {
		logger.Debug("Set default language")
		conf.Lang = constants.DefaultLanguage
		return
	}

	logger.Debug("Set language from config")
	conf.Lang = config.Lang
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

func (conf *AppConfig) FromConfig(config Config) error {
	logger.Debug("Do parse class id")
	if err := conf.parseID(config); err != nil {
		return err
	}

	logger.Debug("Do session id")
	if err := conf.parseSessionID(config); err != nil {
		return err
	}

	logger.Debug("Do language")
	conf.parseLanguage(config)

	logger.Debug("Do directory")
	conf.parseDirectory(config)

	return nil
}
