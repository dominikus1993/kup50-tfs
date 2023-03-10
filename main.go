package main

import (
	"fmt"
	"os"

	"github.com/dominikus1993/kup50-tfs/internal/git"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "organization",
				Usage:    "organization url",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "pat",
				Usage:    "pat token",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "project",
				Usage:    "project",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "author",
				Required: false,
				Value:    "Dominik Kotecki",
			},
		},
		Action: func(c *cli.Context) error {
			pat := c.String("pat")
			organization := c.String("organization")
			project := c.String("project")
			author := c.String("author")
			fmt.Println(pat)
			client, err := git.NewAzureDevopsClient(c.Context, organization, pat, project)
			if err != nil {
				return err
			}
			log.WithField("author", author).WithField("project", project).Infoln("Start parsing")
			changes := client.GetChanges(c.Context, author)
			client.DowloadAndSaveChanges(c.Context, changes)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
