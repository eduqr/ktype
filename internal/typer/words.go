package typer

import "github.com/eduqr/ktype/internal/ui"

func isPrevWordCorrect(cells []ui.Cell, cursorIndex int) bool {
	wordIndex := cells[cursorIndex-1].WordIndex
	for i := cursorIndex - 1; i >= 0 && cells[i].WordIndex == wordIndex; i-- {
		if cells[i].State != ui.StatusCorrect {
			return false
		}
	}
	return true
}

func wordHasErrors(cells []ui.Cell, cursorIndex int) bool {
	wordIndex := cells[cursorIndex].WordIndex

	for i := wordStartIndex(cells, cursorIndex); i < len(cells) && cells[i].WordIndex == wordIndex; i++ {
		if cells[i].State == ui.StatusIncorrect || cells[i].State == ui.StatusMissed {
			return true
		}
	}
	return false
}

func wordStartIndex(cells []ui.Cell, cursorIndex int) int {
	wordIndex := cells[cursorIndex].WordIndex

	start := cursorIndex
	for start > 0 && cells[start-1].WordIndex == wordIndex {
		start--
	}

	return start
}

func overflowCountForWord(cells []ui.Cell, cursorIndex int) int {
	wordIndex := cells[cursorIndex].WordIndex
	count := 0
	for i := wordStartIndex(cells, cursorIndex); i < len(cells) && cells[i].WordIndex == wordIndex; i++ {
		if cells[i].IsOverflow {
			count++
		}
	}
	return count
}

func isIndexAtWordStart(cells []ui.Cell, cursorIndex int) bool {
	return cursorIndex == 0 ||
		cells[cursorIndex-1].WordIndex != cells[cursorIndex].WordIndex
}
