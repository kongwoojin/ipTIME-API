package structs

type IpType int

const (
	Unknown IpType = iota
	DynamicIp
	StaticIp
)

type RouterStatus struct {
	Model               string
	IsInternetConnected bool
	IpType              IpType
	Ip                  string
	Fiveghz             WifiInfo
	Twoghz              WifiInfo
	FirmwareVersion     string
	RemoteAccess        bool
	RemoteAccessPort    int
	Uptime              Uptime
}
