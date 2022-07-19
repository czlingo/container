package container

const (
	RUNNING                 = "running"
	STOP                    = "stopped"
	Exit                    = "exited"
	ConfigName              = "config.json"
	ContainerLogFile string = "container.log"
)

type ContainerInfo struct {
	Pid         string   `json:"pid"`
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Command     string   `json:"command"`
	Status      string   `json:"status"`
	Volumes     []string `json:"volumes"`
	CreatedTime string   `json:"createdTime"`
	PortMapping []string `json:"portmapping"` //端口映射
}
