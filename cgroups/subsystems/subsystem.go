package subsystems

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

type ResourceConfig struct {
	MemoryLimit string
	CPUShare    string
	CPUSet      string
}

type SubSystem interface {
	Name() string
	Set(path string, config *ResourceConfig) error
	Apply(path string, pid int) error
	Remove(path string) error
}

var SubSystemList = []SubSystem{
	&CPUSetSubSystem{},
	&MemorySubSystem{},
	&CPUSubSystem{},
}

func FindCgroupMountPoint(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				return fields[4]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return ""
	}
	return ""
}

func GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	cgroupRoot := FindCgroupMountPoint(subsystem)
	_, err := os.Stat(path.Join(cgroupRoot, cgroupPath))
	if err != nil {
		if os.IsNotExist(err) && autoCreate {
			err := os.Mkdir(path.Join(cgroupRoot, cgroupPath), 0755)
			if err != nil {
				return "", fmt.Errorf("error create cgroup %w", err)
			}
		} else {
			return "", fmt.Errorf("cgroup path error %w", err)
		}
	}
	return path.Join(cgroupRoot, cgroupPath), nil
}
