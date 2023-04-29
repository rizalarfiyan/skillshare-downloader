package main

import (
	"context"
	"fmt"
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
		UsageText: "skillshare-dl --class <class> --cookie-file <cookie-path> [args and such]\n",
		ArgsUsage: "[args and such]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "class",
				Aliases:  []string{"c"},
				Usage:    "Identity skillshare class id or skillshare class url",
				Category: "Class:",
			},
			&cli.StringFlag{
				Name:     "cookies",
				Aliases:  []string{"co"},
				Usage:    "String cookies for get content to skillshare",
				Category: "Required Cookies:",
			},
			&cli.StringFlag{
				Name:     "cookie-file",
				Aliases:  []string{"cf"},
				Usage:    "Cookie File (.txt) for get content to skillshare",
				Category: "Required Cookies:",
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
			&cli.IntFlag{
				Name:        "worker",
				Aliases:     []string{"w"},
				Usage:       "Worker for concurrent connection to download the video",
				DefaultText: fmt.Sprint(constants.DefaultWorker),
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
			isVerbose := cliCtx.Bool("verbose")
			if isVerbose {
				logger.SetLevel(logrus.DebugLevel)
			}

			err := services.NewSkillshare(ctx).Run(models.Config{
				UrlOrId:    cliCtx.String("class"),
				Cookies:    cliCtx.String("cookies"),
				CookieFile: cliCtx.String("cookie-file"),
				Lang:       cliCtx.String("language"),
				Dir:        cliCtx.String("directory"),
				Worker:     cliCtx.Int("worker"),
				IsVerbose:  isVerbose,
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
