package fs

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	defaultRootPath = "./root"
)

func NewWorkspace(volumes []string) string {
	lowerDir := createLowerDir()
	diffDir := createDiffDir()
	workDir := createWorkDir()
	mergedDir := mountMergeDir(lowerDir, diffDir, workDir)
	mountVolumes(mergedDir, volumes)
	return mergedDir
}

func createLowerDir() string {
	lowerDir := filepath.Join(defaultRootPath, "busybox")
	if err := os.Mkdir(lowerDir, 0777); err != nil {
		logrus.Errorf("Fail to create lower dir %s, error %v", lowerDir, err)
	}
	cmd := exec.Command("tar", "-xvf", lowerDir+".tar", "-C", lowerDir)
	if err := cmd.Run(); err != nil {
		logrus.Errorf("Fail to tar busybox")
	}
	return lowerDir
}

func createDiffDir() string {
	return mkdir(defaultRootPath, "diff")
}

func createWorkDir() string {
	return mkdir(defaultRootPath, "work")
}

func mountMergeDir(lowPath, diffPath, workPath string) string {
	mountPath := mkdir(defaultRootPath, "merged")

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

func mkdir(path, dir string) string {
	fullPath := filepath.Join(path, dir)
	if err := os.Mkdir(fullPath, 0777); err != nil {
		logrus.Errorf("Fail to creat dir %s error. %v", fullPath, err)
	}
	return fullPath
}

func DestroyWorkspace(volumes []string) {
	umountVolumes(volumes)
	umountMerged()
	deleteWorkDir()
	deleteDiffDir()
	deleteLowerDir()
}

func umountVolumes(volumes []string) {
	for _, volume := range volumes {
		volMapping := strings.Split(volume, ":")
		// TODO:
		mountPoint := filepath.Join(defaultRootPath, "merged", volMapping[1])
		cmd := exec.Command("umount", mountPoint)
		if err := cmd.Run(); err != nil {
			logrus.Errorf("Fail to umount volume %s, error %v", volume, err)
		}
	}
}

func umountMerged() {
	mergedPath := filepath.Join(defaultRootPath, "merged")

	cmd := exec.Command("umount", mergedPath)
	if err := cmd.Run(); err != nil {
		logrus.Errorf("Fail to umount merged dir %s. error %v", mergedPath, err)
	}
	if err := os.RemoveAll(mergedPath); err != nil {
		logrus.Errorf("Fail to delete dir %s. error %v", mergedPath, err)
	}
}

func deleteWorkDir() {
	workPath := filepath.Join(defaultRootPath, "work")
	if err := os.RemoveAll(workPath); err != nil {
		logrus.Errorf("Fail to delete dir %s. error %v", workPath, err)
	}
}

func deleteDiffDir() {
	diffPath := filepath.Join(defaultRootPath, "diff")
	if err := os.RemoveAll(diffPath); err != nil {
		logrus.Errorf("Fail to delete dir %s. error %v", diffPath, err)
	}
}

func deleteLowerDir() {
	lowerPath := filepath.Join(defaultRootPath, "busybox")
	if err := os.RemoveAll(lowerPath); err != nil {
		logrus.Errorf("Fail to delete dir %s. error %v", lowerPath, err)
	}
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
