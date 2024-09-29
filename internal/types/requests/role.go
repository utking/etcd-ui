package requests

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/utking/etcd-ui/internal/providers/etcd/types"
)

type KVPerm struct {
	Key      string `json:"key" form:"key"`
	RangeEnd string `json:"range_end" form:"range_end"`
	Type     types.PermType
	IsRange  bool
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

type RoleGrantPerm struct {
	Name     string         `json:"name" form:"name"`
	Key      string         `json:"key" form:"key"`
	RangeEnd string         `json:"range_end" form:"range_end"`
	SetRange string         `json:"is_range" form:"is_range"`
	Type     types.PermType `json:"type" form:"type"`
}

func (g *RoleGrantPerm) IsRange() bool {
	return strings.EqualFold(g.SetRange, "on")
}

func (g *RoleGrantPerm) Validate() error {
	var errList []error

	if !slices.Contains([]types.PermType{types.PermRead, types.PermWrite, types.PermReadWrite}, g.Type) {
		errList = append(errList, fmt.Errorf("unknown permissions type"))
	}

	if g.IsRange() && g.RangeEnd != "" {
		errList = append(errList, fmt.Errorf("is-range and range-end are mutually exclusive; use just one"))
	}

	return errors.Join(errList...)
}

func (g *RoleGrantPerm) KVPerm() types.KVPerm {
	return types.KVPerm{
		Key:      g.Key,
		RangeEnd: g.RangeEnd,
		Type:     g.Type,
		IsRange:  g.IsRange(),
	}
}
