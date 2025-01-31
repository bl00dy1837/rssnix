package main

import (
	"errors"
	"os"
	"os/exec"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	syscall.Umask(0)
	LoadConfig()

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "opens the config file with $EDITOR",
				Action: func(cCtx *cli.Context) error {
					editor, ok := os.LookupEnv("EDITOR")
					if len(editor) == 0 || !ok {
						return errors.New("$EDITOR environment variable is not set")
					}
					homePath, err := os.UserHomeDir()
					if err != nil {
						log.Error("Failed to get home path")
						os.Exit(1)
					}
					cmd := exec.Command(editor, homePath+"/.config/rssnix/config.ini")
					cmd.Stdin = os.Stdin
					cmd.Stdout = os.Stdout
					return cmd.Run()
				},
			},
			{
				Name:    "update",
				Aliases: []string{"u"},
				Usage:   "update given feed(s) or all feeds if no argument is given",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() == 0 {
						UpdateAllFeeds()
					}
					for i := 0; i < cCtx.Args().Len(); i++ {
						UpdateFeed(cCtx.Args().Get(i))
					}
					return nil
				},
			},
			{
				Name:    "open",
				Aliases: []string{"o"},
				Usage:   "open given feed's directory or root feeds directory if no argument is given",
				Action: func(cCtx *cli.Context) error {
					var path string
					if cCtx.Args().Len() == 0 {
						path = Config.FeedDirectory
					} else {
						path = Config.FeedDirectory + "/" + cCtx.Args().Get(0)
					}
					cmd := exec.Command(Config.Viewer, path)
					cmd.Stdin = os.Stdin
					cmd.Stdout = os.Stdout
					return cmd.Run()
				},
			},
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add a given feed to config",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() != 2 {
						return errors.New("exactly two arguments are required, first being feed name, second being URL")
					}
					homePath, err := os.UserHomeDir()
					if err != nil {
						log.Error("Failed to get home path")
						os.Exit(1)
					}
					file, err := os.OpenFile(homePath+"/.config/rssnix/config.ini", os.O_APPEND|os.O_WRONLY, 0644)
					if err != nil {
						return err
					}
					defer file.Close()
					_, err = file.WriteString("\n" + cCtx.Args().Get(0) + " = " + cCtx.Args().Get(1))
					return err
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Error(err)
	}
}
