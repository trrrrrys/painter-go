package main

import (
	"bufio"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// options
var (
	settingFile string
)

func init() {
	flag.StringVar(&settingFile, "f", "", "setting yaml file")
}

//go:embed example.yaml
var yamlData []byte

const (
	// man console_codes
	colorRed    = iota + 31 // 31
	colorGreen              // 32
	colorYellow             // 33
	colorBlue               // 33
)

var out = os.Stdout

var colorMap = map[string]int{
	"red":    colorRed,
	"green":  colorGreen,
	"yellow": colorYellow,
	"blue":   colorBlue,
}

type Palette struct {
	Keywords []*Keyword `yaml:"Keywords"`
}

func (s *Palette) Painting(str string) {
	for _, v := range s.Keywords {
		if strings.Contains(str, v.Word) {
			fmt.Fprintf(out, "\x1b[%vm%s\x1b[0m\n", colorMap[v.Color], str)
			return
		}
	}
	fmt.Fprintf(out, "%s\n", str)
}

type Keyword struct {
	Word  string `yaml:"Keyword"`
	Color string `yaml:"Color"`
}

func main() {
	flag.Parse()
	var b []byte
	_, err := os.Stat(settingFile)
	if err != nil {
		b = yamlData
	} else {
		f, err := os.ReadFile(settingFile)
		if err != nil {
			panic(err)
		}
		b = f
	}
	var v Palette
	if err := yaml.Unmarshal(b, &v); err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()
		v.Painting(s)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
