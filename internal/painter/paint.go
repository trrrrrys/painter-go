package painter

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

//go:embed default.yaml
var defaultPalette []byte

type Palette struct {
	Regex        bool            `yaml:"Regex"`
	Keywords     []*KeywordColor `yaml:"Keywords"`
	keywordRegex []*regexp.Regexp
	out          io.Writer
}

type PaletteOptions struct {
	PaletteBytes []byte
	Output       io.Writer
}

type PaletteOption func(*PaletteOptions)

func WithPaletteBytes(b []byte) func(*PaletteOptions) {
	return func(o *PaletteOptions) {
		o.PaletteBytes = b
	}
}
func WithPaletteOutput(w io.Writer) func(*PaletteOptions) {
	return func(o *PaletteOptions) {
		o.Output = w
	}
}

func NewPalette(opts ...PaletteOption) *Palette {
	var o PaletteOptions
	for _, opt := range opts {
		opt(&o)
	}
	if len(o.PaletteBytes) == 0 {
		o.PaletteBytes = defaultPalette
	}
	var p Palette
	if err := yaml.Unmarshal(o.PaletteBytes, &p); err != nil {
		return nil
	}
	if o.Output != nil {
		p.out = o.Output
	}
	p.init()
	return &p
}

func (p *Palette) init() {
	if p.Regex {
		for _, v := range p.Keywords {
			p.keywordRegex = append(p.keywordRegex, regexp.MustCompile(v.Keyword))
		}
	}
	if p.out == nil {
		p.out = os.Stdout
	}
}

func (p *Palette) setColor(str string) string {
	if p.Regex {
		for i, r := range p.keywordRegex {
			if r.MatchString(str) {
				return fmt.Sprint(colorMap[p.Keywords[i].Color].Sprintf(str))
			}
		}
	} else {
		for _, v := range p.Keywords {
			return fmt.Sprint(colorMap[v.Color].Sprintf(str))
		}
	}
	return str
}

func (p *Palette) Painting(str string) {
	fmt.Fprint(p.out, p.setColor(str)+"\n")
}

func (p *Palette) Write(b []byte) (int, error) {
	fmt.Fprint(p.out, p.setColor(string(b)))
	return len(b), nil
}

type KeywordColor struct {
	Keyword string `yaml:"Keyword"`
	Color   string `yaml:"Color"`
}
