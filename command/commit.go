package command

import (
	"os/exec"

	"github.com/sirupsen/logrus"
)

func Commit(imageName string) {
	mergedPath := "./root/merged"
	imageTar := "./root/" + imageName + ".tar"
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mergedPath, ".").CombinedOutput(); err != nil {
		logrus.Errorf("Tar folder %s error %v", mergedPath, err)
	}
}
