package requests

import (
	"encoding/json"

	"github.com/utking/etcd-ui/internal/providers/etcd/types"
)

type KVPerm struct {
	Key      string `json:"key" form:"key"`
	RangeEnd string `json:"range_end" form:"range_end"`
	Type     types.PermType
}

func (p KVPerm) String() string {
	b64Val, err := json.Marshal(&p)
	if err != nil {
		return ""
	}

	return string(b64Val)
}

func (p KVPerm) From(in string) (KVPerm, error) {
	var item KVPerm

	err := json.Unmarshal([]byte(in), &item)
	if err != nil {
		return item, err
	}

	return item, nil
}

type RoleCreate struct {
	Name  string   `json:"name" form:"name"`
	Perms []KVPerm `json:"perms" form:"perms"`
}

type RoleRevokePerm struct {
	Name       string   `json:"name" form:"name"`
	PermHashes []string `json:"perms" form:"perms"` // base64 encoded [key, range)
}

func (rev *RoleRevokePerm) KVPerm() []KVPerm {
	result := make([]KVPerm, 0, len(rev.PermHashes))

	for _, item := range rev.PermHashes {
		result = append(result, KVPerm{Key: item})
	}

	return result
}
