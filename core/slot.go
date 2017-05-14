package core

import "time"

func GetTime() float64 {
	mainNetStart := time.Date(2017, 3, 21, 13, 00, 0, 0, time.UTC)
	now := time.Now()

	diff := now.Sub(mainNetStart)
	return diff.Seconds()
}
