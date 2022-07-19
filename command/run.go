package command

import (
	"czlingo/my-docker/cgroups"
	"czlingo/my-docker/cgroups/subsystems"
	"czlingo/my-docker/container"
	"czlingo/my-docker/fs"
	"czlingo/my-docker/network"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func Run(tty bool, comArray, volumes, envSlice []string, res *subsystems.ResourceConfig,
	containerName, imageName string, nw string, portmapping []string) {
	containerID := randStringBytes(10)
	if containerName == "" {
		containerName = containerID
	}

	parent, writePipe := container.NewParentProcess(tty, containerName, envSlice)
	if parent == nil {
		logrus.Errorf("New parent process error")
		return
	}

	// cd /
	parent.Dir = fs.NewWorkspace(volumes, containerName, imageName)

	if err := parent.Start(); err != nil {
		logrus.Error(err)
	}

	containerName, err := recordContainerInfo(parent.Process.Pid, comArray, volumes, containerID, containerName)
	if err != nil {
		logrus.Errorf("Record container info error %v", err)
		return
	}

	cgroupManager := cgroups.NewCgroupManager(containerName)
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	if nw != "" {
		// config container network
		network.Init()
		containerInfo := &container.ContainerInfo{
			Id:          containerID,
			Pid:         strconv.Itoa(parent.Process.Pid),
			Name:        containerName,
			PortMapping: portmapping,
		}
		if err := network.Connect(nw, containerInfo); err != nil {
			log.Errorf("Error Connect Network %v", err)
			return
		}
	}

	sendInitCommand(comArray, writePipe)
	if tty {
		// maybe stop container, not remove
		parent.Wait()
		// deleteContainerInfo(containerName)
		fs.DestroyWorkspace(volumes, containerName)
		cgroupManager.Destroy()
	}
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	logrus.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}

func recordContainerInfo(containerPID int, commandArray, volumes []string, containerID, containerName string) (string, error) {
	createTime := time.Now().Format("2006-01-02 15:04:05")
	command := strings.Join(commandArray, "")

	containerInfo := &container.ContainerInfo{
		Id:          containerID,
		Pid:         strconv.Itoa(containerPID),
		Command:     command,
		CreatedTime: createTime,
		Name:        containerName,
		Volumes:     volumes,
		Status:      container.RUNNING,
	}

	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		logrus.Errorf("Record container info error %v", err)
		return "", err
	}

	dirUrl := fmt.Sprintf(fs.DefaultInfoLocation, containerName)
	if err := os.MkdirAll(dirUrl, 0622); err != nil {
		logrus.Errorf("Mkdir error %s error %v", dirUrl, err)
		return "", err
	}
	fileName := dirUrl + "/" + container.ConfigName
	file, err := os.Create(fileName)
	if err != nil {
		logrus.Errorf("Create file %s error %v", fileName, err)
		return "", err
	}
	defer file.Close()

	if _, err := file.Write(jsonBytes); err != nil {
		logrus.Errorf("File write string error %v", err)
		return "", err
	}
	return containerName, nil
}

func randStringBytes(n int) string {
	letterBytes := "1234567890"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// func deleteContainerInfo(containerName string) {
// 	dirURL := fmt.Sprintf(fs.DefaultInfoLocation, containerName)
// 	if err := os.RemoveAll(dirURL); err != nil {
// 		logrus.Errorf("Remove dir %s error %v", dirURL, err)
// 	}
// }
