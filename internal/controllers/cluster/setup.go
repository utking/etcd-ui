package cluster

import (
	"github.com/labstack/echo/v4"
	types "github.com/utking/etcd-ui/internal/types/common"
)

func Setup(app *echo.Echo, webMenu *types.WebMenu) {
	router := app.Group("/cluster")

	router.GET("/stats", clusterIndex)
	router.GET("/auth/:action", clusterFlipAuth)
	router.GET("/elect/:id", electNewLeader)

	router.GET("/leases", indexLeases)
	router.GET("/lease/:id", infoLease)

	router.Match([]string{"GET", "POST"}, "/lease/create", createLease)
	router.GET("/lease/edit/:id", createLease)
	router.POST("/lease/delete", deleteLease)

	router.GET("/user", infoUser)
	router.GET("/users", indexUsers)
	router.POST("/user/edit", revokeUserGroups)
	router.POST("/user/delete", deleteUser)
	router.Match([]string{"GET", "POST"}, "/user/roles/add", addUserGroups)
	router.Match([]string{"GET", "POST"}, "/user/create", createUser)
	router.Match([]string{"GET", "POST"}, "/user/passwd", passwdUser)

	router.GET("/role", infoRole)
	router.GET("/roles", indexRoles)
	router.GET("/role/edit/:name", editRolePermissions)
	router.POST("/role/revoke/:name", revokeRolePermissions)
	router.POST("/role/grant/:name", grantRolePermissions)
	router.POST("/role/delete", deleteRole)
	router.Match([]string{"GET", "POST"}, "/role/create", createRole)

	router.GET("/key", infoKey)
	router.GET("/keys", listKeys)

	router.GET("/key/edit", createKey)
	router.POST("/key/delete", deleteKey)
	router.Match([]string{"GET", "POST"}, "/key/create", createKey)

	webMenu.UserItems = append(webMenu.UserItems, map[string][]types.WebMenuItem{
		"Cluster": {
			{
				Type:    "item",
				Title:   "Stats",
				URIPath: "/cluster/stats",
			},
		},
		"Keys": {
			{
				Type:    "item",
				Title:   "Add",
				URIPath: "/cluster/key/create",
			},
			{
				Title: "-",
			},
			{
				Type:    "item",
				Title:   "List",
				URIPath: "/cluster/keys",
			},
		},
		"Leases": {
			{
				Type:    "item",
				Title:   "Add",
				URIPath: "/cluster/lease/create",
			},
			{
				Title: "-",
			},
			{
				Type:    "item",
				Title:   "List",
				URIPath: "/cluster/leases",
			},
		},
		"Users": {
			{
				Type:    "item",
				Title:   "Add",
				URIPath: "/cluster/user/create",
			},
			{
				Title: "-",
			},
			{
				Type:    "item",
				Title:   "List",
				URIPath: "/cluster/users",
			},
		},
		"Roles": {
			{
				Type:    "item",
				Title:   "Add",
				URIPath: "/cluster/role/create",
			},
			{
				Title: "-",
			},
			{
				Type:    "item",
				Title:   "List",
				URIPath: "/cluster/roles",
			},
		},
	})
}
