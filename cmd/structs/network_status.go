package structs

import (
	"encoding/json"
)

type IpType int
type ExternalNetworkState int
type DHCP int
type WifiStatus int
type RemoteManagement int

const (
	UnknownIpType IpType = iota
	DynamicIp
	StaticIp
)

const (
	UnknownExternalNetworkState ExternalNetworkState = iota
	Connected
	Disconnected
)

const (
	UnknownDHCP DHCP = iota
	Enabled
	Disabled
)

const (
	UnknownWifiStatus WifiStatus = iota
	WifiEnabledWithEncryption
	WifiEnabledWithoutEncryption
	WifiDisabled
)

const (
	UnknownRemoteManagementStatus RemoteManagement = iota
	RemoteManagementEnabled
	RemoteManagementDisabled
)

type NetworkStatus struct {
	External ExternalStatus   `json:"w1external"`
	Internal InternalStatus   `json:"internal"`
	Wireless []WirelessStatus `json:"wireless"`
	Extra    ExtraStatus      `json:"extra"`
}

type ExternalStatus struct {
	IpType        IpType               `json:"info"`
	ExternalIp    string               `json:"extip"`
	ConnectedTime string               `json:"time"`
	State         ExternalNetworkState `json:"state"`
}

type InternalStatus struct {
	Ip      string `json:"ip"`
	Dhcp    DHCP   `json:"dhcp"`
	StartIp string `json:"sip"`
	EndIp   string `json:"eip"`
}

type WirelessStatus struct {
	Ssid     string     `json:"ssid"`
	Status   WifiStatus `json:"mode"`
	Mhz      string     `json:"mhzstr"`
	Password string     `json:"wpapsk"`
}

type ExtraStatus struct {
	RemoteFlag RemoteManagement `json:"remoteflag"`
	RemotePort string           `json:"remoteport"`
	Version    string           `json:"version"`
	Uptime     string           `json:"time"`
}

func (i *IpType) UnmarshalJSON(data []byte) error {
	var ipTypeStr string
	if err := json.Unmarshal(data, &ipTypeStr); err != nil {
		return err
	}

	switch ipTypeStr {
	case "dynamic":
		*i = DynamicIp
	case "static":
		*i = StaticIp
	default:
		*i = UnknownIpType
	}

	return nil
}

func (e *ExternalNetworkState) UnmarshalJSON(data []byte) error {
	var networkStateStr string
	if err := json.Unmarshal(data, &networkStateStr); err != nil {
		return err
	}

	switch networkStateStr {
	case "D_DESC_WAN_PORT_DISCONNECTED.":
		*e = Disconnected
	case "D_SYSINFO_INTERNET_CONN_SUCCESS.":
		*e = Connected
	default:
		*e = UnknownExternalNetworkState
	}

	return nil
}

func (d *DHCP) UnmarshalJSON(data []byte) error {
	var dhcpStr string
	if err := json.Unmarshal(data, &dhcpStr); err != nil {
		return err
	}

	switch dhcpStr {
	case "D_NETCONF_INTERNAL_DHCP_RUNNING.":
		*d = Enabled
	case "D_NETCONF_INTERNAL_DHCP_STOPPED.":
		*d = Disabled
	default:
		*d = UnknownDHCP
	}

	return nil
}

func (w *WifiStatus) UnmarshalJSON(data []byte) error {
	var wifiStatusStr string
	if err := json.Unmarshal(data, &wifiStatusStr); err != nil {
		return err
	}

	switch wifiStatusStr {
	case "D_STOPPED_STATUS.":
		*w = WifiDisabled
	case "D_STARTED_STATUS.D_SYSTEM_INFO_WIRELESS_ENC_ENABLE.":
		*w = WifiEnabledWithEncryption
	case "D_STARTED_STATUS.D_SYSTEM_INFO_WIRELESS_ENC_DISABLE.":
		*w = WifiEnabledWithoutEncryption
	default:
		*w = UnknownWifiStatus
	}

	return nil
}

func (r *RemoteManagement) UnmarshalJSON(data []byte) error {
	var remoteManagementStr string
	if err := json.Unmarshal(data, &remoteManagementStr); err != nil {
		return err
	}

	switch remoteManagementStr {
	case "D_REMOTEMGMT_START_COMMENT.":
		*r = RemoteManagementEnabled
	case "D_REMOTEMGMT_STOP_COMMENT.":
		*r = RemoteManagementDisabled
	default:
		*r = UnknownRemoteManagementStatus
	}

	return nil
}
