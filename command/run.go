package command

import (
	"czlingo/my-docker/cgroups"
	"czlingo/my-docker/cgroups/subsystems"
	"czlingo/my-docker/container"
	"czlingo/my-docker/fs"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func Run(tty bool, comArray []string, res *subsystems.ResourceConfig) {
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		logrus.Errorf("New parent process error")
		return
	}
	// cd /
	parent.Dir = fs.NewWorkspace()
	defer fs.DestroyWorkspace()

	if err := parent.Start(); err != nil {
		logrus.Error(err)
	}
	// use mydocker-cgroup as cgroup name
	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(comArray, writePipe)
	parent.Wait()
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	logrus.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
