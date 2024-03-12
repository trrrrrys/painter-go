package painter

import (
	"bufio"
	"errors"
	"flag"
	"os"
)

func Run() error {
	flag.Parse()
	b, _ := os.ReadFile(settingFile)
	palette := NewPalette(WithPaletteBytes(b))
	if palette == nil {
		return errors.New("failed to parse yaml")
	}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()
		palette.Paint(s)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
