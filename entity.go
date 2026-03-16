// Package main provides simple game structures for a console-based
// tower defense prototype. Entities represent anything with a location
// and optional movement behavior (enemies, towers, etc.).
package main

// Entity is a generic object that occupies a grid cell and can track its
// progress along the path. It is embedded by more specific types like
// Enemy and Tower to share common functionality.
//
// Fields:
//   position    - [row, col] location on the map
//   symbol      - single-character string used when rendering
//   pathIndex   - index into the track path slice (used by movers)
//   moveCounter - simple tick counter for controlling movement speed
//
// Getters and setters provide controlled access to these fields.
type Entity struct {
	position    [2]int // Current position on the track
	symbol      string
	pathIndex   int // Current position in the path (0 = start)
	moveCounter int // Counter to control movement speed
}

// getPosition returns the current grid coordinates of the entity.
func (entity *Entity) getPosition() [2]int {
	return entity.position
}

// setPosition updates the entity's coordinates to the given values.
// Pointers are used to match the calling convention of the original
// code; nil checks are not performed.
func (entity *Entity) setPosition(x *int, y *int) {
	entity.position[0] = *x
	entity.position[1] = *y
}

// getPathIndex returns the current index of the entity within the
// track's path slice. Used primarily by moving entities such as
// enemies.
func (entity *Entity) getPathIndex() int {
	return entity.pathIndex
}

// setPathIndex sets the entity's position along the path. Negative
// values indicate the entity is not currently following the path.
func (entity *Entity) setPathIndex(index int) {
	entity.pathIndex = index
}

// getSymbol returns the character used to represent the entity when
// rendering the map.
func (entity *Entity) getSymbol() string {
	return entity.symbol
}

// setSymbol assigns a new rendering symbol to the entity.
func (entity *Entity) setSymbol(symbol string) {
	entity.symbol = symbol
}
