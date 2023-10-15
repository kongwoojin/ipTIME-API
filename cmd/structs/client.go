package structs

type ConnectionType int

const (
	UnknownConnectionType ConnectionType = iota
	Wired
	Wireless
)

type Client struct {
	IP                 string         `json:"ip"`
	Mac                string         `json:"mac"`
	Hostname           string         `json:"hostname"`
	ConnectionType     ConnectionType `json:"connectionType"`
	IsManuallyAssigned bool           `json:"isManuallyAssigned"`
}
