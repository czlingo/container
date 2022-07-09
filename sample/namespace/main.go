package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

const cgroupMemoryHierarchyMount = "/sys/fs/cgroup/memory"

func main() {
	cmd := exec.Command("sh") //用来指定被fork出来的新进程内的初始命令.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// https://github.com/czlingo/notes/tree/main/container/docker.md
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUSER |
			syscall.CLONE_NEWNET,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
