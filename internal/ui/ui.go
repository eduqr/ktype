package ui

import (
	"fmt"
	"strings"

	"github.com/eduqr/ktype/internal/stats"
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
	"github.com/guptarohit/asciigraph"
)

var finalTitleLines = []string{
	" _    _                    ",
	"| | _| |_ _   _ _ __   ___ ",
	"| |/ / __| | | | '_ \\ / _ \\",
	"|   <| |_| |_| | |_) |  __/",
	"|_|\\_\\\\__|\\__, | .__/ \\___|",
	"          |___/|_|         ",
}

var finalOptionLines = []string{
	"[r] repeat",
	"[tab] new challenge",
	"[esc] quit",
}

func InitScreen() (tcell.Screen, error) {
	defStyle := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)

	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	if err := screen.Init(); err != nil {
		return nil, err
	}

	screen.SetStyle(defStyle)
	screen.DisablePaste()
	screen.SetTitle("ktype")
	screen.Clear()

	return screen, nil
}

func DrawChallenge(screen tcell.Screen, cells []Cell, styles Styles) {
	screen.Clear()

	for _, cell := range cells {
		DrawCell(screen, cell, styles)
	}
}

func DrawCell(screen tcell.Screen, cell Cell, styles Styles) {
	screen.Put(cell.X, cell.Y, cell.Expected, styles.CellStyle(cell.State))
}

func DrawCursorCell(screen tcell.Screen, cell Cell, styles Styles) {
	screen.Put(cell.X, cell.Y, cell.Expected, styles.CursorStyle(cell.State))
}

func ShowFinalScreen(screen tcell.Screen, styles Styles, stats stats.Stats, samples []float64) {
	screen.Clear()
	screen.HideCursor()

	layout := NewLayout(screen)
	y := layout.Top

	y = drawTitle(screen, layout, y, styles.Normal)
	y += 2
	y = drawResults(screen, layout, y, stats, styles)
	y += 2
	y = showWPMSamples(screen, layout, y, styles, samples)
	y += 2
	y = showOptions(screen, layout, y, styles)

	screen.Show()
}

func ShowAbortRoundScreen(screen tcell.Screen, styles Styles) {
	screen.Clear()
	screen.HideCursor()

	layout := NewLayout(screen)
	y := layout.Top

	y = drawTitle(screen, layout, y, styles.Normal)
	y += 2
	y = showAbortRoundMessage(screen, layout, y, styles)
	y += 2
	y = showOptions(screen, layout, y, styles)

	screen.Show()
}

func drawResults(screen tcell.Screen, layout Layout, y int, stats stats.Stats, styles Styles) int {
	result := fmt.Sprintf("Acc: %.2f, WPM: %.2f, Raw WPM: %.2f, Time: %.2fs", stats.Accuracy, stats.WPM, stats.RawWPM, stats.TimeElapsed.Seconds())

	return drawLines(layout, screen, []string{result}, y, styles.Normal)
}

func showOptions(screen tcell.Screen, layout Layout, y int, styles Styles) int {
	if layout.ContentWidth < 60 {
		return drawLinesLeft(layout, screen, finalOptionLines, y, styles.Normal)
	}

	options := strings.Join(finalOptionLines, "     ")
	return drawLines(layout, screen, []string{options}, y, styles.Normal)
}

func drawTitle(screen tcell.Screen, layout Layout, y int, style tcell.Style) int {
	return drawLines(layout, screen, finalTitleLines, y, style)
}

func showAbortRoundMessage(screen tcell.Screen, layout Layout, y int, styles Styles) int {
	return drawLines(layout, screen, []string{"Challenge aborted"}, y, styles.Incorrect)
}

func showWPMSamples(screen tcell.Screen, layout Layout, y int, styles Styles, samples []float64) int {
	if len(samples) < 2 {
		return drawLines(layout, screen, []string{"Not enough data to display the graph"}, y, styles.Normal)
	}

	tickCount := 6
	if len(samples) <= 10 {
		tickCount = len(samples)
	}

	graph := asciigraph.Plot(samples,
		asciigraph.Height(min(4, len(samples))),
		asciigraph.Width(max(layout.ContentWidth-6, 10)),
		asciigraph.XAxisRange(1, float64(len(samples))),
		asciigraph.XAxisTickCount(tickCount),
		asciigraph.XAxisValueFormatter(func(v float64) string {
			return fmt.Sprintf("%.0f", v)
		}),
		asciigraph.Precision(0),
	)

	return drawLinesLeft(layout, screen, strings.Split(graph, "\n"), y, styles.Normal)
}
