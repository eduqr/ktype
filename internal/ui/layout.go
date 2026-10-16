package ui

import "github.com/gdamore/tcell/v3"

const (
	MinScreenWidth  = 40
	MinScreenHeight = 10
)

type Layout struct {
	ScreenWidth  int
	ScreenHeight int
	ContentWidth int
	Top          int
	Left         int
}

func NewLayout(screen tcell.Screen) Layout {
	width, height := screen.Size()

	contentWidth := min(width-4, 120)
	if contentWidth < 1 {
		contentWidth = 1
	}

	left := (width - contentWidth) / 2

	if left < 0 {
		left = 0
	}

	top := height / 10

	return Layout{
		ScreenWidth:  width,
		ScreenHeight: height,
		ContentWidth: contentWidth,
		Top:          top,
		Left:         left,
	}
}

func (layout Layout) IsSupported() bool {
	return layout.ScreenWidth >= MinScreenWidth && layout.ScreenHeight >= MinScreenHeight
}

func (layout Layout) useShortTitle() bool {
	return layout.ScreenHeight < 25
}

func (layout Layout) useStackedOptions() bool {
	return layout.ContentWidth < 60
}

func (layout Layout) canShowGraph(y int, reservedLines int) bool {
	const minGraphWidth = 36
	const minGraphLines = 7

	availableLines := layout.ScreenHeight - y - reservedLines
	return layout.ContentWidth >= minGraphWidth && availableLines >= minGraphLines
}

func (layout Layout) optionsLineCount() int {
	if layout.useStackedOptions() {
		return 3
	}
	return 1
}

func (layout Layout) gap() int {
	if layout.ScreenHeight < 20 {
		return 1
	}
	return 2
}

func (layout Layout) hasRoom(y int, lines int) bool {
	return y+lines <= layout.ScreenHeight
}

func textWidth(s string) int {
	return len([]rune(s))
}

func centerX(layout Layout, text string) int {
	return layout.Left + max((layout.ContentWidth-textWidth(text))/2, 0)
}

func drawLines(layout Layout, screen tcell.Screen, lines []string, y int, style tcell.Style) int {
	for _, line := range lines {
		x := centerX(layout, line)
		screen.PutStrStyled(x, y, line, style)
		y++
	}

	return y
}

func drawLinesLeft(layout Layout, screen tcell.Screen, lines []string, y int, style tcell.Style) int {
	for _, line := range lines {
		screen.PutStrStyled(layout.Left, y, line, style)
		y++
	}

	return y
}
