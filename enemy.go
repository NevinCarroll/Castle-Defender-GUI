package main

import "github.com/gopxl/pixel/v2"

type Enemy struct {
	pos       pixel.Vec
	pathPoint int
	speed     float64
	health    float64
	maxHealth float64
}

func NewEnemy(start pixel.Vec, speed float64, health float64) *Enemy {
	return &Enemy{pos: start, pathPoint: 0, speed: speed, health: health, maxHealth: health}
}

func (e *Enemy) Update(dt float64, path []pixel.Vec) {
	if e.pathPoint >= len(path)-1 {
		return
	}
	target := path[e.pathPoint+1]
	dir := target.Sub(e.pos)
	distance := dir.Len()
	if distance < 1 {
		e.pathPoint++
		if e.pathPoint < len(path) {
			e.pos = path[e.pathPoint]
		}
		return
	}
	move := dir.Unit().Scaled(e.speed * dt)
	if move.Len() >= distance {
		e.pos = target
		e.pathPoint++
	} else {
		e.pos = e.pos.Add(move)
	}
}

func (e *Enemy) ReachedEnd(path []pixel.Vec) bool {
	return e.pathPoint >= len(path)-1
}

func (e *Enemy) InRange(point pixel.Vec, radius float64) bool {
	return e.pos.Sub(point).Len() <= radius
}

func (e *Enemy) TakeDamage(amount float64) {
	e.health -= amount
}

func (e *Enemy) IsDead() bool {
	return e.health <= 0
}
