package main

import (
    "fmt"
    "time"

    "github.com/gopxl/pixel/v2"
    "github.com/gopxl/pixel/v2/backends/opengl"
    "golang.org/x/image/colornames"
    "golang.org/x/image/font/basicfont"
)

const (
    windowWidth  = 1024
    windowHeight = 768
    towerCost    = 100
    enemyReward  = 50
)

type Enemy struct {
    pos       pixel.Vec
    pathPoint int
    speed     float64
    health    float64
    maxHealth float64
}

func (e *Enemy) update(dt float64, path []pixel.Vec) {
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

func (e *Enemy) reachedEnd(path []pixel.Vec) bool {
    return e.pathPoint >= len(path)-1
}

func (e *Enemy) inRange(point pixel.Vec, radius float64) bool {
    return e.pos.Sub(point).Len() <= radius
}

type Tower struct {
    pos           pixel.Vec
    radius        float64
    damage        float64
    attackCadence float64
    cooldown      float64
}

func (t *Tower) update(dt float64, enemies []*Enemy) {
    t.cooldown -= dt
    if t.cooldown > 0 {
        return
    }

    var target *Enemy
    closest := 1e9
    for _, e := range enemies {
        if e.health <= 0 {
            continue
        }
        d := e.pos.Sub(t.pos).Len()
        if d <= t.radius && d < closest {
            target = e
            closest = d
        }
    }
    if target != nil {
        target.health -= t.damage
        t.cooldown = t.attackCadence
    }
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

    for !win.Closed() {
        dt := time.Since(last).Seconds()
        last = time.Now()

        if win.JustPressed(pixel.KeyEscape) {
            break
        }

        if win.JustPressed(pixel.MouseButtonLeft) {
            pos := win.MousePosition()
            if gold >= towerCost && !inPathArea(pos, path) {
                towers = append(towers, &Tower{pos: pos, radius: 120, damage: 1.0, attackCadence: 0.16})
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
            en.update(dt, path)
            if en.reachedEnd(path) {
                lives--
                enemies = append(enemies[:i], enemies[i+1:]...)
                i--
                continue
            }
            if en.health <= 0 {
                gold += enemyReward
                enemies = append(enemies[:i], enemies[i+1:]...)
                i--
            }
        }

        for _, tower := range towers {
            tower.update(dt, enemies)
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

        for i := 0; i < len(path)-1; i++ {
            target := pixel.NewSprite(pixel.PictureDataFromImage(pixel.NewImage(1, 1, colornames.White)), pixel.R(0, 0, 1, 1))
            target.Draw(win, pixel.IM.Scaled(pixel.ZV, 1).Moved(path[i]))
            line := pixel.NewLine(path[i], path[i+1])
            line.Draw(win, pixel.IM.Scaled(pixel.ZV, 4), colornames.Yellow)
        }

        for _, p := range path {
            handle := pixel.NewSprite(pixel.PictureDataFromImage(pixel.NewImage(1, 1, colornames.White)), pixel.R(0, 0, 1, 1))
            handle.Draw(win, pixel.IM.Scaled(pixel.ZV, 8).Moved(p), colornames.White)
        }

        for _, tower := range towers {
            towerCircle := pixel.NewSprite(pixel.PictureDataFromImage(pixel.NewImage(1, 1, colornames.Gold)), pixel.R(0, 0, 1, 1))
            towerCircle.Draw(win, pixel.IM.Scaled(pixel.ZV, 20).Moved(tower.pos))
            rangeCircle := pixel.NewCircle(tower.pos, tower.radius)
            rangeCircle.Draw(win, pixel.IM.Scaled(pixel.ZV, 1), colornames.Royalblue)
        }

        for _, enemy := range enemies {
            enemyCircle := pixel.NewSprite(pixel.PictureDataFromImage(pixel.NewImage(1, 1, colornames.Crimson)), pixel.R(0, 0, 1, 1))
            enemyCircle.Draw(win, pixel.IM.Scaled(pixel.ZV, 12).Moved(enemy.pos))
            health := enemy.health / enemy.maxHealth
            healthRect := pixel.R(enemy.pos.X-10, enemy.pos.Y+12, enemy.pos.X-10+20*health, enemy.pos.Y+16)
            bar := pixel.NewSprite(pixel.PictureDataFromImage(pixel.NewImage(int(20*health), 4, colornames.Limegreen)), pixel.R(0,0, float64(20*health),4))
            bar.Draw(win, pixel.IM.Moved(pixel.V{X: enemy.pos.X - 10, Y: enemy.pos.Y + 12}))
            _ = healthRect
        }

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
    t := ap.Dot(ab) / ab.Dot(ab)
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
        *enemies = append(*enemies, &Enemy{pos: path[0], pathPoint: 0, speed: 90 + float64(i*15), health: 7, maxHealth: 7})
    }
}

func main() {
    opengl.Run(run)
}