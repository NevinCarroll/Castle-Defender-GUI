// Package main implements a simple tower defense game built with Pixel.
// It supports menu/tutorial/gameplay/game-over states, tower placement, enemy waves,
// and attack logic for standard, rapid, and sniper towers.
package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"github.com/gopxl/pixel/v2/ext/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

// Game window constants.
const (
	windowWidth  = 1024
	windowHeight = 768
	enemyReward  = 25
)

// towerTypeName returns the display label for the selected TowerType.
func towerTypeName(t TowerType) string {
	switch t {
	case TowerTypeStandard:
		return "1: Standard"
	case TowerTypeRapid:
		return "2: Rapid"
	case TowerTypeSniper:
		return "3: Sniper"
	default:
		return "None"
	}
}

// Laser represents a temporary visual beam fired by a tower.
type Laser struct {
	start pixel.Vec
	end   pixel.Vec
	time  float64
}

// gameState defines the finite states used by the main game loop.
type gameState int

const (
	stateMenu gameState = iota
	stateTutorial
	statePlay
	stateGameOver
)

// drawCenteredText writes the provided lines centered horizontally on the screen at yStart.
func drawCenteredText(target pixel.Target, txt *text.Text, lines []string, yStart float64) {
	txt.Clear()
	txt.Color = colornames.White
	for i, line := range lines {
		txt.Dot = pixel.V(windowWidth/2-float64(len(line))*3, yStart-float64(i)*30)
		txt.WriteString(line)
	}
	txt.Draw(target, pixel.IM.Moved(pixel.ZV))
}

