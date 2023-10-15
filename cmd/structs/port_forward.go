package structs

type PortForward struct {
	Name              string
	IP                string
	Protocol          string
	ExternalPortStart int
	ExternalPortEnd   int
	InternalPortStart int
	InternalPortEnd   int
}
