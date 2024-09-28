package types

type PermType uint8

func (p PermType) String() string {
	switch p {
	case PermRead:
		return "Read"
	case PermWrite:
		return "Write"
	case PermReadWrite:
		return "ReadWrite"
	default:
		return "Unknown"
	}
}

const (
	PermRead      = PermType(0)
	PermWrite     = PermType(1)
	PermReadWrite = PermType(2)
)

type KVPerm struct {
	Key      string
	RangeEnd string
	Type     PermType
}

type RoleInfo struct {
	Name  string
	Perms []KVPerm
}
