package fs

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

const (
	defaultRootPath = "./root"
)

func NewWorkspace() string {
	lowerDir := createLowerDir()
	diffDir := createDiffDir()
	workDir := createWorkDir()
	return mountMergeDir(lowerDir, diffDir, workDir)
}

// TODO: 暂未使用, 根据镜像创建lower dir
func createLowerDir() string {
	return filepath.Join(defaultRootPath, "busybox")
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

func mkdir(path, dir string) string {
	fullPath := filepath.Join(path, dir)
	if err := os.Mkdir(fullPath, 0777); err != nil {
		logrus.Errorf("Fail to creat dir %s error. %v", fullPath, err)
	}
	return fullPath
}

func DestroyWorkspace() {
	umountMerged()
	deleteWorkDir()
	deleteDiffDir()
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
