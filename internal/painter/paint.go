package painter

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

//go:embed default.yaml
var defaultPalette []byte

type Palette struct {
	Regex            bool            `yaml:"Regex"`
	Keywords         []*KeywordColor `yaml:"Keywords"`
	keywordRegex     []*regexp.Regexp
	out              io.Writer
	enableJSONIndent bool
}

type PaletteOptions struct {
	paletteBytes     []byte
	output           io.Writer
	enableJSONIndent bool
}

type PaletteOption func(*PaletteOptions)

func WithPaletteBytes(b []byte) func(*PaletteOptions) {
	return func(o *PaletteOptions) {
		o.paletteBytes = b
	}
}
func WithPaletteOutput(w io.Writer) func(*PaletteOptions) {
	return func(o *PaletteOptions) {
		o.output = w
	}
}

func WithEnableJSONIndent(enabled bool) func(*PaletteOptions) {
	return func(o *PaletteOptions) {
		o.enableJSONIndent = enabled
	}
}

func NewPalette(opts ...PaletteOption) *Palette {
	var o PaletteOptions
	for _, opt := range opts {
		opt(&o)
	}
	if len(o.paletteBytes) == 0 {
		o.paletteBytes = defaultPalette
	}
	var p Palette
	if err := yaml.Unmarshal(o.paletteBytes, &p); err != nil {
		return nil
	}
	p.init(o)
	return &p
}

func (p *Palette) init(options PaletteOptions) {
	if p.Regex {
		for _, v := range p.Keywords {
			p.keywordRegex = append(p.keywordRegex, regexp.MustCompile(v.Keyword))
		}
	}
	if options.output != nil {
		p.out = options.output
	} else {
		p.out = os.Stdout
	}
	p.enableJSONIndent = options.enableJSONIndent
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
	if p.enableJSONIndent {
		var buf bytes.Buffer
		err := json.Indent(&buf, b, "", "  ")
		if err == nil {
			return fmt.Fprint(p.out, p.setColor(buf.String()))
		}
	}
	return fmt.Fprint(p.out, p.setColor(string(b)))
}

type KeywordColor struct {
	Keyword string `yaml:"Keyword"`
	Color   string `yaml:"Color"`
}
