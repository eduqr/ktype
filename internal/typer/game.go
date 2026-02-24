package typer

import (
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

func StartGame(screen tcell.Screen, challenge []string) {
	for {
		screen.Show()

		event := <-screen.EventQ()

		switch ev := event.(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			}
			if ev.Key() == tcell.KeyRune {
				screen.Put(10, 10, ev.Str(), tcell.StyleDefault)
				if ev.Str() == "a" {
					correctStyle := tcell.StyleDefault.Foreground(color.Green)
					screen.PutStrStyled(10, 11, "es una a", correctStyle)
				} else {
					// borrar el mensaje anterior
					for i := 0; i < len("es una a"); i++ {
						screen.Put(10+i, 11, " ", tcell.StyleDefault)
					}
				}
			}
		}
	}
}
