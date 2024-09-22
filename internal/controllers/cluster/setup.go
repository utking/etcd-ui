package cluster

import (
	"github.com/labstack/echo/v4"
	types "github.com/utking/etcd-ui/internal/types/common"
)

func Setup(app *echo.Echo, webMenu *types.WebMenu) {
	router := app.Group("/cluster")

	router.GET("/stats", clusterIndex)
	router.GET("/elect/:id", electNewLeader)

	router.GET("/leases/list", indexLeases)
	router.GET("/lease/:id", infoLease)

	router.GET("/lease/create", createLease)
	router.GET("/lease/edit/:id", createLease)
	router.POST("/lease/create", createLease)
	router.POST("/lease/delete", deleteLease)

	router.GET("/users/list", indexUsers)
	router.GET("/user", infoUser)
	router.GET("/user/edit/:name", createUser)
	router.GET("/user/create", createUser)
	router.POST("/user/create", createUser)
	router.POST("/user/delete", deleteUser)

	router.GET("/roles/list", indexRoles)
	router.GET("/role", infoRole)
	router.GET("/role/edit/:name", createRole)
	router.GET("/role/create", createRole)
	router.POST("/role/create", createRole)
	router.POST("/role/delete", deleteRole)

	router.GET("/keys/list", listKeys)
	router.GET("/key", infoKey)

	router.GET("/key/create", createKey)
	router.GET("/key/edit", createKey)
	router.POST("/key/create", createKey)
	router.POST("/key/delete", deleteKey)

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
				URIPath: "/cluster/keys/list",
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
				URIPath: "/cluster/leases/list",
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
				URIPath: "/cluster/users/list",
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
				URIPath: "/cluster/roles/list",
			},
		},
	})
}
