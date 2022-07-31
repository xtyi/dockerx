package main

import (
	"github.com/sirupsen/logrus"
	"github.com/xtyi/docker/container"
	"gopkg.in/urfave/cli.v1"
)

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",
	Action: func(context *cli.Context) error {
		logrus.Infof("init come on")
		err := container.RunContainerInitProcess()
		return err
	},
}
