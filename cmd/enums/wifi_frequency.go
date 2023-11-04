package enums

import "github.com/kongwoojin/ipTIME-API/cmd/structs"

type WifiFrequency int

const (
	F2_4GHZ = iota + 1
	F5GHZ
)

func (f WifiFrequency) GetSSID(status *structs.RouterStatus) string {
	ssid := ""
	for _, wifi := range status.NetworkStatus.Wireless {
		if wifi.Mhz == f.GetFullFrequency() {
			ssid = wifi.Ssid
			break
		}
	}
	return ssid
}

func (f WifiFrequency) GetFullFrequency() string {
	if f == F2_4GHZ {
		return "2.4 GHz"
	} else {
		return "5 GHz"
	}
}

func (f WifiFrequency) GetFrequency() string {
	if f == F2_4GHZ {
		return "2g"
	} else {
		return "5g"
	}
}
