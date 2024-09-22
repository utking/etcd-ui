package types

type KVPerm struct {
	Key      string
	RangeEnd string
}

type RoleInfo struct {
	Name    string
	KVRead  []KVPerm
	KVWrite []KVPerm
}
