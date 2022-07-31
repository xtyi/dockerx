package subsystems

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type CPUSetSubSystem struct {
}

func (s *CPUSetSubSystem) Set(cgroupPath string, config *ResourceConfig) error {
	subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, true)
	if err != nil {
		return err
	}
	if config.CPUSet == "" {
		return nil
	}
	err = ioutil.WriteFile(path.Join(subsysCgroupPath, "cpuset.cpus"), []byte(config.CPUSet), 0644)
	if err != nil {
		return fmt.Errorf("set cgroup cpuset fail %w", err)
	}
	return nil
}

func (s *CPUSetSubSystem) Remove(cgroupPath string) error {
	subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false)
	if err != nil {
		return err
	}
	return os.RemoveAll(subsysCgroupPath)
}

func (s *CPUSetSubSystem) Apply(cgroupPath string, pid int) error {
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

func (s *CPUSetSubSystem) Name() string {
	return "cpuset"
}
