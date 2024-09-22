package requests

import "github.com/utking/etcd-ui/internal/providers/etcd/types"

type UserCreate struct {
	Name     string             `json:"name" form:"name"`
	Password types.SensitiveStr `json:"password" form:"password"`
	Roles    []string           `json:"roles" form:"roles"`
}
