package main

import (
	"fmt"
	"os"

	"github.com/dominikus1993/kup50-tfs/internal/diff"
	"github.com/dominikus1993/kup50-tfs/internal/dir"
	"github.com/dominikus1993/kup50-tfs/internal/git"
	"github.com/dominikus1993/kup50-tfs/internal/zip"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "organization",
				Usage:    "The URL of your Azure DevOps organization (e.g., https://dev.azure.com/myorg)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "pat",
				Usage:    "Your Azure DevOps Personal Access Token",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "project",
				Usage:    "The name of the project within your organization",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "author",
				Usage:    `(optional) - The name of the author whose changes you want to download. If not provided, defaults to "Dominik Kotecki" me xD`,
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
			differ := diff.NewBlobDiffer()
			client, err := git.NewAzureDevopsClient(c.Context, organization, pat, project, differ)
			if err != nil {
				return err
			}
			changes := client.GetChanges(c.Context, author)
			client.DowloadAndSaveChanges(c.Context, changes)
			err = zip.ZipDir("kup", "kup.zip")
			if err != nil {
				return err
			}
			dir.Rm("kup")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
