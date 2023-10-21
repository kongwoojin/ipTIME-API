package structs

type PortForward struct {
	Name              string
	Enabled           bool
	IP                string
	Protocol          string
	ExternalPortStart int
	ExternalPortEnd   int
	InternalPortStart int
	InternalPortEnd   int
}
