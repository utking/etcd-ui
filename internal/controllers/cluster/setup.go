package cluster

import (
	"github.com/labstack/echo/v4"
	types "github.com/utking/etcd-ui/internal/types/common"
)

func Setup(app *echo.Echo, webMenu *types.WebMenu) {
	router := app.Group("/cluster")

	router.GET("/stats", clusterIndex)
	router.GET("/elect/:id", electNewLeader)

	router.GET("/leases", indexLeases)
	router.GET("/lease/:id", infoLease)

	router.GET("/lease/create", createLease)
	router.GET("/lease/edit/:id", createLease)
	router.POST("/lease/create", createLease)
	router.POST("/lease/delete", deleteLease)

	router.GET("/users", indexUsers)
	router.GET("/user/:name", infoUser)

	router.GET("/roles", indexRoles)
	router.GET("/role/:name", infoRole)

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
			{
				Title: "-",
			},
			{
				Type:    "item",
				Title:   "Keys",
				URIPath: "/cluster/keys/list",
			},
			{
				Type:    "item",
				Title:   "Leases",
				URIPath: "/cluster/leases",
			},
			{
				Title: "-",
			},
			{
				Type:    "item",
				Title:   "Roles",
				URIPath: "/cluster/roles",
			},
			{
				Type:    "item",
				Title:   "Users",
				URIPath: "/cluster/users",
			},
		},
	})
}
