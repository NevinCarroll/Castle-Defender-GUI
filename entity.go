package main

import "github.com/gopxl/pixel/v2"

type Entity struct {
	pos pixel.Vec
}

func (e *Entity) Position() pixel.Vec {
	return e.pos
}

func (e *Entity) SetPosition(p pixel.Vec) {
	e.pos = p
}
