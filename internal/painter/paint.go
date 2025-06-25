package painter

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

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
	enableDefaultSetting *bool
	enableRegex          *bool
	Keywords             []*KeywordColor
	paletteBytes         []byte
	output               io.Writer
	enableJSONIndent     bool
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

func WithEnableDefaultSetting(enabled bool) func(*PaletteOptions) {
	return func(o *PaletteOptions) {
		enabled := enabled
		o.enableDefaultSetting = &enabled
	}
}

func WithEnableRegex(enabled bool) func(*PaletteOptions) {
	return func(o *PaletteOptions) {
		enabled := enabled
		o.enableRegex = &enabled
	}
}

func WithKeyword(keyword, color string, onlyKeyword bool) func(*PaletteOptions) {
	return func(o *PaletteOptions) {
		o.Keywords = append(o.Keywords, &KeywordColor{Keyword: keyword, Color: color, Partial: onlyKeyword})
	}
}

func NewPalette(opts ...PaletteOption) *Palette {
	var o PaletteOptions
	for _, opt := range opts {
		opt(&o)
	}
	var p Palette
	if o.enableDefaultSetting == nil || *o.enableDefaultSetting {
		if len(o.paletteBytes) == 0 {
			o.paletteBytes = defaultPalette
		}
		if err := yaml.Unmarshal(o.paletteBytes, &p); err != nil {
			return nil
		}
	}
	if o.enableRegex != nil {
		p.Regex = *o.enableRegex
	}
	if len(o.Keywords) > 0 {
		p.Keywords = append(p.Keywords, o.Keywords...)
	}
	if p.Regex {
		for _, v := range p.Keywords {
			p.keywordRegex = append(p.keywordRegex, regexp.MustCompile(v.Keyword))
		}
	}
	if o.output != nil {
		p.out = o.output
	} else {
		p.out = os.Stdout
	}
	p.enableJSONIndent = o.enableJSONIndent
	return &p
}

func (p *Palette) setColor(str string) string {
	coloredStr := strings.ReplaceAll(str, "%", "%%")
	if p.Regex {
		for i, r := range p.keywordRegex {
			if r.MatchString(coloredStr) {
				c := colorMap[p.Keywords[i].Color]
				if p.Keywords[i].Partial {
					coloredStr = r.ReplaceAllStringFunc(coloredStr, func(match string) string {
						subs := r.FindStringSubmatch(match)
						if len(subs) > 1 {
							return strings.ReplaceAll(match, subs[1], c.Sprintf(subs[1]))
						}
						return c.Sprintf(match)
					})
				} else {
					return c.Sprintf(coloredStr)
				}
			}
		}
	} else {
		for _, v := range p.Keywords {
			if strings.Contains(coloredStr, v.Keyword) {
				if v.Partial {
					coloredStr = strings.ReplaceAll(coloredStr, v.Keyword, colorMap[v.Color].Sprintf(v.Keyword))
				} else {
					return fmt.Sprint(colorMap[v.Color].Sprintf(coloredStr))
				}
			}
		}
	}
	return coloredStr
}

func (p *Palette) Paint(str string) {
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
	Partial bool   `yaml:"Partial"`
}
