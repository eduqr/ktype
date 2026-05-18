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

	left := (width - contentWidth) / 2

	if left < 0 {
		left = 0
	}

	top := height / 6

	return Layout{
		ScreenWidth:  width,
		ScreenHeight: height,
		ContentWidth: contentWidth,
		Top:          top,
		Left:         left,
	}
}
