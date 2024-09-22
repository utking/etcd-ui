package types

import (
	"go.etcd.io/etcd/api/v3/etcdserverpb"
)

type MemberRecord struct {
	Name       string
	Version    string
	PeerURLs   []string
	ClientURLs []string
	Health     EndpointStatusRecord
	ID         uint64
}

type EndpointStatusRecord struct {
	Version  string
	Errors   []string
	ID       uint64
	DBSize   int64
	IsMaster bool
}

type ClusterStats struct {
	Members   []*etcdserverpb.Member
	ClusterID uint64
	MemberID  uint64
	RaftTerm  uint64
}
