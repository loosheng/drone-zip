package main

import (
	"fmt"
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
			EnvVar: "ZIP_FILES",
		},
		cli.StringFlag{
			Name:   "output",
			Usage:  "output zip file path eg: /app/a.zip",
			EnvVar: "OUTPUT_PATH",
		},
	}
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}

}

func run(c *cli.Context) error {
	fmt.Printf("input: %v",c.StringSlice("input"))
	fmt.Println(c.String("output"))
	plugin := Plugin{
		Input: c.StringSlice("input"),
		Output: c.String("output"),
	}

	return plugin.Exec()
}
