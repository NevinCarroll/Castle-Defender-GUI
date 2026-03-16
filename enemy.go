// Package main defines core types for a console tower defense game,
// including Enemy which moves along a predefined path and can take damage.
package main

// Enemy represents a hostile unit advancing along the track. It embeds
// Entity for position and path tracking and adds a health counter that
// towers will decrement when attacking.
//
// Enemies increment their pathIndex on each move and are removed when
// they reach the end or their health drops to zero.
type Enemy struct {
	Entity
	health int
}

// NewEnemy constructs a fresh Enemy with default health and symbol.
// The pathIndex is initialized to 0 so it begins at the start of the
// track when spawned.
func NewEnemy() *Enemy {
	enemy := &Enemy{
		health: 5, 
	}
	enemy.symbol = "E"
	enemy.pathIndex = 0
	return enemy
}

// getHealth returns the current hit points of the enemy.
func (enemy *Enemy) getHealth() int {
	return enemy.health
}

// takeDamage subtracts the provided value from the enemy's health.
// The damage is passed by pointer to match existing call sites.
// 
func (enemy *Enemy) takeDamage(damage *int) {
	enemy.health -= *damage
}

// move advances the enemy one step along the track path. If there is
// no path or the enemy is already at the last cell, the position remains
// unchanged.
func (enemy *Enemy) move(track *Track) {
	path := track.getPath()

	if len(path) == 0 {
		return
	}

	// Move to next position in path (one step per turn in turn-based)
	if enemy.pathIndex < len(path)-1 {
		enemy.pathIndex++
		nextPos := path[enemy.pathIndex]
		enemy.position = nextPos
	}
}

// isAtEnd returns true if the enemy has reached or passed the final
// cell of the track path. Used to determine life loss when enemies
// escape.
func (enemy *Enemy) isAtEnd(track *Track) bool {
	path := track.getPath()
	if len(path) == 0 {
		return false
	}
	return enemy.pathIndex >= len(path)-1
}
