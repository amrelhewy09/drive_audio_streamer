package gui

import (
	"fmt"
	"log"
	"os/exec"
	"syscall"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type Gui struct {
	*widgets.Gauge
	*widgets.Paragraph
	Name string
}

func CreatePlayerGui(itemName string) *Gui {
	return &Gui{
		Gauge:     widgets.NewGauge(),
		Paragraph: widgets.NewParagraph(),
		Name:      itemName,
	}
}

func (r *Gui) RenderMediaPlayerUI(terminateDownloadChan chan int, terminateExecutionChan chan int, cmd *exec.Cmd) {
	if err := ui.Init(); err != nil {
		log.Fatalf("Failed to initialize termui: %v", err)
	}
	r.Paragraph.Text = fmt.Sprintf("Playing: %s", r.Name)
	r.Paragraph.SetRect(0, 0, 80, 3)
	r.Paragraph.BorderStyle.Fg = ui.ColorYellow
	r.Gauge.Title = "Streamable"
	r.Gauge.SetRect(0, 3, 80, 6)
	r.Gauge.Percent = 0
	r.Gauge.BarColor = ui.ColorGreen
	r.Gauge.BorderStyle.Fg = ui.ColorWhite
	r.Gauge.TitleStyle.Fg = ui.ColorCyan
	p := widgets.NewParagraph()
	p.SetRect(0, 6, 80, 9)
	p.Text = "Press q or Ctrl+C to exit and return to your list"

	ui.Render(r.Paragraph, r.Gauge, p)

	r.ListenForEvents(terminateDownloadChan, terminateExecutionChan, cmd)
}

func (r *Gui) UpdateProgress(percent int) {
	if percent >= 100 {
		percent = 100
	}

	r.Gauge.Percent = percent
	ui.Render(r.Gauge)
}

func (r *Gui) Close() {
	ui.Close()
}

func (r *Gui) ListenForEvents(terminateDownloadChan chan int, terminateExecutionChan chan int, cmd *exec.Cmd) {
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			cmd.Process.Signal(syscall.SIGTERM)
			close(terminateExecutionChan)
			close(terminateDownloadChan)
			r.Close()
			return
		}
	}
}
