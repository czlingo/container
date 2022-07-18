package command

import (
	"czlingo/my-docker/container"
	"czlingo/my-docker/fs"

	log "github.com/sirupsen/logrus"
)

func RemoveContainer(containerName string) {
	containerInfo, err := getContainerInfoByName(containerName)
	if err != nil {
		log.Errorf("Get container %s info error %v", containerName, err)
		return
	}
	if containerInfo.Status != container.STOP {
		log.Errorf("Couldn't remove running container")
		return
	}
	fs.DestroyWorkspace(containerInfo.Volumes, containerName)
}
