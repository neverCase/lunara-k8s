package kubernetes

const (
	_port          = 6379
	_targetPort    = 6379
	_containerPort = 6379
)

type Redis struct {
	Name string `json:"name"`

	Image        string `json:"image"`
	ImageVersion string `json:"image_version"`
	Replicas     int    `json:"replicas"`
}

func (r *Redis) NewRedis() {

}
