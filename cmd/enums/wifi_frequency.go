package enums

import "github.com/kongwoojin/ipTIME-API/cmd/structs"

type WifiFrequency int

const (
	F2_4GHZ = iota + 1
	F5GHZ
)

func (f WifiFrequency) GetSSID(status *structs.RouterStatus) string {
	if f == F2_4GHZ {
		return status.Twoghz.Ssid
	} else {
		return status.Fiveghz.Ssid
	}
}

func (f WifiFrequency) GetFrequency() string {
	if f == F2_4GHZ {
		return "2g"
	} else {
		return "5g"
	}
}
