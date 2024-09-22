package static

import "embed"

var (
	//go:embed *
	StaticFiles embed.FS
)
