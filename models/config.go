package models

import (
	"errors"
	"regexp"

	"github.com/rizalarfiyan/skillshare-downloader/constants"
)

type Config struct {
	UrlOrId   string
	SessionId string
	Lang      string
	Dir       string
}

type AppConfig struct {
	ID        string
	SessionId string
	Lang      string
	Dir       string
}

func (conf *AppConfig) parseID(config Config) error {
	if config.UrlOrId == "" {
		return errors.New("class id or url is required")
	}

	isClassId, err := regexp.MatchString(constants.RegexSkillshareClassId, config.UrlOrId)
	if err != nil {
		return err
	}

	if isClassId {
		conf.ID = config.UrlOrId
		return nil
	}

	regex := regexp.MustCompile(constants.RegexSkillshareClassUrl)
	match := regex.FindStringSubmatch(config.UrlOrId)
	if len(match) > 0 {
		language := match[1]
		if language == "" {
			language = constants.DefaultLanguage
		}
		if conf.Lang == "" {
			conf.Lang = language
		}
		conf.ID = match[3]
		return nil
	}

	return errors.New("invalid class id or url")
}

func (conf *AppConfig) parseSessionID(config Config) error {
	if config.SessionId == "" {
		return errors.New("session id is required")
	}

	conf.SessionId = config.SessionId
	return nil
}

func (conf *AppConfig) parseLanguage(config Config) {
	if config.Lang == "" && conf.Lang != "" {
		conf.Lang = constants.DefaultLanguage
		return
	}

	conf.Lang = config.Lang
}

func (conf *AppConfig) parseDirectory(config Config) {
	if config.Dir == "" {
		conf.Dir = constants.DefaultDir
		return
	}

	conf.Dir = config.Dir
}

func (conf *AppConfig) FromConfig(config Config) error {
	if err := conf.parseID(config); err != nil {
		return err
	}

	if err := conf.parseSessionID(config); err != nil {
		return err
	}

	conf.parseLanguage(config)
	conf.parseDirectory(config)

	return nil
}
