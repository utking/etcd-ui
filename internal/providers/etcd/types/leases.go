package types

type LeaseRecord struct {
	HexID      string
	Keys       []string
	ID         int64
	GrantedTTL int64
	TTL        int64
}

func (l LeaseRecord) KeysCounter() int {
	return len(l.Keys)
}
