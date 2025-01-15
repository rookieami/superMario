package main

import (
	"SuperMarioBros/common"
	"SuperMarioBros/entities"
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 160
	screenHeight = 160
	tileSize     = 16
)

var g *Game

func init() {
	floorImg, _, err := ebitenutil.NewImageFromFile("resource/tileset/SMB-Tiles.png")
	if err != nil {
		log.Fatal(err)
	}
	mapJson, err := common.LoadConfig("resource/map/word1.tmj")
	if err != nil {
		log.Fatal(err)
	}
	spriteImg, _, err := ebitenutil.NewImageFromFile("resource/img/ninja.png")
	if err != nil {
		log.Fatal(err)
	}
	g = NewGame()
	g.FloorImg = floorImg
	g.TiledMap = mapJson
	g.Sprite = &entities.Sprite{
		Img: spriteImg,
		X:   16.0,
		Y:   112.0,
	}
	g.Camera = &entities.Camera{
		X: 16.0,
		Y: 112.0,
	}

}

type Game struct {
	TiledMap *common.TileMap
	FloorImg *ebiten.Image
	Sprite   *entities.Sprite
	Camera   *entities.Camera
}

func NewGame() *Game {
	return &Game{}
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	// Player movement with arrow keys
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		if g.Sprite.X <= 0 {
			return nil
		}
		g.Sprite.X -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if g.Sprite.X >= screenWidth {
			fmt.Println("x", g.Sprite.X)
			return nil
		}
		g.Sprite.X += 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.Sprite.Y -= 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.Sprite.Y += 2
	}
	g.Camera.FollowTarget(g.Sprite.X, g.Sprite.Y, screenWidth, screenHeight)
	g.Camera.Constrain(20*16, 10*16, screenWidth, screenHeight)
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}
	tileXCount := g.FloorImg.Bounds().Dx() / tileSize
	xSpacing := 1
	ySpacing := 1
	for _, layer := range g.TiledMap.Layers {
		for index, id := range layer.Data {
			x := index % layer.Width
			y := index / layer.Width
			x *= tileSize
			y *= tileSize

			opts.GeoM.Translate(float64(x), float64(y))
			// 计算瓦片坐标
			srcX := ((id-1)%tileXCount)*tileSize + ((id-1)%tileXCount)*xSpacing
			srcY := ((id-1)/tileXCount)*tileSize + ((id-1)/tileXCount)*ySpacing

			screen.DrawImage(
				g.FloorImg.SubImage(
					image.Rect(srcX, srcY, srcX+tileSize, srcY+tileSize)).(*ebiten.Image),
				&opts,
			)
			opts.GeoM.Reset()
		}
	}
	opts.GeoM.Translate(g.Sprite.X, g.Sprite.Y)
	screen.DrawImage(
		g.Sprite.Img.SubImage(
			image.Rect(0, 0, 16, 16),
		).(*ebiten.Image),
		&opts,
	)
	opts.GeoM.Reset()
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("SuperMario")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
