package main

import (
	"SuperMarioBros/entities"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/lafriks/go-tiled"
)

const (
	tileSize  = 16
	gravity   = 0.5
	jumpSpeed = 8
)

var g *Game

func init() {
	gameMap, err := tiled.LoadFile("resource/map/word1.tmx")
	if err != nil {
		log.Fatal(err)
	}
	floorImg, _, err := ebitenutil.NewImageFromFile("resource/tileset/SMB-Tiles.png")
	if err != nil {
		log.Fatal(err)
	}

	spriteImg, _, err := ebitenutil.NewImageFromFile("resource/img/ninja.png")
	if err != nil {
		log.Fatal(err)
	}
	g = NewGame()
	g.screenWidth = 320
	g.screenHeight = 160
	g.FloorImg = floorImg
	g.TiledMap = gameMap
	g.Sprite = &entities.Sprite{
		Img:       spriteImg,
		X:         16,
		Y:         112,
		VX:        16,
		JumpState: 0,
	}
	g.Camera = &entities.Camera{
		X: g.Sprite.X - float64(g.screenWidth/3),
		Y: g.Sprite.Y - float64(g.screenHeight/2),
	}
	// 约束相机位置
	mapWidth := float64(g.TiledMap.Width * tileSize)
	mapHeight := float64(g.TiledMap.Height * tileSize)
	g.Camera.Constrain(mapWidth, mapHeight, float64(g.screenWidth), float64(g.screenHeight))
}

type Game struct {
	TiledMap     *tiled.Map
	FloorImg     *ebiten.Image
	Sprite       *entities.Sprite
	Camera       *entities.Camera
	screenWidth  int
	screenHeight int
}

func NewGame() *Game {
	return &Game{}
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
func (g *Game) Update() error {
	spriteSpeed := 1
	if g.Sprite.JumpState > 1 {
		spriteSpeed = 2
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		if g.Sprite.VX <= 0 {
			return nil
		}
		g.Sprite.VX -= float64(spriteSpeed * jumpSpeed)
		g.Sprite.IsLeft = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		if g.Sprite.VX >= float64(g.TiledMap.Width-1)*tileSize {
			return nil
		}
		g.Sprite.VX += float64(spriteSpeed * jumpSpeed)
		if g.Sprite.IsLeft {
			g.Sprite.IsLeft = false
		}
	}
	if (g.Sprite.JumpState < 2) && inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if g.Sprite.Y < 2 {
			return nil
		}
		fmt.Println("开始跳跃")
		fmt.Println("--------")
		g.Sprite.Y -= jumpSpeed * 2
		g.Sprite.JumpState += 1
	}
	//计算当前帧位置
	if g.Sprite.X != g.Sprite.VX {
		if g.Sprite.IsLeft {
			//  <-
			g.Sprite.X -= jumpSpeed
		} else {
			//  ->
			g.Sprite.X += jumpSpeed
		}
		// 碰撞检测
		g.checkCollision()
	}
	if g.Sprite.JumpState != 0 {
		g.Sprite.Y += gravity //模拟重力
		// 碰撞检测
		g.checkCollision()
	}
	g.Camera.FollowTarget(g.Sprite.X, g.Sprite.Y, float64(g.screenWidth), float64(g.screenHeight))
	// 约束相机位置
	// 跟踪角色
	mapWidth := float64(g.TiledMap.Width * tileSize)
	mapHeight := float64(g.TiledMap.Height * tileSize)
	g.Camera.Constrain(mapWidth, mapHeight, float64(g.screenWidth), float64(g.screenHeight))

	return nil
}

