package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v1"
)

const usage = `dockerx is a simple container runtime implementation.
			   The purpose of this project is to learn how docker works and how to write a docker by ourselves
			   Enjoy it, just for fun.`

func main() {
	app := cli.NewApp()
	app.Name = "dockerx"
	app.Usage = usage

	app.Commands = []cli.Command{
		initCommand,
		runCommand,
		// listCommand,
		// logCommand,
		// execCommand,
		// stopCommand,
		// removeCommand,
		// commitCommand,
		// networkCommand,
	}

	app.Before = func(context *cli.Context) error {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
