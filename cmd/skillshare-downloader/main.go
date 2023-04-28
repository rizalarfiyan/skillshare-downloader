package main

import (
	"context"
	"os"
	"sort"
	"time"

	"github.com/rizalarfiyan/skillshare-downloader/constants"
	"github.com/rizalarfiyan/skillshare-downloader/logger"
	"github.com/rizalarfiyan/skillshare-downloader/models"
	"github.com/rizalarfiyan/skillshare-downloader/services"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func init() {
	logger.Init()
}

func main() {
	log := logger.Get()
	ctx := context.Background()
	defer func() {
		if rec := recover(); rec != nil {
			log.Fatalln("Panic: ", rec)
		}
	}()

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Print the version of skillshare downloader",
	}

	app := &cli.App{
		Name:     "Skillshare Downloader",
		Usage:    "Download the skillshare video with premium account! 🎉",
		Version:  "v1.0.0",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name: "Rizal Arfiyan",
			},
		},
		HelpName:  "Skillshare Downloader",
		UsageText: "skillshare-dl -c <class> -s <session> [args and such]\n",
		ArgsUsage: "[args and such]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "class",
				Aliases:  []string{"c"},
				Usage:    "Identity skillshare class id or skillshare class url",
				Category: "Required:",
			},
			&cli.StringFlag{
				Name:     "session",
				Aliases:  []string{"s"},
				Usage:    "Session id for get content to skillshare, you can get session id from cookie in browser with key (PHPSESSID)",
				Category: "Required:",
			},
			&cli.StringFlag{
				Name:        "language",
				Aliases:     []string{"l"},
				Usage:       "Language subtitle for download the video",
				DefaultText: constants.DefaultLanguage,
				Category:    "Optional:",
			},
			&cli.StringFlag{
				Name:        "directory",
				Aliases:     []string{"d"},
				Usage:       "Directory name for save the video",
				DefaultText: constants.DefaultDir,
				Category:    "Optional:",
			},
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"vvv"},
				Usage:       "Verbose mode to see all logs",
				DefaultText: "false",
				Category:    "Optional:",
			},
		},
		Action: func(cliCtx *cli.Context) error {
			if cliCtx.Bool("verbose") {
				logger.SetLevel(logrus.DebugLevel)
			}

			err := services.NewSkillshare(ctx).Run(models.Config{
				UrlOrId:   cliCtx.String("class"),
				SessionId: cliCtx.String("session"),
				Lang:      cliCtx.String("language"),
				Dir:       cliCtx.String("directory"),
			})
			if err != nil {
				return err
			}

			return nil
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
