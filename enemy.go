package main

import "github.com/gopxl/pixel/v2"

type EnemyType int

const (
	EnemyTypeDefault EnemyType = iota
	EnemyTypeFast
	EnemyTypeTank
)

type Enemy struct {
	pos       pixel.Vec
	pathPoint int
	speed     float64
	health    float64
	maxHealth float64
	typeID    EnemyType
}

func NewEnemy(start pixel.Vec, typeID EnemyType) *Enemy {
	var speed, health float64
	switch typeID {
	case EnemyTypeFast:
		speed = 150
		health = 4
	case EnemyTypeTank:
		speed = 50
		health = 20
	default:
		speed = 90
		health = 7
	}
	return &Enemy{pos: start, pathPoint: 0, speed: speed, health: health, maxHealth: health, typeID: typeID}
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
