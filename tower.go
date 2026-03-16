// Package main contains a simple tower-defense game logic implemented
// for the console. Types include tracks, enemies, towers and basic game
// management functions.
package main

// Tower represents a defensive unit placed by the player. It embeds an
// Entity to reuse common positioning fields and adds attributes for range,
// damage, cost and attack speed. Towers target enemies within their
// Manhattan-radius and inflict damage each turn.
type Tower struct {
	Entity
	radius     int
	damage     int
	cost       int
	attackRate int
}

// NewTower constructs a Tower at the given grid coordinates. The
// returned pointer has default stat values and a symbol of "T". Its
// pathIndex is set to -1 since towers are stationary and do not follow
// the enemy path.
func NewTower(x int, y int) *Tower {
	tower := &Tower{
		radius:     3,
		damage:     1,
		cost:       100,
		attackRate: 1,
	}
	tower.position = [2]int{x, y}
	tower.symbol = "T"
	tower.pathIndex = -1 // Towers don't follow a path
	return tower
}

// getRadius returns the tower's attack radius measured in Manhattan
// distance (sum of row and column deltas).
func (tower *Tower) getRadius() int {
	return tower.radius
}

// getDamage returns the amount of health the tower removes from an
// enemy on a successful hit.
func (tower *Tower) getDamage() int {
	return tower.damage
}

// getCost returns the gold cost required to place this tower on the
// map.
func (tower *Tower) getCost() int {
	return tower.cost
}

// isEnemyInRange checks whether the provided enemy lies within the
// tower's attack radius. Distance is calculated using Manhattan metric
// (horizontal + vertical steps), matching grid movement rules.
func (tower *Tower) isEnemyInRange(enemy *Enemy) bool {
	towerPos := tower.getPosition()
	enemyPos := enemy.getPosition()

	// Manhattan distance
	distance := abs(towerPos[0]-enemyPos[0]) + abs(towerPos[1]-enemyPos[1])
	return distance <= tower.radius
}

// attackEnemy applies damage to the enemy if it is currently within
// range. This method simply deducts the tower's damage from the enemy's
// health via the enemy.takeDamage helper.
func (tower *Tower) attackEnemy(enemy *Enemy) {
	if tower.isEnemyInRange(enemy) {
		damage := tower.damage
		enemy.takeDamage(&damage)
	}
}

// abs returns the absolute value of an integer. This small helper
// is used for the Manhattan distance calculation.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
