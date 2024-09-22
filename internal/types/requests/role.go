package requests

import "github.com/utking/etcd-ui/internal/providers/etcd/types"

type KVPerm struct {
	Key     string `json:"key" form:"key"`
	IsRange bool   `json:"is_range" form:"is_range"`
}

type RoleCreate struct {
	Name       string   `json:"name" form:"name"`
	ReadPerms  []KVPerm `json:"read_perm" form:"read_perm"`
	WritePerms []KVPerm `json:"write_perm" form:"write_perm"`
}

func (r *RoleCreate) GetReadPerms() []types.KVPerm {
	return []types.KVPerm{}
}

func (r *RoleCreate) GetWritePerms() []types.KVPerm {
	return []types.KVPerm{}
}
