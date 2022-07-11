package container

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

const (
	defaultRootPath = "./root"
)

// type OverlayFS struct {
// 	LowerDir string
// 	UpperDir string
// 	WorkDir  string
// 	Merged   string
// }

func NewWorkspace() string {
	lowerPath := CreateLowerDir()
	diffPath := CreateDiffDir()
	workPath := CreateWorkDir()
	return MountMergeDir(lowerPath, diffPath, workPath)
}

// TODO: 暂未使用, 根据镜像创建lower dir
func CreateLowerDir() string {
	return filepath.Join(defaultRootPath, "busybox")
}

func CreateDiffDir() string {
	return mkdir(defaultRootPath, "diff")
}

func CreateWorkDir() string {
	return mkdir(defaultRootPath, "work")
}

func MountMergeDir(lowPath, diffPath, workPath string) string {
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
