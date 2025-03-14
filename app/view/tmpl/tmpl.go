package tmpl

import "embed"

//go:embed *.html
var Content embed.FS
