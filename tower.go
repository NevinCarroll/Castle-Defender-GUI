package main

import "github.com/gopxl/pixel/v2"

// TowerType enumerates the available tower classes the player can build.
type TowerType int

const (
	TowerTypeStandard TowerType = iota + 1
	TowerTypeRapid
	TowerTypeSniper
)

// Tower stores per-instance placement and attack state.
type Tower struct {
	pos           pixel.Vec
	radius        float64
	damage        float64
	attackCadence float64
	cooldown      float64
	typeID        TowerType
}

// TowerConfig is the static configuration for each TowerType.
type TowerConfig struct {
	Radius        float64
	Damage        float64
	AttackCadence float64
	Cost          int
	Label         string
}

// TowerConfigs maps tower types to their static gameplay values for radius, damage, cadence, and cost.
var TowerConfigs = map[TowerType]TowerConfig{
	TowerTypeStandard: {Radius: 120, Damage: 1.0, AttackCadence: 0.16, Cost: 100, Label: "Standard"},
	TowerTypeRapid:    {Radius: 90, Damage: 0.5, AttackCadence: 0.08, Cost: 100, Label: "Rapid"},
	TowerTypeSniper:   {Radius: 180, Damage: 2.5, AttackCadence: 0.35, Cost: 100, Label: "Sniper"},
}

// NewTower creates a new tower instance for the given type and position.
func NewTower(pos pixel.Vec, typeID TowerType) *Tower {
	cfg, ok := TowerConfigs[typeID] // Gets tower type and its stats
	if !ok {
		cfg = TowerConfigs[TowerTypeStandard] 
		typeID = TowerTypeStandard
	}
	// Returns the new tower type
	return &Tower{pos: pos, radius: cfg.Radius, damage: cfg.Damage, attackCadence: cfg.AttackCadence, typeID: typeID}
}

// Update checks cooldown, finds the closest target in range, and applies damage once ready.
// Returns the enemy that was attacked (or nil if no attack occurred).
func (t *Tower) Update(dt float64, enemies []*Enemy) *Enemy {
	t.cooldown -= dt
	// If on cooldown, don't attack
	if t.cooldown > 0 {
		return nil
	}

	var target *Enemy
	closest := 1e9
	for _, e := range enemies { // Checkes all enemies in range
		if e.IsDead() {
			continue
		}
		d := e.pos.Sub(t.pos).Len() // Take enemy position, subtract by towers position, to figure out distance between the two
		if d <= t.radius && d < closest {
			target = e
			closest = d
		}
	}

	// Damage enemy
	if target != nil {
		target.TakeDamage(t.damage)
		t.cooldown = t.attackCadence
		return target
	}
	return nil
}
