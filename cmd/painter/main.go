package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	// man console_codes
	colorRed    = iota + 31 // 31
	colorGreen              // 32
	colorYellow             // 33
	colorBlue               // 33
)

var (
	stdin  = os.Stdin
	stdout = os.Stdout
)

var colorMap = map[string]int{
	"red":    colorRed,
	"green":  colorGreen,
	"yellow": colorYellow,
	"blue":   colorBlue,
}

type St struct {
	Keywords []*Keyword `yaml:"Keywords"`
}

func (s *St) Print(str string) {
	for _, v := range s.Keywords {
		if strings.Contains(str, v.Word) {
			fmt.Fprintf(stdout, "\x1b[%vm%s\x1b[0m\n", colorMap[v.Color], str)
			return
		}
	}
	fmt.Fprintf(stdout, "%s\n", str)
}

type Keyword struct {
	Word  string `yaml:"Keyword"`
	Color string `yaml:"Color"`
}

func main() {
	f, err := os.ReadFile("./sample.yaml")
	if err != nil {
		panic(err)
	}
	var v St
	if err := yaml.Unmarshal(f, &v); err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(stdin)
	for scanner.Scan() {
		s := scanner.Text()
		v.Print(s)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
