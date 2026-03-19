package main

import "github.com/gopxl/pixel/v2"

// Entity is a base object that tracks a position in world coordinates.
type Entity struct {
	pos pixel.Vec
}

// Position returns the Entity's current position.
func (e *Entity) Position() pixel.Vec {
	return e.pos
}

// SetPosition updates the Entity's position value.
func (e *Entity) SetPosition(p pixel.Vec) {
	e.pos = p
}
