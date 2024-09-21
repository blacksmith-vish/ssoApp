package embed

import "embed"

var (
	//go:embed migrations
	SQLiteMigrations embed.FS

	//go:embed static
	StaticFiles embed.FS
)
