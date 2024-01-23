package utils

import "strings"

func ParseGrid(input string) [][]bool {
	lines := strings.Split(strings.Trim(input, " "), "\n")
	grid := make([][]bool, len(lines))

	for i, line := range lines {
		row := make([]bool, len(line))
		for j, c := range line {
			row[j] = c == '#'
		}
		grid[i] = row
	}

	return grid
}
