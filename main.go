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

// States that the game is in, used to determine what should be displayed
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
	// Take window screen, half it the x and y value, then draw the text there
	for i, line := range lines {
		txt.Dot = pixel.V(windowWidth/2-float64(len(line))*3, yStart-float64(i)*30)
		txt.WriteString(line)
	}
	txt.Draw(target, pixel.IM.Moved(pixel.ZV))
}

// run initializes the window and enters the game loop, handling state transitions, input, updates, and rendering.
func run() {
	// Set up configs of game
	cfg := opengl.WindowConfig{
		Title:  "Castle Defender",
		Bounds: pixel.R(0, 0, windowWidth, windowHeight),
		VSync:  true,
	}
	win, err := opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII) // Set font of text
	txt := text.New(pixel.ZV, atlas) 

	// Set up path for enemies to follow
	path := []pixel.Vec{
		{80, 160}, {240, 160}, {240, 420}, {560, 420}, {560, 220}, {840, 220}, {840, 640}, {940, 640},
	}

	// Varaiables that manage gameplay
	enemies := []*Enemy{}
	towers := []*Tower{}
	lives := 5
	gold := 300
	wave := 0
	spawnTimer := 0.0
	spawnInterval := 2.2
	last := time.Now()

	// For tower placement
	selectedTowerType := TowerTypeStandard
	placePreview := false
	previewPos := pixel.ZV
	lasers := []*Laser{}

	// Player stats to display when game is over
	goldEarned := 0
	enemiesKilled := 0
	towersPlaced := 0
	wavesSurvived := 0
	gameOverReason := ""

	state := stateMenu // Set game state to main menu

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

	// While game is running
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

		// Display correct menu depending on game state
		switch state {
		case stateMenu:
			// Display main menu
			win.Clear(colornames.Darkslategray)
			lines := []string{"Castle Defender", "", "Press ENTER to continue", "", "(Esc to quit)"}
			drawCenteredText(win, txt, lines, 460)
			win.Update()
			if win.JustPressed(pixel.KeyEnter) {
				state = stateTutorial
			}
			continue
		case stateTutorial:
			// Display tutorial
			win.Clear(colornames.Darkslategray)
			lines := []string{
				"TUTORIAL",
				"",
				"Press 1-3 to choose tower type",
				"Then left click to place selected tower",
				"Or right click to cancel tower placement",
				"A tower preview will show you tower range and if you are allowed to place it",
				"Survive as many waves of enemies as you can",
				"or press quit to end game early",
				"",
				"Press ENTER to start playing",
			}
			drawCenteredText(win, txt, lines, 500)
			win.Update()
			if win.JustPressed(pixel.KeyEnter) { // Start game
				state = statePlay
			}
			continue
		case stateGameOver:
			// Display game over screen
			win.Clear(colornames.Black)
			// Display player stats
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
			if win.JustPressed(pixel.KeyEnter) { // Start new game
				resetGame()
				state = statePlay
				continue
			}
			if win.JustPressed(pixel.KeyQ) || win.JustPressed(pixel.KeyEscape) { // Quit game
				return
			}
			continue
		}

		// Select tower type depending on which button the player pressed
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

		// Show tower preview at cursor
		if placePreview {
			previewPos = win.MousePosition()
		}

		// Cancel tower placement
		if win.JustPressed(pixel.MouseButtonRight) {
			placePreview = false
		}

		// Place down tower after player pressed left click while previewing tower placement
		if win.JustPressed(pixel.MouseButtonLeft) && placePreview {
			pos := win.MousePosition() // Place tower at mouse position
			cfg := TowerConfigs[selectedTowerType]
			// Ensure player has enough gold and position is valid
			if gold >= cfg.Cost && isValidPlacement(pos, path, towers) {
				towers = append(towers, NewTower(pos, selectedTowerType))
				gold -= cfg.Cost
				towersPlaced++
			}
			placePreview = false // Stop previewing tower
		}

		// Periodically spawn waves; spawn interval gradually decreases to escalate difficulty.
		spawnTimer += dt
		if spawnTimer >= spawnInterval {
			wave++
			spawnEnemyWave(&enemies, path, wave) // Spawn enemies
			spawnTimer -= spawnInterval 
			spawnInterval *= 0.98 // Decrease spawn interval
			if spawnInterval < 0.6 {
				spawnInterval = 0.6
			}
		}

		// Update all enemies, remove those that reached the end or died.
		for i := 0; i < len(enemies); i++ {
			en := enemies[i]
			en.Update(dt, path) // Move enemies
			// Make player lose life if enemy reaches the end, and despawn them
			if en.ReachedEnd(path) {
				lives--
				enemies = append(enemies[:i], enemies[i+1:]...)
				i--
				continue
			}
			if en.IsDead() { // Give player gold for killing enemy
				gold += enemyReward
				goldEarned += enemyReward
				enemiesKilled++
				enemies = append(enemies[:i], enemies[i+1:]...)
				i--
			}
		}

		// Make tower shot closest enemy if it can
		for _, tower := range towers {
			shot := tower.Update(dt, enemies)
			// Add laser to list so that it will be rendered
			if shot != nil {
				lasers = append(lasers, &Laser{start: tower.pos, end: shot.pos, time: 0.15})
			}
		}

		// Make player lose if they lose all lives
		if lives <= 0 {
			state = stateGameOver
			wavesSurvived = wave
			gameOverReason = "All lives lost"
			continue
		}

		// Clear screen with background
		win.Clear(colornames.Darkslategray)

		// Draw path enemies will follow
		lineDrawer := imdraw.New(nil)
		lineDrawer.Color = colornames.Yellow
		for i := 0; i < len(path)-1; i++ {
			lineDrawer.Push(path[i], path[i+1])
			lineDrawer.Line(4)
		}
		lineDrawer.Draw(win)

		// Draw path nodes
		pathPoints := imdraw.New(nil)
		for _, p := range path {
			pathPoints.Color = colornames.White
			pathPoints.Push(p)
			pathPoints.Circle(4, 0)
		}
		pathPoints.Draw(win)

		// Draw tower preview when player wants to place a tower
		if placePreview {
			previewDrawer := imdraw.New(nil)
			typeCfg := TowerConfigs[selectedTowerType]
			valid := isValidPlacement(previewPos, path, towers)
			// If tower placement is invalid, then display preview as red rather than blue
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

		// Draw towers on screen
		towerDrawer := imdraw.New(nil)
		for _, tower := range towers {
			// Determine tower type and draw corresponding color
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

		// Draw enemies on screen
		enemyDrawer := imdraw.New(nil)
		for _, enemy := range enemies {
			// Determine enemy type and draw corresponding color
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
			alpha := laser.time / 0.15 // Slowly decrease laser alpha as time goes on
			laserDrawer := imdraw.New(nil)
			laserDrawer.Color = colornames.Orange
			laserDrawer.Push(laser.start, laser.end)
			laserDrawer.Line(2)
			laserDrawer.Draw(win)
			_ = alpha
		}

		// Draw text and player stats
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
		if distancePointToSegment(pos, path[i], path[i+1]) < 30 {
			return true
		}
	}
	return false
}

// distancePointToSegment returns the shortest Euclidean distance from point p to the segment [segA, segB].
// It projects p onto the infinite line through the segment, clamps to the endpoints, and computes distance to the closest point.
func distancePointToSegment(p, segA, segB pixel.Vec) float64 {
	segmentVector := segB.Sub(segA)
	pointVector := p.Sub(segA)
	squaredLen := segmentVector.Dot(segmentVector)
	if squaredLen == 0 {
		return p.Sub(segA).Len()
	}
	
	t := pointVector.Dot(segmentVector) / squaredLen
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	closest := segA.Add(segmentVector.Scaled(t))
	return p.Sub(closest).Len()
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
		// Chooses a random enemy type based off random chance
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
