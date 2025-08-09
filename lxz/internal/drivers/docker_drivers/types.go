package docker_drivers

type ContainerData struct {
	ID          string
	Image       string
	Name        string
	Status      string
	RunningTime int64
}
