package requests

type LeaseCreate struct {
	LeaseID int64 `json:"lease_id" form:"lease_id"`
	TTL     int64 `json:"ttl" form:"ttl"`
}
