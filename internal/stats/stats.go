package stats

import (
	"time"
)

type Stats struct {
	Accuracy    float64
	WPM         float64
	RawWPM      float64
	TimeElapsed time.Duration
}

func calculateAccuracy(correctCells, errors int) float64 {
	typed := correctCells + errors

	if typed == 0 {
		return 100
	}

	acc := float64(correctCells) / float64(typed) * 100

	return acc
}

func CalculateWPM(correctCells int, elapsedSeconds float64) float64 {
	if elapsedSeconds <= 0 {
		return 0
	}

	return float64(correctCells) * (60 / elapsedSeconds) / 5
}

func calculateRawWPM(correctCells, errors int, elapsedSeconds float64) float64 {
	if elapsedSeconds <= 0 {
		return 0
	}

	return float64(correctCells+errors) * (60 / elapsedSeconds) / 5
}

func calculateTimeElapsed(start, end time.Time) time.Duration {
	if start.IsZero() || end.IsZero() {
		return 0
	}

	if end.Before(start) {
		return 0
	}

	return end.Sub(start)
}

func ComputeStats(correctCells, errors int, start, end time.Time) (stats Stats) {
	stats.TimeElapsed = calculateTimeElapsed(start, end)
	stats.Accuracy = calculateAccuracy(correctCells, errors)
	stats.WPM = CalculateWPM(correctCells, stats.TimeElapsed.Seconds())
	stats.RawWPM = calculateRawWPM(correctCells, errors, stats.TimeElapsed.Seconds())
	return
}
