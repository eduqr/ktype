package typer

import "github.com/gdamore/tcell/v3"

func handleResultsKey(ev *tcell.EventKey) Action {
	if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
		return ActionQuit
	}

	if ev.Key() == tcell.KeyTab {
		return ActionNewChallenge
	}

	if ev.Key() == tcell.KeyRune {
		input := ev.Str()
		if input == "r" || input == "R" {
			return ActionRepeat
		}
	}

	return ActionNone
}

func WaitResultsAction(screen tcell.Screen, redraw func()) Action {
	for {
		event := <-screen.EventQ()

		switch ev := event.(type) {
		case *tcell.EventKey:
			action := handleResultsKey(ev)
			if action != ActionNone {
				return action
			}
		case *tcell.EventResize:
			screen.Sync()
			redraw()
		}
	}
}
