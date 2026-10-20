package frontend

import "embed"

// Assets embeds the production build of the Vue 3 frontend application.
//
//go:embed all:dist
var Assets embed.FS

