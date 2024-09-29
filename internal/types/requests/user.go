package requests

import (
	"errors"
	"fmt"
	"strings"

	"github.com/utking/etcd-ui/internal/providers/etcd/types"
)

type UserCreate struct {
	Name         string             `json:"name" form:"name"`
	Password     types.SensitiveStr `json:"password" form:"password"`
	Confirmation types.SensitiveStr `json:"confirmation" form:"confirmation"`
}

func (c *UserCreate) normalize() {
	c.Password = types.SensitiveStr(strings.TrimSpace(c.Password.Unwrap()))
	c.Confirmation = types.SensitiveStr(strings.TrimSpace(c.Confirmation.Unwrap()))
	c.Name = strings.TrimSpace(c.Name)
}

func (c *UserCreate) Validate() error {
	var errList []error

	c.normalize()

	if c.Password.Unwrap() != c.Confirmation.Unwrap() {
		errList = append(errList, fmt.Errorf("the password and conrifmation must match"))
	}

	if c.Name == "" || c.Password.Unwrap() == "" {
		errList = append(errList, fmt.Errorf("username and password cannot be empty"))
	}

	return errors.Join(errList...)
}

type UserEditRoles struct {
	Name          string   `json:"name" form:"name"`
	ExistingRoles []string `json:"roles" form:"roles"`
	AddRoles      []string `json:"add_roles" form:"add_roles"`
	RevokeRoles   []string `json:"revoke_roles" form:"revoke_roles"`
}
