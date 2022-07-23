package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	version = "0.0.1"
)

func main() {
	// referenced from  https://github.com/drone-plugins/drone-gitea-release/blob/master/main.go

	app := cli.NewApp()
	app.Name = "drone-compressed-files plugin"
	app.Usage = "drone-compressed-files plugin"
	app.Action = run
	app.Version = version

	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:   "input",
			Usage:  "file input eg: a.txt, /app/a.txt , 'a/*.js'",
			EnvVar: "PLUGIN_INPUT,ZIP_INPUT",
		},
		cli.StringFlag{
			Name:   "output",
			Usage:  "output zip file path eg: /app/a.zip",
			EnvVar: "PLUGIN_OUTPUT,ZIP_OUTPUT",
		},
		cli.StringSliceFlag{
			Name:   "exclude",
			Usage:  "exclude file",
			EnvVar: "PLUGIN_EXCLUDE,ZIP_EXCLUDE",
		},
	}
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}

}

func run(c *cli.Context) error {

	plugin := Plugin{
		Input:   c.StringSlice("input"),
		Output:  c.String("output"),
		Exclude: c.StringSlice("exclude"),
	}

	logrus.Infof("input: %v", plugin.Input)
	logrus.Infof("output: %v\n", plugin.Output)
	logrus.Infof("exclude: %v\n", plugin.Exclude)

	return plugin.Exec()
}
