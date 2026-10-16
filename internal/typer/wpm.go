package typer

import (
	"time"

	"github.com/eduqr/ktype/internal/stats"
)

func (gs *GameState) startWPMSampler() {
	if gs.wpmTicker != nil {
		return
	}

	gs.wpmTicker = time.NewTicker(1 * time.Second)
}

func (gs *GameState) sampleWPM() {
	correctCells := gs.countCorrectCells()
	elapsed := time.Since(gs.StartTime).Seconds()
	sample := stats.CalculateWPM(correctCells, elapsed)
	gs.WPMSamples = append(gs.WPMSamples, sample)
}

func (gs *GameState) stopWPMSampler() {
	if gs.wpmTicker == nil {
		return
	}

	gs.wpmTicker.Stop()
	gs.wpmTicker = nil
}
