package fs

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	DefaultImageRepo    = "./mycontainer/image"
	DefaultRootPath     = "./mycontainer/container"
	DefaultLowerPath    = "./mycontainer/container/%s/lower"
	DefaultMergedPath   = "./mycontainer/container/%s/merged"
	DefaultDiffPath     = "./mycontainer/container/%s/diff"
	DefaultWorkPath     = "./mycontainer/container/%s/work"
	DefaultInfoLocation = "./mycontainer/container/%s/"
)

func NewWorkspace(volumes []string, containerName, imageName string) string {
	lowerDir := createLowerDir(containerName, imageName)
	diffDir := createDiffDir(containerName)
	workDir := createWorkDir(containerName)
	mergedDir := mountMergeDir(containerName, lowerDir, diffDir, workDir)
	mountVolumes(mergedDir, volumes)
	return mergedDir
}

func createLowerDir(containerName, imageName string) string {
	lowerDir := fmt.Sprintf(DefaultLowerPath, containerName)
	if err := os.MkdirAll(lowerDir, 0777); err != nil {
		logrus.Errorf("Fail to create lower dir %s, error %v", lowerDir, err)
	}
	imageURI := filepath.Join(DefaultImageRepo, imageName)
	cmd := exec.Command("tar", "-xvf", imageURI+".tar", "-C", lowerDir)
	if err := cmd.Run(); err != nil {
		logrus.Errorf("Fail to tar busybox")
	}
	return lowerDir
}

func createDiffDir(containerName string) string {
	diffPath := fmt.Sprintf(DefaultDiffPath, containerName)
	if err := os.MkdirAll(diffPath, 0777); err != nil {
		logrus.Errorf("Fail to create diff dir. error %v", err)
	}
	return diffPath
}

func createWorkDir(containerName string) string {
	workPath := fmt.Sprintf(DefaultWorkPath, containerName)
	if err := os.MkdirAll(workPath, 0777); err != nil {
		logrus.Errorf("Fail to create diff dir. error %v", err)
	}
	return workPath
}

func mountMergeDir(containerName, lowPath, diffPath, workPath string) string {
	mountPath := fmt.Sprintf(DefaultMergedPath, containerName)
	if err := os.MkdirAll(mountPath, 0777); err != nil {
		logrus.Errorf("Fail to create mount dir. error %v", err)
	}

	option := "lowerdir=" + lowPath + ",upperdir=" + diffPath + ",workdir=" + workPath
	cmd := exec.Command("mount", "-t", "overlay", "-o", option, "overlay", mountPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Errorf("Fail to mount, %v", err)
	}
	return mountPath
}

func mountVolumes(merged string, volumes []string) {
	for _, volume := range volumes {
		volMapping := strings.Split(volume, ":")
		if len(volMapping) == 2 && volMapping[0] != "" && volMapping[1] != "" {
			hostPath := volMapping[0]
			mountPoint := filepath.Join(merged, volMapping[1])
			if err := os.Mkdir(mountPoint, 0777); err != nil {
				logrus.Errorf("Fail to create volume mount point %s, error %v", mountPoint, err)
			}
			cmd := exec.Command("mount", "--bind", hostPath, mountPoint)
			if err := cmd.Run(); err != nil {
				logrus.Errorf("Fail to mount volume %s. error %v", volume, err)
			}
		}
	}
}

func DestroyWorkspace(volumes []string, containerName string) {
	umountVolumes(containerName, volumes)
	umountMerged(containerName)
	deleteWorkspace(containerName)
}

func umountVolumes(containerName string, volumes []string) {
	mergedPath := fmt.Sprintf(DefaultMergedPath, containerName)
	for _, volume := range volumes {
		volMapping := strings.Split(volume, ":")
		mountPoint := filepath.Join(mergedPath, volMapping[1])
		cmd := exec.Command("umount", mountPoint)
		if err := cmd.Run(); err != nil {
			logrus.Errorf("Fail to umount volume %s, error %v", volume, err)
		}
	}
}

func umountMerged(containerName string) {
	mergedPath := fmt.Sprintf(DefaultMergedPath, containerName)

	cmd := exec.Command("umount", mergedPath)
	if err := cmd.Run(); err != nil {
		logrus.Errorf("Fail to umount merged dir %s. error %v", mergedPath, err)
	}
	if err := os.RemoveAll(mergedPath); err != nil {
		logrus.Errorf("Fail to delete dir %s. error %v", mergedPath, err)
	}
}

func deleteWorkspace(containerName string) {
	workPath := filepath.Join(DefaultRootPath, containerName)
	if err := os.RemoveAll(workPath); err != nil {
		logrus.Errorf("Fail to delete dir %s. error %v", workPath, err)
	}
}

// func pathExists(path string) (bool, error) {
// 	_, err := os.Stat(path)
// 	if err == nil {
// 		return true, nil
// 	}
// 	if os.IsNotExist(err) {
// 		return false, nil
// 	}
// 	return false, err
// }
