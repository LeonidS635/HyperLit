package static

import _ "embed"

//go:embed head.html
var Head string

//go:embed style.html
var Style string

//go:embed script.html
var Script string

//go:embed empty_state.html
var EmptyState string
