package assets

import "embed"

// Existing global FS — keep as-is
//go:embed *
var FS embed.FS

// Embed only SQL migration files
//go:embed migrations/*.sql
var MigrationFS embed.FS

//go:embed seeders/*.sql
var SeederFS embed.FS
