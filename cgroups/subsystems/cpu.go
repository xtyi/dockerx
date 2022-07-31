package subsystems

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type CPUSubSystem struct {
}

func (s *CPUSubSystem) Set(cgroupPath string, config *ResourceConfig) error {
	subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, true)
	if err != nil {
		return err
	}
	if config.CPUShare == "" {
		return nil
	}
	err = ioutil.WriteFile(path.Join(subsysCgroupPath, "cpu.shares"), []byte(config.CPUShare), 0644)
	if err != nil {
		return fmt.Errorf("set cgroup cpu share fail %w", err)
	}
	return nil
}

func (s *CPUSubSystem) Remove(cgroupPath string) error {
	subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false)
	if err != nil {
		return err
	}
	return os.RemoveAll(subsysCgroupPath)
}

func (s *CPUSubSystem) Apply(cgroupPath string, pid int) error {
	subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false)
	if err != nil {
		return fmt.Errorf("get cgroup %s error: %w", cgroupPath, err)
	}
	err = ioutil.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		return fmt.Errorf("set cgroup proc fail %w", err)
	}
	return nil
}

func (s *CPUSubSystem) Name() string {
	return "cpu"
}
