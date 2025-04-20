package cmd

import (
	"fmt"
)

func truncateOrPad(s string, width int) string {
	if len(s) > width {
		if width > 3 {
			return s[:width-3] + "..."
		}
		return s[:width]
	}
	return fmt.Sprintf("%-*s", width, s)
}
