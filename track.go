// Package main implements a simple console-based track
// representation and parsing for a racing or path-following game.
//
// The Track type encapsulates a fixed-size grid layout and can
// extract an ordered sequence of coordinates representing the
// path defined by "1" cells within that grid.
package main

// Track holds the map layout and the extracted path from that layout.
// layout is a 10×15 grid of strings where "1" marks path segments and
// "*" may also be treated as part of the active path. The path slice
// contains coordinates in the order they are traversed.
//
// Methods on Track allow setting the layout, retrieving the layout and
// path, and internally parsing the layout to compute the path.
type Track struct {
	layout [10][15]string
	path   [][2]int // Ordered path coordinates
}

// setLayout assigns a new grid to the track and triggers a
// re-parsing of the path contained within that grid. The provided
// layout must be exactly 10 rows by 15 columns.
func (track *Track) setLayout(layout [10][15]string) {
	track.layout = layout
	track.parsePath()
}

// getLayout returns a copy of the track's current grid layout.
// Callers should treat the returned array as read-only; modifying it
// does not affect the internal state of the Track instance.
func (track *Track) getLayout() [10][15]string {
	return track.layout
}

// getPath returns the slice of coordinates representing the parsed
// path. Each element is a two‑element array containing row and column
// indices. The slice is ordered from the start of the path to its end.
func (track *Track) getPath() [][2]int {
	return track.path
}

// parsePath examines the layout grid, locates the first cell marked
// "1" (searching from the top-left), and then follows adjacent cells
// marked "1" or "*" in a straight path. The result is stored in the
// track.path field as an ordered list of [row, column] pairs.
//
// The algorithm avoids backtracking by remembering the previous cell
// and only considers new neighbours. If no starting cell is found the
// path slice is left empty.
func (track *Track) parsePath() {
	var path [][2]int

	// Find starting position (1 in top-left area)
	var startRow, startCol int = -1, -1
	for i := 0; i < ROW_COUNT; i++ {
		for j := 0; j < COLUMN_COUNT; j++ {
			if track.layout[i][j] == "1" {
				startRow, startCol = i, j
				break
			}
		}
		if startRow != -1 {
			break
		}
	}

	if startRow == -1 {
		return // No path found
	}

	path = append(path, [2]int{startRow, startCol})

	// Follow the path
	prevRow, prevCol := -1, -1
	currentRow, currentCol := startRow, startCol

	for {
		var found bool = false
		// Check all 4 directions
		directions := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

		for _, dir := range directions {
			newRow := currentRow + dir[0]
			newCol := currentCol + dir[1]

			// Skip if going back
			if newRow == prevRow && newCol == prevCol {
				continue
			}

			// Check bounds
			if newRow < 0 || newRow >= ROW_COUNT || newCol < 0 || newCol >= COLUMN_COUNT {
				continue
			}

			cell := track.layout[newRow][newCol]
			if cell == "1" || cell == "*" {
				path = append(path, [2]int{newRow, newCol})
				prevRow, prevCol = currentRow, currentCol
				currentRow, currentCol = newRow, newCol
				found = true
				break
			}
		}

		if !found {
			break
		}
	}

	track.path = path
}
