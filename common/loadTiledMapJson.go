package common

import (
	"encoding/json"
	"os"
)

// TileMap represents the structure of the tiled map JSON.
type TileMap struct {
	Compressionlevel int        `json:"compressionlevel"`
	Version          string     `json:"version"`
	TiledVersion     string     `json:"tiledversion"`
	Orientation      string     `json:"orientation"`
	RenderOrder      string     `json:"renderorder"`
	Width            int        `json:"width"`
	Height           int        `json:"height"`
	TileWidth        int        `json:"tilewidth"`
	TileHeight       int        `json:"tileheight"`
	Infinite         bool       `json:"infinite"`
	NextLayerID      int        `json:"nextlayerid"`
	NextObjectID     int        `json:"nextobjectid"`
	TileSets         []*TileSet `json:"tilesets"`
	Layers           []*Layer   `json:"layers"`
	Tiles            map[int]*Tile
}

// TileSet represents the structure of a tileset in the tiled map JSON.
type TileSet struct {
	FirstGID       int         `json:"firstgid"`
	Columns        int         `json:"columns"`
	FillMode       string      `json:"fillmode"`
	Image          string      `json:"image"`
	ImageHeight    int         `json:"imageheight"`
	ImageWidth     int         `json:"imagewidth"`
	Margin         int         `json:"margin"`
	Name           string      `json:"name"`
	Properties     []*Property `json:"properties"`
	Spacing        int         `json:"spacing"`
	TileCount      int         `json:"tilecount"`
	TileRenderSize string      `json:"tilerendersize"`
	Tiles          []*Tile     `json:"tiles"`
	Tilewidth      int         `json:"tilewidth"`
	TileHeight     int         `json:"tileheight"`
}
type Property struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value bool   `json:"value"`
}
type Tile struct {
	ID         int `json:"id"`
	Properties []*Property
}

func (t *Tile) GetPropertyBool(name string) bool {
	for _, property := range t.Properties {
		if property.Name == name && property.Type == "bool" {
			return property.Value
		}
	}
	return false
}

// Layer represents the structure of a layer in the tiled map JSON.
type Layer struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Opacity int    `json:"opacity"`
	Type    string `json:"type"`
	Visible bool   `json:"visible"`
	X       int    `json:"x"`
	Y       int    `json:"y"`
	Data    []int  `json:"data"`
}

// LoadConfig reads and parses the JSON configuration file.
func LoadConfig(filePath string) (*TileMap, error) {
	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// 解析JSON数据
	tileMap := &TileMap{}
	err = json.Unmarshal(data, tileMap)
	if err != nil {
		return nil, err
	}
	tileMap.Tiles = make(map[int]*Tile)
	for _, tileSet := range tileMap.TileSets {
		for _, tile := range tileSet.Tiles {
			tileMap.Tiles[tile.ID] = tile
		}
	}
	return tileMap, nil
}
