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
	app.Usage = "Remove a reference between obsolete Docker images and their blobs for later removing by garbage collector"
	app.Version = "1.0.0"
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
			Name:     "max-age",
			Usage:    "Keep images whose age is less than this value (in days). Older images will be unlinked with their blobs.",
			EnvVar:   "MAX_AGE",
			Required: true,
		},
		cli.IntFlag{
			Name:     "keep-tag",
			Usage:    "Number of tags to be preserved. If all the tags are outdated (by max-age), this parameter is to make sure that latest tags will not be unlinked.",
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
	maxAge := c.Int("max-age")
	keepTag := c.Int("keep-tag")
	verbose := c.Bool("verbose")

	dkService := dkservice.New(dockerUsername, dockerPassword, host, verbose)
	dkService.SweepOutdatedImages(repo, maxAge, keepTag)

	return nil
}
