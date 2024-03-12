package painter

import "github.com/trrrrrys/painter-go/internal/painter"

type (
	PaletteOption = painter.PaletteOption
)

var NewPalette = painter.NewPalette

var (
	WithPaletteBytes         = painter.WithPaletteBytes
	WithPaletteOutput        = painter.WithPaletteOutput
	WithEnableJSONIndent     = painter.WithEnableJSONIndent
	WithEnableDefaultSetting = painter.WithEnableDefaultSetting
	WithEnableRegex          = painter.WithEnableRegex
	WithKeyword              = painter.WithKeyword
)
