package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
	"image/color"
	"math"
)

type Game struct {
	events        []Point
	snake         *Snake
	food          *Food
	score         int
	gameOver      bool
	updateCounter int
	delay         int
}

func (g *Game) Restart() {
	g.snake = NewSnake()
	g.score = 0
	g.events = []Point{}
	g.gameOver = false
	g.food = NewFood()
	g.delay = startDelay
}

func (g *Game) Update() error {
	if g.gameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			g.Restart()
		}
	}

	{
		if ebiten.IsKeyPressed(ebiten.KeyLeft) &&
			g.snake.Direction.X == 0 {
			g.events = append(g.events, Point{X: -1, Y: 0})
		} else if ebiten.IsKeyPressed(ebiten.KeyRight) &&
			g.snake.Direction.X == 0 {
			g.events = append(g.events, Point{X: 1, Y: 0})
		} else if ebiten.IsKeyPressed(ebiten.KeyUp) &&
			g.snake.Direction.Y == 0 {
			g.events = append(g.events, Point{X: 0, Y: -1})
		} else if ebiten.IsKeyPressed(ebiten.KeyDown) &&
			g.snake.Direction.Y == 0 {
			g.events = append(g.events, Point{X: 0, Y: 1})
		}
	}

	{
		g.updateCounter++
		if g.updateCounter < g.delay {
			return nil
		}
		g.updateCounter = 0

		if len(g.events) > 0 {
			x := Point{}
			x, g.events = g.events[0], g.events[1:]
			g.snake.Direction = x
		}

		g.snake.Move()

		g.events = []Point{}
	}

	{
		head := g.snake.Body[0]

		if head.X < 0 || head.Y < 0 || head.X >= screenWidth/tileSize || head.Y >= screenHeight/tileSize {
			g.gameOver = true
			g.delay = startDelay
		}

		for _, part := range g.snake.Body[1:] {
			if head.X == part.X && head.Y == part.Y {
				g.gameOver = true
				g.delay = startDelay
			}
		}

		if head.X == g.food.Position.X && head.Y == g.food.Position.Y {
			g.score++
			g.snake.GrowCounter += 1
			g.food = NewFood()

			if g.delay > minDelay {
				if math.Mod(float64(g.score), float64(5)) == 0 {
					g.delay--
				}
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Gray{})

	for _, part := range g.snake.Body {
		vector.DrawFilledRect(
			screen,
			float32(part.X*tileSize),
			float32(part.Y*tileSize),
			tileSize,
			tileSize,
			color.RGBA{R: 0, G: 255, B: 0, A: 255},
			false,
		)
	}

	vector.DrawFilledRect(
		screen,
		float32(g.food.Position.X*tileSize),
		float32(g.food.Position.Y*tileSize),
		tileSize,
		tileSize,
		color.RGBA{R: 255, G: 0, B: 0, A: 255},
		false,
	)

	face := basicfont.Face7x13

	if g.gameOver {
		text.Draw(
			screen,
			"Game over",
			face,
			screenWidth/2-40,
			screenHeight/2,
			color.White,
		)
		text.Draw(
			screen,
			"Press 'R' to restart",
			face,
			screenWidth/2-60,
			screenHeight/2+16,
			color.White,
		)
	}

	scoreText := fmt.Sprintf("Score: %d", g.score)
	text.Draw(
		screen,
		scoreText,
		face,
		5,
		screenHeight-5,
		color.White,
	)
}

func (g *Game) Layout(int, int) (int, int) {
	return screenWidth, screenHeight
}
