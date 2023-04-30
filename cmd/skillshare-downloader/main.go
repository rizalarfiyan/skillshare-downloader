package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/rizalarfiyan/skillshare-downloader/constants"
	"github.com/rizalarfiyan/skillshare-downloader/logger"
	"github.com/rizalarfiyan/skillshare-downloader/models"
	"github.com/rizalarfiyan/skillshare-downloader/services"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// these information will be collected when build, by `-ldflags "-X main.appVersion=0.1"`
const notSet string = "not set"

var (
	appVersion = notSet
	buildTime  = notSet
	gitCommit  = notSet
	gitRef     = notSet
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

	cli.VersionPrinter = func(cCtx *cli.Context) {
		var authors []string
		if cCtx.App.Authors != nil {
			for _, author := range cCtx.App.Authors {
				authors = append(authors, author.Name)
			}
		}

		var author string
		if len(authors) > 0 {
			author += "by "
			author += strings.Join(authors, ", ")
		}

		fmt.Printf("%s %s\n\n", cCtx.App.Name, author)
		fmt.Println(strings.Replace(cCtx.App.Version, "\n\t ", "\n", -1))
	}

	cli.AppHelpTemplate = fmt.Sprintf("%s\n\n%s", constants.SplashScreen, cli.AppHelpTemplate)

	app := &cli.App{
		Name:     "Skillshare Downloader",
		Usage:    "Download the skillshare video with premium account! 🎉 \nDO NOT use this project for piracy! \nI'm not responsible for the use of this program, this is only for personal and educational purpose.\nBefore any usage please read the Skillshare Terms of Service. https://skillshare.com/terms",
		Version:  fmt.Sprintf("version = %s \n\t build_time = %s \n\t git_commit = %s \n\t git_ref = %s", appVersion, buildTime, gitCommit, gitRef),
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
