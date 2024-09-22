package common

// WebMenu is the struct for the web menu.
// It has two slices of maps, one for the user items and one for the admin items.
type WebMenu struct {
	UserItems  []map[string][]WebMenuItem
	AdminItems []map[string][]WebMenuItem
}

// WebMenuItem is the struct for a Web menu item.
type WebMenuItem struct {
	Type     string
	Title    string
	URIPath  string
	Target   string
	SubItems []WebMenuItem
}
