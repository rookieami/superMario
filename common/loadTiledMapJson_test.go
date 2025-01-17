package common

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {

	gmap, err := LoadConfig("../resource/map/word1.tmj")
	if err != nil {
		t.Error(err)
		return
	}
	for _, TileSet := range gmap.TileSets {
		for _, tile := range TileSet.Tiles {
			if tile.GetPropertyBool("CanPassed") {
				t.Log("tile ", tile.ID)
			}
		}
	}
}