// run initializes the window and enters the game loop, handling state transitions, input, updates, and rendering.
func run() {
	cfg := opengl.WindowConfig{
		Title:  "Pixel Tower Defense",
		Bounds: pixel.R(0, 0, windowWidth, windowHeight),
		VSync:  true,
	}
	win, err := opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	txt := text.New(pixel.ZV, atlas)

	path := []pixel.Vec{
		{80, 160}, {240, 160}, {240, 420}, {560, 420}, {560, 220}, {840, 220}, {840, 640}, {940, 640},
	}

	enemies := []*Enemy{}
	towers := []*Tower{}
	lives := 5
	gold := 300
	wave := 0
	spawnTimer := 0.0
	spawnInterval := 2.2
	last := time.Now()

	selectedTowerType := TowerTypeStandard
	placePreview := false
	previewPos := pixel.ZV
	lasers := []*Laser{}

	goldEarned := 0
	enemiesKilled := 0
	towersPlaced := 0
	wavesSurvived := 0
	gameOverReason := ""

	state := stateMenu

	// resetGame resets all gameplay variables for a new session.
	resetGame := func() {
		enemies = []*Enemy{}
		towers = []*Tower{}
		lives = 5
		gold = 300
		wave = 0
		spawnTimer = 0.0
		spawnInterval = 2.2
		selectedTowerType = TowerTypeStandard
		placePreview = false
		previewPos = pixel.ZV
		lasers = []*Laser{}
		goldEarned = 0
		enemiesKilled = 0
		towersPlaced = 0
		wavesSurvived = 0
		gameOverReason = ""
	}

	resetGame()

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		// Global shortcut: ESC closes the game unless playing, in which case it ends the run.
		if win.JustPressed(pixel.KeyEscape) {
			if state == statePlay {
				state = stateGameOver
				wavesSurvived = wave
				gameOverReason = "Player quit"
			} else {
				break
			}
		}

		switch state {
		case stateMenu:
			win.Clear(colornames.Darkslategray)
			lines := []string{"PIXEL TOWER DEFENSE", "", "Press ENTER to continue", "", "(Esc to quit)"}
			drawCenteredText(win, txt, lines, 460)
			win.Update()
			if win.JustPressed(pixel.KeyEnter) {
				state = stateTutorial
			}
			continue
		case stateTutorial:
			win.Clear(colornames.Darkslategray)
			lines := []string{
				"TUTORIAL",
				"",
				"Press 1-3 to choose tower type",
				"Then left click to place selected tower",
				"Don't place towers on the path",
				"Survive waves of enemies",
				"",
				"Press ENTER to start playing",
			}
			drawCenteredText(win, txt, lines, 500)
			win.Update()
			if win.JustPressed(pixel.KeyEnter) {
				state = statePlay
			}
			continue
		case stateGameOver:
			win.Clear(colornames.Black)
			lines := []string{
				"GAME OVER",
				"",
				fmt.Sprintf("Reason: %s", gameOverReason),
				"",
				fmt.Sprintf("Gold earned: %d", goldEarned),
				fmt.Sprintf("Enemies killed: %d", enemiesKilled),
				fmt.Sprintf("Towers placed: %d", towersPlaced),
				fmt.Sprintf("Waves survived: %d", wavesSurvived),
				"",
				"Press ENTER to start a new game",
				"Press Q or ESC to quit",
			}
			drawCenteredText(win, txt, lines, 500)
			win.Update()
			if win.JustPressed(pixel.KeyEnter) {
				resetGame()
				state = statePlay
				continue
			}
			if win.JustPressed(pixel.KeyQ) || win.JustPressed(pixel.KeyEscape) {
				return
			}
			continue
		}

		if win.JustPressed(pixel.Key1) {
			selectedTowerType = TowerTypeStandard
			placePreview = true
		}
		if win.JustPressed(pixel.Key2) {
			selectedTowerType = TowerTypeRapid
			placePreview = true
		}
		if win.JustPressed(pixel.Key3) {
			selectedTowerType = TowerTypeSniper
			placePreview = true
		}

		if placePreview {
			previewPos = win.MousePosition()
		}

		if win.JustPressed(pixel.MouseButtonRight) {
			placePreview = false
		}

		if win.JustPressed(pixel.MouseButtonLeft) && placePreview {
			pos := win.MousePosition()
			cfg := TowerConfigs[selectedTowerType]
			if gold >= cfg.Cost && isValidPlacement(pos, path, towers) {
				towers = append(towers, NewTower(pos, selectedTowerType))
				gold -= cfg.Cost
				towersPlaced++
			}
			placePreview = false
		}

		// Periodically spawn waves; spawn interval gradually decreases to escalate difficulty.
		spawnTimer += dt
		if spawnTimer >= spawnInterval {
			wave++
			spawnEnemyWave(&enemies, path, wave)
			spawnTimer -= spawnInterval
			spawnInterval *= 0.98
			if spawnInterval < 0.6 {
				spawnInterval = 0.6
			}
		}

		// Update all enemies, remove those that reached the end or died.
		for i := 0; i < len(enemies); i++ {
			en := enemies[i]
			en.Update(dt, path)
			if en.ReachedEnd(path) {
				lives--
				enemies = append(enemies[:i], enemies[i+1:]...)
				i--
				continue
			}
			if en.IsDead() {
				gold += enemyReward
				goldEarned += enemyReward
				enemiesKilled++
				enemies = append(enemies[:i], enemies[i+1:]...)
				i--
			}
		}

		for _, tower := range towers {
			shot := tower.Update(dt, enemies)
			if shot != nil {
				lasers = append(lasers, &Laser{start: tower.pos, end: shot.pos, time: 0.15})
			}
		}

		if lives <= 0 {
			state = stateGameOver
			wavesSurvived = wave
			gameOverReason = "All lives lost"
			continue
		}

		win.Clear(colornames.Darkslategray)

		lineDrawer := imdraw.New(nil)
		lineDrawer.Color = colornames.Yellow
		for i := 0; i < len(path)-1; i++ {
			lineDrawer.Push(path[i], path[i+1])
			lineDrawer.Line(4)
		}
		lineDrawer.Draw(win)

		pathPoints := imdraw.New(nil)
		for _, p := range path {
			pathPoints.Color = colornames.White
			pathPoints.Push(p)
			pathPoints.Circle(4, 0)
		}
		pathPoints.Draw(win)

		if placePreview {
			previewDrawer := imdraw.New(nil)
			typeCfg := TowerConfigs[selectedTowerType]
			valid := isValidPlacement(previewPos, path, towers)
			if valid {
				previewDrawer.Color = colornames.Royalblue
			} else {
				previewDrawer.Color = colornames.Red
			}
			previewDrawer.Push(previewPos)
			previewDrawer.Circle(10, 0)
			previewDrawer.Push(previewPos)
			previewDrawer.Circle(typeCfg.Radius, 1)
			previewDrawer.Draw(win)
		}

		towerDrawer := imdraw.New(nil)
		for _, tower := range towers {
			switch tower.typeID {
			case TowerTypeStandard:
				towerDrawer.Color = colornames.Gold
			case TowerTypeRapid:
				towerDrawer.Color = colornames.Limegreen
			case TowerTypeSniper:
				towerDrawer.Color = colornames.Mediumblue
			default:
				towerDrawer.Color = colornames.Gold
			}
			towerDrawer.Push(tower.pos)
			towerDrawer.Circle(10, 0)
			towerDrawer.Color = colornames.White
			towerDrawer.Push(tower.pos)
			towerDrawer.Circle(2, 0)
		}
		towerDrawer.Draw(win)

		enemyDrawer := imdraw.New(nil)
		for _, enemy := range enemies {
			switch enemy.typeID {
			case EnemyTypeFast:
				enemyDrawer.Color = colornames.Plum
			case EnemyTypeTank:
				enemyDrawer.Color = colornames.Darkred
			default:
				enemyDrawer.Color = colornames.Crimson
			}
			enemyDrawer.Push(enemy.pos)
			enemyDrawer.Circle(7, 0)
			enemyDrawer.Line(3)
		}
		enemyDrawer.Draw(win)

		// Draw lasers as short traces and fade them out
		for i := len(lasers) - 1; i >= 0; i-- {
			laser := lasers[i]
			laser.time -= dt
			if laser.time <= 0 {
				lasers = append(lasers[:i], lasers[i+1:]...)
				continue
			}
			alpha := laser.time / 0.15
			laserDrawer := imdraw.New(nil)
			laserDrawer.Color = colornames.Orange
			laserDrawer.Push(laser.start, laser.end)
			laserDrawer.Line(2)
			laserDrawer.Draw(win)
			_ = alpha
		}

		txt.Clear()
		txt.Dot = pixel.V(10, windowHeight-20)
		txt.Color = colornames.White
		txt.WriteString(fmt.Sprintf("Lives: %d   Gold: %d   Wave: %d   Towers: %d   Enemies: %d   Selected: %s", lives, gold, wave, len(towers), len(enemies), towerTypeName(selectedTowerType)))
		txt.Draw(win, pixel.IM.Moved(pixel.ZV))

		win.Update()
	}
}
// inPathArea returns true when the position is within 30 pixels of any path segment.
func inPathArea(pos pixel.Vec, path []pixel.Vec) bool {
	for i := 0; i < len(path)-1; i++ {
		if pointToSegmentDistance(pos, path[i], path[i+1]) < 30 {
			return true
		}
	}
	return false
}

