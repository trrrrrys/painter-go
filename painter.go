package painter

import "github.com/trrrrrys/painter-go/internal/painter"

type (
	PaletteOption = painter.PaletteOption
)

var (
	NewPalette = painter.NewPalette

	// options
	WithPaletteBytes     = painter.WithPaletteBytes
	WithPaletteOutput    = painter.WithPaletteOutput
	WithEnableJSONIndent = painter.WithEnableJSONIndent
)
