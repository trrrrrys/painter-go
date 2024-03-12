package painter

import "fmt"

type Color string

func (c Color) Sprintf(format string, args ...interface{}) string {
	return fmt.Sprintf("\x1b[%vm%s\x1b[0m", c, fmt.Sprintf(format, args...))
}

const (
	// man console_codes
	ColorWhite   Color = "0;30"
	ColorRed     Color = "0;31"
	ColorGreen   Color = "0;32"
	ColorYellow  Color = "1;33"
	ColorBlue    Color = "0;34"
	ColorMagenta Color = "0;35"
	ColorCyan    Color = "0;36"
	ColorGray    Color = "1;30"
)

var colorMap = map[string]Color{
	"white":   ColorWhite,
	"red":     ColorRed,
	"green":   ColorGreen,
	"yellow":  ColorYellow,
	"blue":    ColorBlue,
	"magenta": ColorMagenta,
	"cyan":    ColorCyan,
	"gray":    ColorGray,
}
