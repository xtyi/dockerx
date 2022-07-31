package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/xtyi/docker/cgroups"
	"github.com/xtyi/docker/cgroups/subsystems"
	"github.com/xtyi/docker/container"
	"gopkg.in/urfave/cli.v1"
)

var runCommand = cli.Command{
	Name: "run",
	Usage: `Create a container with namespace and cgroups limit
			dockerx run -it [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "it",
			Usage: "enable tty",
		},
		cli.BoolFlag{
			Name:  "d",
			Usage: "deteach container",
		},
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name:  "cpushare",
			Usage: "cpushare limit",
		},
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "container name",
		},
		cli.StringFlag{
			Name:  "v",
			Usage: "volume",
		},
		cli.StringSliceFlag{
			Name:  "e",
			Usage: "set environment",
		},
		cli.StringFlag{
			Name:  "net",
			Usage: "container network",
		},
		cli.StringSliceFlag{
			Name:  "p",
			Usage: "port mapping",
		},
	},
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container command")
		}
		var cmds []string
		for _, arg := range context.Args() {
			cmds = append(cmds, arg)
		}

		tty := context.Bool("it")
		conf := &subsystems.ResourceConfig{
			MemoryLimit: context.String("m"),
			CPUSet:      context.String("cpuset"),
			CPUShare:    context.String("cpushare"),
		}
		Run(cmds, tty, conf)
		return nil
	},
}

func Run(cmds []string, tty bool, conf *subsystems.ResourceConfig) {
	// 创建 init 进程结构体
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		logrus.Errorf("new parent process error")
		return
	}
	// 启动 init 进程
	if err := parent.Start(); err != nil {
		logrus.Error(err)
	}
	cgroupManager := cgroups.NewCgroupManager("dockerx-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(conf)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(cmds, writePipe)
	parent.Wait()
	mntURL := "/root/mnt/"
	rootURL := "/root/"
	container.DeleteWorkSpace(rootURL, mntURL)
	os.Exit(-1)
}

func sendInitCommand(cmds []string, writePipe *os.File) {
	command := strings.Join(cmds, " ")
	logrus.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
