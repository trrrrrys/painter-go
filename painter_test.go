package painter_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/trrrrrys/painter-go"
	internal "github.com/trrrrrys/painter-go/internal/painter"
)

func TestPainter(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    string
		options []painter.PaletteOption
	}{
		{
			"1. Hello world",
			"Hello world !!",
			"Hello world !!",
			nil,
		},
		{
			"2. Debug log",
			"DEBUG: this is debug log.",
			internal.ColorYellow.Sprintf("DEBUG: this is debug log."),
			nil,
		},
		{
			"3. Error log",
			"ERROR: this is error log.",
			internal.ColorRed.Sprintf("ERROR: this is error log."),
			nil,
		},
		{
			"4. Invalid format",
			"Error: this is error log.",
			"Error: this is error log.",
			nil,
		},
		{
			"5. json",
			`{"key": "value"}`,
			"{\n  \"key\": \"value\"\n}",
			[]painter.PaletteOption{
				painter.WithEnableJSONIndent(true),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			tt.options = append(tt.options, painter.WithPaletteOutput(&buf))
			p := painter.NewPalette(
				tt.options...,
			)
			fmt.Fprintf(p, tt.in)
			if diff := cmp.Diff(buf.String(), tt.want); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
