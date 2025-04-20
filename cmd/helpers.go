package cmd

func calculateColumnWidths(totalWidth int, ratios []float64) []int {
	widths := make([]int, len(ratios))
	padding := 3 * (len(ratios) - 1) // space between columns (" | ")
	available := totalWidth - padding

	for i, r := range ratios {
		widths[i] = int(float64(available) * r)
	}

	// Fix rounding errors to ensure full width is used
	sum := 0
	for _, w := range widths {
		sum += w
	}
	widths[len(widths)-1] += available - sum

	return widths
}
