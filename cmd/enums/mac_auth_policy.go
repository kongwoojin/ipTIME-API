package enums

type MacAuthPolicy int

const (
	MacAuthDisabled = iota + 1
	WhiteList
	BlackList
)

func (m MacAuthPolicy) GetPolicy() string {
	switch m {
	case 1:
		return "open"
	case 2:
		return "deny" // Not a typo. accept and deny are swapped in ipTIME routers
	case 3:
		return "accept"
	default:
		return "open"
	}
}
