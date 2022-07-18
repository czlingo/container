package command

import (
	"czlingo/my-docker/fs"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func Commit(containerName, imageName string) {
	mergedPath := fmt.Sprintf(fs.DefaultMergedPath, containerName)
	imageTar := filepath.Join(fs.DefaultImageRepo, imageName+".tar")
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mergedPath, ".").CombinedOutput(); err != nil {
		logrus.Errorf("Tar folder %s error %v", mergedPath, err)
	}
}
