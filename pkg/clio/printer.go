package clio

import (
	"fmt"
	"strings"
)

const (
	HeaderBarBuffer    = 5
	DefaultHeaderWidth = 64
	DefaultKVWidth     = 32
	DefaultIndent      = 2
)

type Printer struct {
	HeaderWidth int
	KVWidth     int
	Indent      int
	LeftIndent  string
}

func NewPrinter(headerWidth int, kvWidth int, indent int) *Printer {
	return &Printer{
		HeaderWidth: headerWidth,
		KVWidth:     kvWidth,
		Indent:      indent,
		LeftIndent:  strings.Repeat(" ", indent),
	}
}

func (p *Printer) PrintHeader(title string) {
	//4
	rightBarLength := p.HeaderWidth - HeaderBarBuffer - len(title)
	for rightBarLength < 0 {
		rightBarLength++
	}
	bars := strings.Repeat("=", rightBarLength)
	fmt.Printf("=== %s %s\n", title, bars)
}

func (p *Printer) PrintKVIndent(key string, value interface{}) {
	// fmt.Printf("  RSSI: %5d\n", data.CellLTEStats[0].Stat.RSSICurrent)
	leftPad := p.KVWidth - (p.Indent + len(key))
	for leftPad < 0 {
		leftPad++
	}
	fmt.Printf("%s%s:%*v\n", p.LeftIndent, key, leftPad, value)
}

func (p *Printer) PrintKV(key string, value interface{}) {
	fmt.Printf("%s: %v\n", key, value)
}
