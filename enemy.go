package main

import "github.com/gopxl/pixel/v2"

// EnemyType defines different enemy behaviors and strength tiers.
type EnemyType int

const (
	EnemyTypeDefault EnemyType = iota
	EnemyTypeFast
	EnemyTypeTank
)

// Enemy holds the moving enemy's current position, health, and path progress.
type Enemy struct {
	pos       pixel.Vec
	pathPoint int
	speed     float64
	health    float64
	maxHealth float64
	typeID    EnemyType
}

// NewEnemy constructs a new enemy with stats based on type.
func NewEnemy(start pixel.Vec, typeID EnemyType) *Enemy {
	var speed, health float64
	switch typeID {
	case EnemyTypeFast:
		speed = 150
		health = 3
	case EnemyTypeTank:
		speed = 50
		health = 15
	default:
		speed = 90
		health = 5
	}
	return &Enemy{pos: start, pathPoint: 0, speed: speed, health: health, maxHealth: health, typeID: typeID}
}

// Update moves the enemy towards the next path node based on dt (delta time).
func (e *Enemy) Update(dt float64, path []pixel.Vec) {
	if e.pathPoint >= len(path)-1 {
		return
	}
	target := path[e.pathPoint+1] // Path node enemy is moving towards
	dir := target.Sub(e.pos) // Get direction enemy should move
	distance := dir.Len() // Get distance 
	if distance < 1 { // If enemy is close enough, make them start heading to the next point
		e.pathPoint++
		if e.pathPoint < len(path) {
			e.pos = path[e.pathPoint] // Set position to the point they just reached
		}
		return
	}
	move := dir.Unit().Scaled(e.speed * dt) // Moved in the direction by the amount of speed times the amount time since the last frame
	if move.Len() >= distance { // If enemy would move over point, assign their position to the point and set next path point
		e.pos = target
		e.pathPoint++
	} else { // Move enemy
		e.pos = e.pos.Add(move)
	}
}

// ReachedEnd returns true when the enemy has reached the final node in path.
func (e *Enemy) ReachedEnd(path []pixel.Vec) bool {
	return e.pathPoint >= len(path)-1
}

// InRange returns whether the enemy is within radius of a point.
func (e *Enemy) InRange(point pixel.Vec, radius float64) bool {
	return e.pos.Sub(point).Len() <= radius
}

// TakeDamage subtracts hit points from the enemy.
func (e *Enemy) TakeDamage(amount float64) {
	e.health -= amount
}

// IsDead reports whether the enemy has no health remaining.
func (e *Enemy) IsDead() bool {
	return e.health <= 0
}
