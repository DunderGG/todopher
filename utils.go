// utils.go contains miscellaneous helper functions and terminal UI components.
// It includes the progress bar, ASCII banner, and time conversion utilities.
package main

import (
	"fmt"
	"time"
)

// printProgressBar renders a simple ASCII progress bar in the terminal.
func printProgressBar(current, total int) {
	width := 40
	percent := float64(current) / float64(total)
	filled := int(percent * float64(width))

	bar := "["
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "="
		} else if i == filled {
			bar += ">"
		} else {
			bar += " "
		}
	}
	bar += "]"

	fmt.Printf("\rScanning: %s %d/%d (%d%%)", bar, current, total, int(percent*100))
}

// printIntro displays an ASCII art banner.
func printIntro() {
	fmt.Println("======================================================")
	intro := `
   ______ ___  ____  ____  ____  __ __  ______ ____
  /_  __/ __ \/ __ \/ __ \/ __ \/ // / / ____/ __  \
   / / / / / / / / / / / / /_/ / _  / / __/ / /_/  /
  / / / /_/ / /_/ / /_/ / ____/ // / / /___/__/ ,_/ 
 /_/  \____/_____/\____/_/   /_//_/ /_____/_/ |_|
`
	fmt.Println(intro)
	fmt.Println("======================================================")
}

// stringsToDate parses a YYYY-MM-DD string into a time.Time object.
func stringsToDate(s string) (time.Time, error) {
	return time.Parse("2006-01-02", s)
}

// timeSinceInDays calculates the number of days between now and the provided time.
func timeSinceInDays(t time.Time) float64 {
	return time.Since(t).Hours() / 24
}
