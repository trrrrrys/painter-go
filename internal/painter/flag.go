package painter

import "flag"

// options
var (
	settingFile string
)

func init() {
	flag.StringVar(&settingFile, "f", "", "setting yaml file")
}
