package common

import (
	"encoding/json"
	"fmt"
	"os"
)

// Tilemap 结构体用于解析 tilemap.tmj 文件
type TileMap struct {
	CompressionLevel int        `json:"compressionlevel"`
	Height           int        `json:"height"`
	Infinite         bool       `json:"infinite"`
	Layers           []*Layer   `json:"layers"`
	NextLayerID      int        `json:"nextlayerid"`
	NextObjectID     int        `json:"nextobjectid"`
	Orientation      string     `json:"orientation"`
	RenderOrder      string     `json:"renderorder"`
	TiledVersion     string     `json:"tiledversion"`
	TileHeight       int        `json:"tileheight"`
	Tilesets         []*Tileset `json:"tilesets"`
	TileWidth        int        `json:"tilewidth"`
	Type             string     `json:"type"`
	Version          string     `json:"version"`
	Width            int        `json:"width"`
}

// Layer 结构体用于解析 layers 数组中的每个图层
type Layer struct {
	Data    []int   `json:"data"`
	Height  int     `json:"height"`
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Opacity float64 `json:"opacity"`
	Type    string  `json:"type"`
	Visible bool    `json:"visible"`
	Width   int     `json:"width"`
	X       int     `json:"x"`
	Y       int     `json:"y"`
}

// Tileset 结构体用于解析 tilesets 数组中的每个瓦片集
type Tileset struct {
	Columns    int     `json:"columns"`
	FirstGID   int     `json:"firstgid"`
	Grid       *Grid   `json:"grid"`
	Margin     int     `json:"margin"`
	Name       string  `json:"name"`
	Spacing    int     `json:"spacing"`
	TileCount  int     `json:"tilecount"`
	TileHeight int     `json:"tileheight"`
	Tiles      []*Tile `json:"tiles"`
	TileWidth  int     `json:"tilewidth"`
}

// Grid 结构体用于解析 grid 对象
type Grid struct {
	Height      int    `json:"height"`
	Orientation string `json:"orientation"`
	Width       int    `json:"width"`
}

// Tile 结构体用于解析 tiles 数组中的每个瓦片
type Tile struct {
	ID          int    `json:"id"`
	Image       string `json:"image"`
	ImageHeight int    `json:"imageheight"`
	ImageWidth  int    `json:"imagewidth"`
}

func LoadConfig(filename string) (*TileMap, error) {
	if contents, err := os.ReadFile(filename); err != nil {
		return nil, err
	} else {
		m := &TileMap{}
		err = json.Unmarshal(contents, &m)
		if err != nil {
			return nil, err
		}
		fmt.Println(m)
		return m, nil
	}
}