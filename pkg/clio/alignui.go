package clio

import (
	"fmt"

	"github.com/asciifaceman/gomo/pkg/clients"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// AlignUI wraps the alignment UI
type AlignUI struct {
	version string
	Client  *clients.Polling
	Motd    *widgets.Paragraph
}

// NewAlignUI creates and returns a new alignment UI. pollFrequency in seconds
func NewAlignUI(version string, hostname string, timeout int, pollFrequency int) (*AlignUI, error) {
	c, err := clients.NewPolling(hostname, timeout, pollFrequency)
	if err != nil {
		return nil, err
	}

	a := &AlignUI{
		Client:  c,
		version: version,
	}

	return a, nil
}

// BuildUI does the initial setup of the UI
func (a *AlignUI) BuildUI() {
	a.Motd = widgets.NewParagraph()
	a.Motd.Title = fmt.Sprintf("Gomo Alignment: %s", a.version)
	a.Motd.Text = "PRESS q TO QUIT"
	a.Motd.SetRect(0, 0, 50, 5)
	a.Motd.TextStyle.Fg = ui.ColorWhite
	a.Motd.BorderStyle.Fg = ui.ColorCyan
}

// Run wraps the runtime providing a listener for response and draw callback

// Draw redraws the UI
