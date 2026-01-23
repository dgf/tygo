package test

import (
	"fmt"
	"time"
)

type Result struct {
	Duration               time.Duration
	WordsPerMinute         int // WPM = (total keys pressed / 5) / duration in minutes
	AccuracyPercent        int // AP = (correct keys pressed / total keys pressed) * 100
	AdjustedWordsPerMinute int // AWPM = WPM * AP
}

func (r Result) String() string {
	return fmt.Sprintf("%s\r\nWPM %4d\r\nACC  %3d%%\r\nAWPM %3d",
		r.Duration, r.WordsPerMinute, r.AccuracyPercent, r.AdjustedWordsPerMinute)
}

func Calc(duration time.Duration, grid Grid) Result {
	totalKeysPressed := 0
	correctKeysPressed := 0

	for _, row := range grid {
		for _, cell := range row {
			if cell == nil {
				break
			}

			totalKeysPressed += len(cell.Inputs)

			for _, i := range cell.Inputs {
				if i == cell.Rune {
					correctKeysPressed++
				}
			}
		}
	}

	wpm := (float64(totalKeysPressed / 5)) / duration.Minutes()
	accuracy := float64(correctKeysPressed) / float64(totalKeysPressed)

	return Result{
		Duration:               duration,
		WordsPerMinute:         int(wpm),
		AccuracyPercent:        int(100 * accuracy),
		AdjustedWordsPerMinute: int(wpm * accuracy),
	}
}