// checkCollision 检测角色与地图之间的碰撞
func (g *Game) checkCollision() {
	// 获取角色的边界框
	spriteRect := image.Rect(int(g.Sprite.X), int(g.Sprite.Y), int(g.Sprite.X)+16, int(g.Sprite.Y)+16)
	fmt.Println("角色位置:", spriteRect)
	// 计算玩家周围需要检测的瓦片范围
	// startX := int(g.Sprite.X)/tileSize - 1
	// startY := int(g.Sprite.Y)/tileSize - 1
	// endX := startX + 3 // 检测3个瓦片宽度
	// endY := startY + 3 // 检测3个瓦片高度

	// for _, layer := range g.TiledMap.Layers {
	// 	for y := startY; y < endY; y++ {
	// 		for x := startX; x < endX; x++ {
	// 			if x <= 0 || y <= 0 || x >= g.TiledMap.Width || y >= g.TiledMap.Height {
	// 				continue
	// 			}

	// 			index := y*layer.Width + x
	// 			id := layer.Data[index]
	// 			if id == 0 {
	// 				continue
	// 			}

	// 			// 获取瓦片的边界框
	// 			tileRect := image.Rect(x*tileSize, y*tileSize, (x+1)*tileSize, (y+1)*tileSize)

	// 			// 检测角色与瓦片是否相交
	// 			if spriteRect.Overlaps(tileRect) {
	// 				// 处理碰撞
	// 				if spriteRect.Min.Y < tileRect.Max.Y && spriteRect.Max.Y > tileRect.Min.Y {
	// 					if g.Sprite.VX < 0 {
	// 						// 左侧碰撞
	// 						g.Sprite.X = float64(tileRect.Max.X)
	// 					} else if g.Sprite.VX > 0 {
	// 						// 右侧碰撞
	// 						g.Sprite.X = float64(tileRect.Min.X) - 16
	// 					}
	// 				}
	// 				if spriteRect.Min.X < tileRect.Max.X && spriteRect.Max.X > tileRect.Min.X {
	// 					if g.Sprite.Y < float64(tileRect.Max.Y) && g.Sprite.Y > float64(tileRect.Min.Y) {
	// 						// 下方碰撞（地面）
	// 						g.Sprite.Y = float64(tileRect.Max.Y)
	// 						g.Sprite.JumpState = 0

	// 					} else if g.Sprite.Y+16 > float64(tileRect.Min.Y) && g.Sprite.Y+16 < float64(tileRect.Max.Y) {
	// 						// 上方碰撞（头部）
	// 						g.Sprite.Y = float64(tileRect.Min.Y) - 16
	// 						g.Sprite.JumpState = 0
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {
	opts := ebiten.DrawImageOptions{}
	tileXCount := g.FloorImg.Bounds().Dx() / tileSize
	xSpacing := 1
	ySpacing := 1

	// 计算屏幕可见区域的瓦片范围
	startX := int(g.Camera.X) / tileSize
	startY := int(g.Camera.Y) / tileSize
	endX := startX + g.screenWidth/tileSize
	endY := startY + g.screenHeight/tileSize

	// 计算缩放比例
	scaleX := float64(g.screenWidth) / float64(g.TiledMap.Width*tileSize)
	scaleY := float64(g.screenHeight) / float64(g.TiledMap.Height*tileSize)
	scale := math.Min(scaleX, scaleY)
	opts.GeoM.Scale(scale, scale)
	fmt.Println("endx endy", endX, endY)
	// 渲染地图
	for _, layer := range g.TiledMap.Layers {
		for y := startY; y < endY; y++ {
			for x := startX; x < endX; x++ {
				if x < 0 || y < 0 || x > g.TiledMap.Width || y > g.TiledMap.Height {
					fmt.Println("x y w h", x, y, g.TiledMap.Width, g.TiledMap.Height)
					continue
				}

				index := y*g.TiledMap.Width + x
				if index >= len(layer.Tiles) {
					fmt.Println("index w x limit", index, g.TiledMap.Width, x, len(layer.Tiles))
					continue
				}
				tile := layer.Tiles[index]
				if tile == nil {
					fmt.Println("tile id ", tile.ID, index)
					continue
				}
				// 计算瓦片在图集中的位置
				srcX := (int(tile.ID) % tileXCount) * (tileSize + xSpacing)
				srcY := (int(tile.ID) / tileXCount) * (tileSize + ySpacing)
				opts.GeoM.Translate(float64(x*tileSize)-g.Camera.X, float64(y*tileSize)-g.Camera.Y)
				screen.DrawImage(
					g.FloorImg.SubImage(
						image.Rect(srcX, srcY, srcX+tileSize, srcY+tileSize)).(*ebiten.Image),
					&opts,
				)
				opts.GeoM.Reset()
				opts.GeoM.Scale(scale, scale)
			}
		}
	}

	// 计算角色的边界框
	spriteRect := image.Rect(int(g.Sprite.X), int(g.Sprite.Y), int(g.Sprite.X)+16, int(g.Sprite.Y)+16)
	greenColor := color.RGBA{0, 255, 0, 255}
	// 绘制角色边界框
	vector.DrawFilledRect(screen,
		float32(float64(spriteRect.Min.X)-g.Camera.X), float32(float64(spriteRect.Min.Y)-g.Camera.Y),
		float32(float64(spriteRect.Dx())), float32(float64(spriteRect.Dy())), greenColor, true)
	opts.GeoM.Reset()

	opts.GeoM.Translate(g.Sprite.X-g.Camera.X, g.Sprite.Y-g.Camera.Y)
	screen.DrawImage(
		g.Sprite.Img.SubImage(
			image.Rect(0, 0, 16, 16),
		).(*ebiten.Image),
		&opts,
	)
	opts.GeoM.Reset()

	// Display FPS.
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()))
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.screenWidth, g.screenHeight

}
func main() {
	ebiten.SetWindowSize(g.screenWidth, g.screenHeight)
	ebiten.SetWindowTitle("SuperMario")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
