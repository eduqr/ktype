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

var shortFinalTitle = []string{
	"ktype",
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
		screen.Fini()
		return nil, err
	}

	screen.SetStyle(defStyle)
	screen.DisablePaste()
	screen.EnableMouse(tcell.MouseButtonEvents)
	screen.SetTitle("ktype")
	screen.Clear()

	return screen, nil
}

func drawChallenge(screen tcell.Screen, cells []Cell, styles Styles) {
	screen.Clear()
	layout := NewLayout(screen)
	if !layout.IsSupported() {
		showMinimumScreenSizeMessage(screen, layout, styles)
		screen.Show()
		return
	}

	for _, cell := range cells {
		drawCell(screen, cell, styles)
	}
}

func DrawGame(screen tcell.Screen, cells []Cell, styles Styles, layout Layout, cursorIndex int) {
	recalcChallengePositions(cells, layout)
	drawChallenge(screen, cells, styles)
	if layout.IsSupported() {
		drawCursor(screen, cells[cursorIndex])
	} else {
		screen.HideCursor()
	}
	screen.Show()
}

func drawCursor(screen tcell.Screen, cell Cell) {
	screen.ShowCursor(cell.X, cell.Y)
}

func drawCell(screen tcell.Screen, cell Cell, styles Styles) {
	text := cell.Expected
	if cell.IsOverflow {
		text = cell.OverflowChar
	}

	if cell.IsTerminator {
		return
	}

	screen.Put(cell.X, cell.Y, text, styles.cellStyle(cell.State))
}

func ShowFinalScreen(screen tcell.Screen, styles Styles, stats stats.Stats, samples []float64) {
	screen.Clear()
	screen.HideCursor()

	layout := NewLayout(screen)
	if !layout.IsSupported() {
		showMinimumScreenSizeMessage(screen, layout, styles)
		screen.Show()
		return
	}

	y := layout.Top
	gap := layout.gap()

	y = drawTitle(screen, layout, y, styles.Title)
	y += gap
	y = drawResults(screen, layout, y, stats, styles)
	y += gap
	y = showWPMSamples(screen, layout, y, styles, samples)
	y += gap
	y = showOptions(screen, layout, y, styles)

	screen.Show()
}

func ShowAbortRoundScreen(screen tcell.Screen, styles Styles) {
	screen.Clear()
	screen.HideCursor()

	layout := NewLayout(screen)
	if !layout.IsSupported() {
		showMinimumScreenSizeMessage(screen, layout, styles)
		screen.Show()
		return
	}

	y := layout.Top
	gap := layout.gap()

	y = drawTitle(screen, layout, y, styles.Title)
	y += gap
	y = showAbortRoundMessage(screen, layout, y, styles)
	y += gap
	y = showOptions(screen, layout, y, styles)

	screen.Show()
}

func drawResults(screen tcell.Screen, layout Layout, y int, stats stats.Stats, styles Styles) int {
	const separator = "  "

	results := []string{
		fmt.Sprintf("Acc: %.2f", stats.Accuracy),
		fmt.Sprintf("WPM: %.2f", stats.WPM),
		fmt.Sprintf("Raw WPM: %.2f", stats.RawWPM),
		fmt.Sprintf("Time: %.2fs", stats.TimeElapsed.Seconds()),
	}
	resultStyles := []tcell.Style{styles.Accuracy, styles.WPM, styles.RawWPM, styles.Time}
	result := strings.Join(results, separator)

	if textWidth(result) <= layout.ContentWidth {
		x := centerX(layout, result)
		for i, result := range results {
			if i > 0 {
				x += textWidth(separator)
			}
			screen.PutStrStyled(x, y, result, resultStyles[i])
			x += textWidth(result)
		}
		return y + 1
	}

	for i, result := range results {
		screen.PutStrStyled(layout.Left, y+i, result, resultStyles[i])
	}

	return y + len(results)
}

func showOptions(screen tcell.Screen, layout Layout, y int, styles Styles) int {
	if layout.useStackedOptions() {
		return drawLinesLeft(layout, screen, finalOptionLines, y, styles.Actions)
	}

	options := strings.Join(finalOptionLines, "     ")
	return drawLines(layout, screen, []string{options}, y, styles.Actions)
}

func drawTitle(screen tcell.Screen, layout Layout, y int, style tcell.Style) int {
	if layout.useShortTitle() {
		return drawLines(layout, screen, shortFinalTitle, y, style)
	}
	return drawLines(layout, screen, finalTitleLines, y, style)
}

func showAbortRoundMessage(screen tcell.Screen, layout Layout, y int, styles Styles) int {
	return drawLines(layout, screen, []string{"Challenge aborted"}, y, styles.Error)
}

func showWPMSamples(screen tcell.Screen, layout Layout, y int, styles Styles, samples []float64) int {
	reservedLines := layout.optionsLineCount() + layout.gap()

	if len(samples) < 2 {
		if !layout.hasRoom(y, reservedLines+1) {
			return y
		}
		return drawLines(layout, screen, []string{"Not enough data to display the graph"}, y, styles.Error)
	}

	if !layout.canShowGraph(y, reservedLines) {
		if !layout.hasRoom(y, reservedLines+1) {
			return y
		}
		return drawLines(layout, screen, []string{"Not enough space to display the graph"}, y, styles.Error)
	}

	tickCount := min(6, len(samples))

	graphHeight := min(4, len(samples))
	minimum, maximum := samples[0], samples[0]
	for _, sample := range samples[1:] {
		minimum = min(minimum, sample)
		maximum = max(maximum, sample)
	}
	sampleRange := maximum - minimum
	if sampleRange >= 1 && int(sampleRange) < graphHeight {
		graphHeight = int(sampleRange)
	}

	graph := asciigraph.Plot(samples,
		asciigraph.Height(graphHeight),
		asciigraph.Width(max(layout.ContentWidth-6, 10)),
		asciigraph.XAxisRange(1, float64(len(samples))),
		asciigraph.XAxisTickCount(tickCount),
		asciigraph.XAxisValueFormatter(func(v float64) string {
			return fmt.Sprintf("%.0f", v)
		}),
		asciigraph.YAxisValueFormatter(func(v float64) string {
			return fmt.Sprintf("%.0f", v)
		}),
	)

	return drawLinesLeft(layout, screen, strings.Split(graph, "\n"), y, styles.Graph)
}

func showMinimumScreenSizeMessage(screen tcell.Screen, layout Layout, styles Styles) {
	message := fmt.Sprintf("Terminal too small. Resize to at least %dx%d.", MinScreenWidth, MinScreenHeight)
	if textWidth(message) > layout.ContentWidth {
		message = fmt.Sprintf("Resize to %dx%d", MinScreenWidth, MinScreenHeight)
	}
	drawLines(layout, screen, []string{message}, max(layout.ScreenHeight/2, 0), styles.Error)
}
