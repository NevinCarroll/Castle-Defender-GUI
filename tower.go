package main

import "github.com/gopxl/pixel/v2"

type Tower struct {
	pos           pixel.Vec
	radius        float64
	damage        float64
	attackCadence float64
	cooldown      float64
}

func NewTower(pos pixel.Vec) *Tower {
	return &Tower{pos: pos, radius: 120, damage: 1.0, attackCadence: 0.16}
}

func (t *Tower) Update(dt float64, enemies []*Enemy) {
	t.cooldown -= dt
	if t.cooldown > 0 {
		return
	}

	var target *Enemy
	closest := 1e9
	for _, e := range enemies {
		if e.IsDead() {
			continue
		}
		d := e.pos.Sub(t.pos).Len()
		if d <= t.radius && d < closest {
			target = e
			closest = d
		}
	}

	if target != nil {
		target.TakeDamage(t.damage)
		t.cooldown = t.attackCadence
	}
}
