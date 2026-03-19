package main

import (
	"fmt"
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
	towerCost    = 100
	enemyReward  = 50
)

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
				"Left click to place towers (cost 100)",
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

		if win.JustPressed(pixel.MouseButtonLeft) {
			pos := win.MousePosition()
			if gold >= towerCost && !inPathArea(pos, path) {
				towers = append(towers, NewTower(pos))
				gold -= towerCost
			}
		}

		spawnTimer += dt
		if spawnTimer >= spawnInterval {
			wave++
			spawnEnemyWave(&enemies, path)
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

		towerDrawer := imdraw.New(nil)
		for _, tower := range towers {
			towerDrawer.Color = colornames.Gold
			towerDrawer.Push(tower.pos)
			towerDrawer.Circle(10, 0)
			towerDrawer.Color = colornames.Royalblue
			towerDrawer.Push(tower.pos)
			towerDrawer.Circle(2, 0)
			towerDrawer.Push(tower.pos.Add(pixel.Vec{X: tower.radius, Y: 0}), tower.pos.Add(pixel.Vec{X: -tower.radius, Y: 0}))
			towerDrawer.Line(1)
		}
		towerDrawer.Draw(win)

		enemyDrawer := imdraw.New(nil)
		for _, enemy := range enemies {
			enemyDrawer.Color = colornames.Crimson
			enemyDrawer.Push(enemy.pos)
			enemyDrawer.Circle(7, 0)
			enemyDrawer.Line(3)
		}
		enemyDrawer.Draw(win)

		txt.Clear()
		txt.Dot = pixel.V(10, windowHeight-20)
		txt.Color = colornames.White
		txt.WriteString(fmt.Sprintf("Lives: %d   Gold: %d   Wave: %d   Towers: %d   Enemies: %d", lives, gold, wave, len(towers), len(enemies)))
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

func spawnEnemyWave(enemies *[]*Enemy, path []pixel.Vec) {
	for i := 0; i < 2; i++ {
		*enemies = append(*enemies, NewEnemy(path[0], 90+float64(i*15), 7))
	}
}

func main() {
	opengl.Run(run)
}
