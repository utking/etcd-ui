package types

type SensitiveStr string

func (s SensitiveStr) String() string {
	return "***"
}

func (s SensitiveStr) Unwrap() string {
	return string(s)
}

type KVItem struct {
	Key      string
	Value    SensitiveStr
	LeaseTTL int64
	LeaseID  int64
}
