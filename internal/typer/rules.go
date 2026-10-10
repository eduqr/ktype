package typer

import "github.com/eduqr/ktype/internal/config"

type Rules struct {
	StopOn          config.StopOn
	MaxWordOverflow int
}
