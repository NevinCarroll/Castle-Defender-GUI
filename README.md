# Castle Defender (Go + Pixel)

A simple tower-defense game implemented in Go using the Pixel game library.

## Overview

In this game, you defend a castle path by placing towers that automatically shoot enemies. Enemies travel a fixed path, and each wave gets harder over time. Survive as many waves as possible.

### Features

- Menu, tutorial, gameplay, and game-over states
- Three tower types:
  - Standard: medium range/damage cadence
  - Rapid: fast fire rate, lower damage
  - Sniper: long range, high damage, slower rate
- Randomized waves with enemy types:
  - Default
  - Fast
  - Tank
- Enemy path following, health/damage logic, and end-of-path life loss
- Real-time Tower placement preview with valid/invalid placement feedback
- Player stats and game-over summary

## Controls

- `ENTER` to advance from menu/tutorial and start a new game
- `1`, `2`, `3` to select tower type
- Left-click to place selected tower (if valid)
- Right-click to cancel placement
- `ESC` to quit (or end the game during play)
- `Q` from game-over to quit

## Gameplay

- Start with 5 lives and 300 gold.
- Each enemy kill gives +25 gold.
- Each wave spawns more enemies and the wave spawn interval slightly decreases over time.
- You lose when lives reach 0.

## Project Structure

- `main.go` - game loop, window rendering, state management, enemy wave logic, UI
- `tower.go` - tower types, configs, targeting/attack logic
- `enemy.go` - enemy types, movement/path following, health logic
- `go.mod` - Go module and dependencies

## Requirements

- Go 1.25+ (or compatible Go version)
- GPU/OpenGL support (Pixel requires OpenGL context)
- Dependencies are managed by Go modules (`go mod tidy` will fetch them)

## Run

From project root:

```bash
go run .
```

If dependencies are missing, run:

```bash
go mod tidy
```

Then run again.

## Build

```bash
go build -o castle-defender .
```

