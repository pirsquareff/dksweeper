package main

import (
	"log"
	"os"
	"time"

	"github.com/pirsquareff/dksweeper/src/dkservice"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Docker Registry Sweeper"
	app.Usage = "Remove a reference between obsolete Docker images and their blobs for later deleting by garbage collection"
	app.Version = "1.0.1"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Parinthorn Saithong",
			Email: "parinthorn.sa@gmail.com",
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "username, u",
			Usage:  "Username to access Docker Registry",
			EnvVar: "DOCKER_USERNAME",
		},
		cli.StringFlag{
			Name:   "password, p",
			Usage:  "Password to access Docker Registry",
			EnvVar: "DOCKER_PASSWORD",
		},
		cli.StringFlag{
			Name:     "host",
			Usage:    "Docker Registry host with a protocol (http, https)",
			EnvVar:   "REGISTRY_HOST",
			Required: true,
		},
		cli.StringFlag{
			Name:     "repo, r",
			Usage:    "Repository for cleanup",
			EnvVar:   "REPOSITORY",
			Required: true,
		},
		cli.IntFlag{
			Name:     "older-than",
			Usage:    "Delete images older than this value (in days)",
			EnvVar:   "OLDER_THAN",
			Required: true,
		},
		cli.IntFlag{
			Name:     "keep-tag",
			Usage:    "Number of tags to be preserved. This parameter is to make sure that latest tags will not be unlinked.",
			EnvVar:   "KEEP_TAG",
			Required: true,
		},
		cli.BoolFlag{
			Name:     "verbose",
			Usage:    "Print info during processing",
			EnvVar:   "VERBOSE",
			Required: false,
		},
	}

	app.Action = run

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	dockerUsername := c.String("username")
	dockerPassword := c.String("password")
	host := c.String("host")
	repo := c.String("repo")
	olderThan := c.Int("older-than")
	keepTag := c.Int("keep-tag")
	verbose := c.Bool("verbose")

	dkService := dkservice.New(dockerUsername, dockerPassword, host, verbose)
	dkService.SweepOutdatedImages(repo, olderThan, keepTag)

	return nil
}
