package main

import (
	"image/color"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 640
	screenHeight = 480

	rect_width  = 8
	rect_height = 16
	rect_axel   = float64(0.0000002)
	rect_res    = float64(0.9)
)

type Game struct {
	pressedKeys  []ebiten.Key
	pressedMouse []ebiten.MouseButton
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

var (
	rect_posx                                    = float64(screenWidth) / 2
	rect_posy                                    = float64(screenHeight) / 2
	rect_movx, rect_movy, rect_axelx, rect_axely = float64(0), float64(0), float64(0), float64(0)

	prevUpdTime = time.Now()
)

func (g *Game) Update() error {
	timeDelta := float64(time.Since(prevUpdTime))
	prevUpdTime = time.Now()

	g.pressedKeys = inpututil.AppendPressedKeys(g.pressedKeys[:0])

	rect_axelx = 0
	rect_axely = 0

	axel := rect_axel

	for _, key := range g.pressedKeys { //trajectory decision
		switch key.String() {
		case "ArrowDown":
			rect_axely = axel
		case "ArrowUp":
			rect_axely = -axel
		case "ArrowRight":
			rect_axelx = axel
		case "ArrowLeft":
			rect_axelx = -axel
		}
	}

	rect_movx += rect_axelx //moving
	rect_movy += rect_axely

	rect_movx *= rect_res //resistance & braking
	rect_movy *= rect_res

	rect_posx += rect_movx * timeDelta //time relation
	rect_posy += rect_movy * timeDelta

	const minX = rect_width
	const minY = rect_height
	const maxX = screenWidth - rect_width
	const maxY = screenHeight - rect_height

	if rect_posx >= maxX || rect_posx <= minX {
		if rect_posx > maxX {
			rect_posx = maxX
		} else if rect_posx < minX {
			rect_posx = minX
		}

		rect_movx *= -1
	}

	if rect_posy >= maxY || rect_posy <= minY {
		if rect_posy > maxY {
			rect_posy = maxY
		} else if rect_posy < minY {
			rect_posy = minY
		}

		rect_movy *= -1
	}

	return nil
}

//no shaders

func (g *Game) drawRectangle(screen *ebiten.Image, xcoord, ycoord, width, height float32, CLR color.Color) {
	var path vector.Path

	path.MoveTo(xcoord, ycoord)
	//path.Arc(xcoord,ycoord,) -- for circular objects

	x1coord := xcoord + width
	y1coord := ycoord + height
	for X := xcoord; X < x1coord; X++ {
		for Y := ycoord; Y < y1coord; Y++ {
			screen.Set(int(X), int(Y), CLR)
		}
	}
	ebitenutil.DebugPrint(screen, "!")
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "!")
	purpleClr := color.RGBA{255, 0, 255, 255}

	g.drawRectangle(screen, float32(rect_posx), float32(rect_posy), rect_width, rect_height, purpleClr)

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		ebiten.SetWindowTitle(strconv.FormatFloat(rect_posx, 'f', 5, 64))
		cursorX, cursorY := ebiten.CursorPosition()
		redClr := color.RGBA{255, 0, 0, 255}

		if cursorX < int(rect_posx) && cursorY < int(rect_posy) { //ЛВУ
			g.drawRectangle(screen, float32(cursorX), float32(cursorY), float32(rect_posx)-float32(rect_width), float32(rect_posy)-float32(rect_height), redClr)
		} else if cursorX < int(rect_posx) && cursorY > int(rect_posy) { //ЛНУ
			g.drawRectangle(screen, float32(rect_posx), float32(rect_posy), float32(cursorX)-float32(rect_posx), float32(cursorY)-float32(rect_posy), redClr)
		} else if cursorX > int(rect_posx) && cursorY < int(rect_posy) { //ПВУ
			g.drawRectangle(screen, float32(cursorX)-float32(rect_posx), float32(rect_posy), float32(rect_posx), float32(cursorY)-float32(rect_posy), redClr)
		} else { //ПНУ
			g.drawRectangle(screen, float32(rect_posx), float32(rect_posy), float32(cursorX)-float32(rect_posx), float32(cursorY)-float32(rect_posy), redClr)
		}

		//ebitenutil.DrawLine(g, rect_posx, rect_posy, float64(cursorX), float64(cursorY), redClr) -- doesn't work
	}
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Rectangle w/ physics")

	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}

/*
20.02
correct: now a red rectangle is drawn from the character to the cursor position, BUT:
the cursor position is counted from 0:0 and the rectangle (not straight, looool) is not built in the direction with negative coordinates relative to the character

*/
