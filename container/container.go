package container

const (
	RUNNING                    = "running"
	STOP                       = "stopped"
	Exit                       = "exited"
	DefaultInfoLocation        = "/var/run/mycontainer/%s/"
	ConfigName                 = "config.json"
	ContainerLogFile    string = "container.log"
)

type ContainerInfo struct {
	Pid         string `json:"pid"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	Command     string `json:"command"`
	CreatedTime string `json:"createdTime"`
	Status      string `json:"status"`
}
