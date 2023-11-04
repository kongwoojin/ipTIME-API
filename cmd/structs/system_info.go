package structs

type SystemInfo struct {
	ConnectInfo []ConnectionInfo `json:"connectinfo"`
	CanLogout   string           `json:"canlogout"`
	FirmVersion string           `json:"firmversion"`
	ProductName string           `json:"productname"`
}

type ConnectionInfo struct {
	WanName     string `json:"wanname"`
	WirelessWan string `json:"wirelesswan"`
	WanType     string `json:"wantype"`
	WanStatus   string `json:"wanstatus"`
	IpAddr      string `json:"ipaddr"`
}
