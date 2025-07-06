package static

import "embed"

//go:embed *
var publicFiles embed.FS

func FS() embed.FS {
	return publicFiles
}
