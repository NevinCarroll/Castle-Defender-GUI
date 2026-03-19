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

const (
	windowWidth  = 1024
	windowHeight = 768
	enemyReward  = 50
)

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

type gameState int

const (
	stateMenu gameState = iota
	stateTutorial
	statePlay
)

func drawCenteredText(target pixel.Target, txt *text.Text, lines []string, yStart float64) {
	txt.Clear()
	txt.Color = colornames.White
	for i, line := range lines {
		txt.Dot = pixel.V(windowWidth/2-float64(len(line))*3, yStart-float64(i)*30)
		txt.WriteString(line)
	}
	txt.Draw(target, pixel.IM.Moved(pixel.ZV))
}

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

	state := stateMenu

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		if win.JustPressed(pixel.KeyEscape) {
			break
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
			if gold >= cfg.Cost && !inPathArea(pos, path) {
				towers = append(towers, NewTower(pos, selectedTowerType))
				gold -= cfg.Cost
			}
			placePreview = false
		}

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
				enemies = append(enemies[:i], enemies[i+1:]...)
				i--
			}
		}

		for _, tower := range towers {
			tower.Update(dt, enemies)
		}

		if lives <= 0 {
			win.Clear(colornames.Black)
			txt.Clear()
			txt.Dot = pixel.V(320, 380)
			txt.Color = colornames.White
			txt.WriteString("GAME OVER")
			txt.Draw(win, pixel.IM.Moved(pixel.ZV))
			win.Update()
			time.Sleep(2 * time.Second)
			break
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
			previewDrawer.Color = colornames.Royalblue
			typeCfg := TowerConfigs[selectedTowerType]
			previewDrawer.Push(previewPos)
			previewDrawer.Circle(10, 0)
			previewDrawer.Color = colornames.Royalblue
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
			towerDrawer.Push(tower.pos.Add(pixel.Vec{X: tower.radius, Y: 0}), tower.pos.Add(pixel.Vec{X: -tower.radius, Y: 0}))
			towerDrawer.Line(1)
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

		txt.Clear()
		txt.Dot = pixel.V(10, windowHeight-20)
		txt.Color = colornames.White
		txt.WriteString(fmt.Sprintf("Lives: %d   Gold: %d   Wave: %d   Towers: %d   Enemies: %d   Selected: %s", lives, gold, wave, len(towers), len(enemies), towerTypeName(selectedTowerType)))
		txt.Draw(win, pixel.IM.Moved(pixel.ZV))

		win.Update()
	}
}

func inPathArea(pos pixel.Vec, path []pixel.Vec) bool {
	for i := 0; i < len(path)-1; i++ {
		if pointToSegmentDistance(pos, path[i], path[i+1]) < 30 {
			return true
		}
	}
	return false
}

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

func main() {
	opengl.Run(run)
}
