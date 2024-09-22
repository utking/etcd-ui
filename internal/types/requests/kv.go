package requests

type KVCreate struct {
	Key     string `json:"key" form:"key"`
	Value   string `json:"value" form:"value"`
	LeaseID uint64 `json:"lease_id" form:"lease_id"`
	TTL     int64  `json:"ttl" form:"ttl"`
}
