package ui

import "github.com/gdamore/tcell/v3"

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
