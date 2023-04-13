package alignui

import (
	"fmt"
	"time"

	"github.com/asciifaceman/gomo/pkg/clients"
	"github.com/asciifaceman/gomo/pkg/models"
	"github.com/asciifaceman/gomo/pkg/radiofreq"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const (
	S5G = iota
	SLTE
	LHEADER      = "lheader"
	RHEADER      = "rheader"
	PLOT5G       = "plot5G"
	KEY5GSNR     = "key5GSNR"
	KEY5GRSRP    = "key5Grsrp"
	KEY5GRSRQ    = "key5Grsrq"
	KEY5GBAND    = "key5Gband"
	KEY5GCELLID  = "key5Gcellid"
	PLOTLTE      = "plotLTE"
	KEYLTESNR    = "keyLTEsnr"
	KEYLTERSRP   = "keyLTErsrp"
	KEYLTERSRQ   = "keyLTErsrq"
	KEYLTEBAND   = "keyLTEband"
	KEYLTECELLID = "keyLTEcellid"
	ALERTS       = "alerts"
)

var (
	graphStats   = []string{KEY5GSNR, KEY5GRSRP, KEY5GRSRQ, KEYLTESNR, KEYLTERSRP, KEYLTERSRQ}
	runningSteps = []string{"|", "/", "--", "\\", "|", "/", "--", "\\"}
	bandmap      = radiofreq.BandMap.Map()
)

type QualityStat struct {
	slice []float64
	max   float64
}

func NewQualityStat() *QualityStat {
	return &QualityStat{
		slice: make([]float64, 0),
		max:   float64(0),
	}
}

type AlignmentUI struct {
	tw       int
	th       int
	grid     *ui.Grid
	elements map[string]ui.Drawable
	events   <-chan ui.Event
	client   *clients.Polling
	data     chan *models.FastmileReturn
	cancel   chan interface{}
	stats    map[string]*QualityStat
	silent   bool
}

// New initializes the UI and prepares to run
// if you must abandon before running be sure to call AlignmentUI.Close()
func New(hostname string, timeout int, pollFrequency int, silent bool) (*AlignmentUI, error) {
	if err := ui.Init(); err != nil {
		return nil, err
	}

	tw, th := ui.TerminalDimensions()

	p, err := clients.NewPolling(hostname, timeout, pollFrequency)
	if err != nil {
		return nil, err
	}

	a := &AlignmentUI{
		tw:       tw,
		th:       th,
		elements: make(map[string]ui.Drawable),
		events:   ui.PollEvents(),
		client:   p,
		data:     make(chan *models.FastmileReturn, 1),
		cancel:   make(chan interface{}),
		stats:    make(map[string]*QualityStat),
		silent:   silent,
	}

	for _, stat := range graphStats {
		a.stats[stat] = NewQualityStat()
	}

	a.buildUI()

	return a, nil
}

func (a *AlignmentUI) buildUI() {
	a.elements[LHEADER] = widgets.NewParagraph()
	a.elements[LHEADER].(*widgets.Paragraph).Title = "Gomo Alignment"
	a.elements[LHEADER].(*widgets.Paragraph).Text = `Ideal signal is all lines at 1.0
	However that is extremely unlikely so you 
	should seek consistently highest values.`
	a.elements[LHEADER].(*widgets.Paragraph).TextStyle.Fg = ui.ColorCyan

	a.elements[RHEADER] = widgets.NewParagraph()
	a.elements[RHEADER].(*widgets.Paragraph).Text = `PRESS q TO QUIT`
	a.elements[RHEADER].(*widgets.Paragraph).TextStyle.Fg = ui.ColorGreen

	a.elements[PLOT5G] = widgets.NewPlot()
	a.elements[PLOT5G].(*widgets.Plot).Title = " 5G "
	a.elements[PLOT5G].(*widgets.Plot).Data = make([][]float64, 3)
	a.elements[PLOT5G].(*widgets.Plot).AxesColor = ui.ColorWhite
	a.elements[PLOT5G].(*widgets.Plot).LineColors[0] = ui.ColorRed
	a.elements[PLOT5G].(*widgets.Plot).LineColors[1] = ui.ColorBlue
	a.elements[PLOT5G].(*widgets.Plot).LineColors[2] = ui.ColorYellow
	a.elements[PLOT5G].(*widgets.Plot).Marker = widgets.MarkerDot

	a.elements[KEY5GSNR] = widgets.NewParagraph()
	a.elements[KEY5GSNR].(*widgets.Paragraph).Text = "SNR"
	a.elements[KEY5GSNR].(*widgets.Paragraph).TextStyle.Fg = ui.ColorRed
	a.elements[KEY5GRSRP] = widgets.NewParagraph()
	a.elements[KEY5GRSRP].(*widgets.Paragraph).Text = "RSRP"
	a.elements[KEY5GRSRP].(*widgets.Paragraph).TextStyle.Fg = ui.ColorBlue
	a.elements[KEY5GRSRQ] = widgets.NewParagraph()
	a.elements[KEY5GRSRQ].(*widgets.Paragraph).Text = "RSRQ"
	a.elements[KEY5GRSRQ].(*widgets.Paragraph).TextStyle.Fg = ui.ColorYellow
	a.elements[KEY5GBAND] = widgets.NewParagraph()
	a.elements[KEY5GBAND].(*widgets.Paragraph).Title = " Band "
	a.elements[KEY5GBAND].(*widgets.Paragraph).Text = "N/A"
	a.elements[KEY5GCELLID] = widgets.NewParagraph()
	a.elements[KEY5GCELLID].(*widgets.Paragraph).Title = " CellID "
	a.elements[KEY5GCELLID].(*widgets.Paragraph).Text = "N/A"

	a.elements[PLOTLTE] = widgets.NewPlot()
	a.elements[PLOTLTE].(*widgets.Plot).Title = " LTE "
	a.elements[PLOTLTE].(*widgets.Plot).Data = make([][]float64, 3)
	a.elements[PLOTLTE].(*widgets.Plot).AxesColor = ui.ColorWhite
	a.elements[PLOTLTE].(*widgets.Plot).LineColors[0] = ui.ColorRed
	a.elements[PLOTLTE].(*widgets.Plot).LineColors[1] = ui.ColorBlue
	a.elements[PLOTLTE].(*widgets.Plot).LineColors[2] = ui.ColorYellow
	a.elements[PLOTLTE].(*widgets.Plot).Marker = widgets.MarkerDot

	a.elements[KEYLTESNR] = widgets.NewParagraph()
	a.elements[KEYLTESNR].(*widgets.Paragraph).Text = "SNR"
	a.elements[KEYLTESNR].(*widgets.Paragraph).TextStyle.Fg = ui.ColorRed
	a.elements[KEYLTERSRP] = widgets.NewParagraph()
	a.elements[KEYLTERSRP].(*widgets.Paragraph).Text = "RSRP"
	a.elements[KEYLTERSRP].(*widgets.Paragraph).TextStyle.Fg = ui.ColorBlue
	a.elements[KEYLTERSRQ] = widgets.NewParagraph()
	a.elements[KEYLTERSRQ].(*widgets.Paragraph).Text = "RSRQ"
	a.elements[KEYLTERSRQ].(*widgets.Paragraph).TextStyle.Fg = ui.ColorYellow
	a.elements[KEYLTEBAND] = widgets.NewParagraph()
	a.elements[KEYLTEBAND].(*widgets.Paragraph).Title = " Band "
	a.elements[KEYLTEBAND].(*widgets.Paragraph).Text = "N/A"
	a.elements[KEYLTECELLID] = widgets.NewParagraph()
	a.elements[KEYLTECELLID].(*widgets.Paragraph).Title = " CellID "
	a.elements[KEYLTECELLID].(*widgets.Paragraph).Text = "N/A"

	a.elements[ALERTS] = widgets.NewParagraph()
	a.elements[ALERTS].(*widgets.Paragraph).Title = " Alerts "
	a.elements[ALERTS].(*widgets.Paragraph).Text = ""
	a.elements[ALERTS].(*widgets.Paragraph).TextStyle.Fg = ui.ColorRed

	a.grid = ui.NewGrid()
	a.grid.SetRect(0, 0, a.tw, a.th)

	a.grid.Set(
		ui.NewRow(1.0/8,
			ui.NewCol(1.0/2, a.elements[LHEADER]),
			ui.NewCol(1.0/2, a.elements[RHEADER]),
		),
		ui.NewRow(1.0/3,
			ui.NewCol(3.0/4, a.elements[PLOT5G]),
			ui.NewCol(1.0/4,
				ui.NewRow(1.0/5, a.elements[KEY5GSNR]),
				ui.NewRow(1.0/5, a.elements[KEY5GRSRP]),
				ui.NewRow(1.0/5, a.elements[KEY5GRSRQ]),
				ui.NewRow(1.0/5, a.elements[KEY5GBAND]),
				ui.NewRow(1.0/5, a.elements[KEY5GCELLID]),
			),
		),
		ui.NewRow(1.0/3,
			ui.NewCol(3.0/4, a.elements[PLOTLTE]),
			ui.NewCol(1.0/4,
				ui.NewRow(1.0/5, a.elements[KEYLTESNR]),
				ui.NewRow(1.0/5, a.elements[KEYLTERSRP]),
				ui.NewRow(1.0/5, a.elements[KEYLTERSRQ]),
				ui.NewRow(1.0/5, a.elements[KEYLTEBAND]),
				ui.NewRow(1.0/5, a.elements[KEYLTECELLID]),
			),
		),
		ui.NewRow(1.0/8, a.elements[ALERTS]),
	)

}

// Run runs the UI with polling
func (a *AlignmentUI) Run() error {

	defer ui.Close()

	a.Draw(&models.FastmileReturn{}, 0)

	go a.client.Start(a.data, a.cancel)

	tick := time.NewTicker(time.Second).C
	tickCount := 0

	for {
		select {
		case e := <-a.events:
			switch e.ID {
			case "q", "<C-c>":
				close(a.cancel)
				close(a.data)
				return nil
			}
		case <-tick:
			tickCount += 1
			if tickCount >= len(runningSteps) {
				tickCount = 0
			}
			a.Draw(&models.FastmileReturn{}, tickCount)

		case d := <-a.data:
			a.Draw(d, tickCount)
		}
	}

}

func (a *AlignmentUI) Draw(data *models.FastmileReturn, count int) {
	a.elements[RHEADER].(*widgets.Paragraph).Text = fmt.Sprintf("PRESS q TO QUIT %s", runningSteps[count])

	if data.Error != nil {
		a.elements[ALERTS].(*widgets.Paragraph).TextStyle.Fg = ui.ColorRed
		a.elements[ALERTS].(*widgets.Paragraph).Text = data.Error.Error()
		ui.Render(a.grid)
		return
	}

	if data.Body != nil {
		a.elements[ALERTS].(*widgets.Paragraph).TextStyle.Fg = ui.ColorGreen
		a.elements[ALERTS].(*widgets.Paragraph).Text = "Received data..."

		max5G := (a.elements[PLOT5G].(*widgets.Plot).Max.X / 5) * 4

		stat5G := data.Stat5G()
		a.HandleStat(stat5G.SNRQuality(0, 1), KEY5GSNR, max5G)
		a.elements[KEY5GSNR].(*widgets.Paragraph).Text = fmt.Sprintf("SNR (peak: %f)", a.stats[KEY5GSNR].max)
		a.HandleStat(stat5G.RSRPQuality(0, 1), KEY5GRSRP, max5G)
		a.elements[KEY5GRSRP].(*widgets.Paragraph).Text = fmt.Sprintf("RSRP (peak: %f)", a.stats[KEY5GRSRP].max)
		a.HandleStat(stat5G.RSRQQuality(0, 1), KEY5GRSRQ, max5G)
		a.elements[KEY5GRSRQ].(*widgets.Paragraph).Text = fmt.Sprintf("RSRQ (peak: %f)", a.stats[KEY5GRSRQ].max)

		a.elements[KEY5GBAND].(*widgets.Paragraph).Text = fmt.Sprintf("%s - %vGHz", stat5G.Band, bandmap[stat5G.Band])
		if !(a.silent) {
			a.elements[KEY5GCELLID].(*widgets.Paragraph).Text = stat5G.PhysicalCellID
		}

		maxLTE := (a.elements[PLOTLTE].(*widgets.Plot).Max.X / 5) * 4

		statLTE := data.StatLTE()
		a.HandleStat(statLTE.SNRQuality(0, 1), KEYLTESNR, maxLTE)
		a.elements[KEYLTESNR].(*widgets.Paragraph).Text = fmt.Sprintf("SNR (peak: %f)", a.stats[KEYLTESNR].max)
		a.HandleStat(statLTE.RSRPQuality(0, 1), KEYLTERSRP, maxLTE)
		a.elements[KEYLTERSRP].(*widgets.Paragraph).Text = fmt.Sprintf("RSRP (peak: %f)", a.stats[KEYLTERSRP].max)
		a.HandleStat(statLTE.RSRQQuality(0, 1), KEYLTERSRQ, maxLTE)
		a.elements[KEYLTERSRQ].(*widgets.Paragraph).Text = fmt.Sprintf("RSRQ (peak: %f)", a.stats[KEYLTERSRQ].max)

		a.elements[KEYLTEBAND].(*widgets.Paragraph).Text = fmt.Sprintf("%s - %vGHz", statLTE.Band, bandmap[statLTE.Band])
		if !(a.silent) {
			a.elements[KEYLTECELLID].(*widgets.Paragraph).Text = statLTE.PhysicalCellID
		}

		// Update plots
		a.elements[PLOT5G].(*widgets.Plot).Data[0] = a.stats[KEY5GSNR].slice
		a.elements[PLOT5G].(*widgets.Plot).Data[1] = a.stats[KEY5GRSRP].slice
		a.elements[PLOT5G].(*widgets.Plot).Data[2] = a.stats[KEY5GRSRQ].slice

		a.elements[PLOTLTE].(*widgets.Plot).Data[0] = a.stats[KEYLTESNR].slice
		a.elements[PLOTLTE].(*widgets.Plot).Data[1] = a.stats[KEYLTERSRP].slice
		a.elements[PLOTLTE].(*widgets.Plot).Data[2] = a.stats[KEYLTERSRQ].slice
	}

	ui.Render(a.grid)
}

func (a *AlignmentUI) HandleStat(val float64, key string, max int) {
	if val > a.stats[key].max {
		a.stats[key].max = val
	}

	a.stats[key].slice = append(a.stats[key].slice, val)
	if len(a.stats[key].slice) > max {
		a.stats[key].slice = a.stats[key].slice[1:]
	}
}

// Close cancels and closes the termui
func (a *AlignmentUI) Close() {
	ui.Close()
}
