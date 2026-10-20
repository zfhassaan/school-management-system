package migrations

import "embed"

//go:embed *.sql
var SQLiteFS embed.FS