// pointToSegmentDistance returns the shortest Euclidean distance from point p to segment [a,b].
func pointToSegmentDistance(p, a, b pixel.Vec) float64 {
	ab := b.Sub(a)
	ap := p.Sub(a)
	denom := ab.Dot(ab)
	if denom == 0 {
		return p.Sub(a).Len()
	}
	t := ap.Dot(ab) / denom
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	proj := a.Add(ab.Scaled(t))
	return p.Sub(proj).Len()
}

// isTowerOverlap returns true if a tower is too close to an existing tower.
func isTowerOverlap(pos pixel.Vec, towers []*Tower) bool {
	for _, t := range towers {
		if pos.Sub(t.pos).Len() < 24 {
			return true
		}
	}
	return false
}

// isValidPlacement returns true if pos is not on the path and not overlapping another tower.
func isValidPlacement(pos pixel.Vec, path []pixel.Vec, towers []*Tower) bool {
	if inPathArea(pos, path) {
		return false
	}
	if isTowerOverlap(pos, towers) {
		return false
	}
	return true
}

// spawnEnemyWave queues a new wave of enemies at the path start with randomized types.
func spawnEnemyWave(enemies *[]*Enemy, path []pixel.Vec, wave int) {
	enemiesThisWave := 2 + (wave / 5)
	for i := 0; i < enemiesThisWave; i++ {
		r := rand.Float64()
		typeID := EnemyTypeDefault
		if r < 0.25 {
			typeID = EnemyTypeTank
		} else if r < 0.6 {
			typeID = EnemyTypeFast
		}
		*enemies = append(*enemies, NewEnemy(path[0], typeID))
	}
}

// main starts the opengl Pixel runtime and launches run().
func main() {
	opengl.Run(run)
}
