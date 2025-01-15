package entities

import "math"

type Camera struct {
	X float64
	Y float64
}

func (c *Camera) FollowTarget(targetX, targetY, screenWidth, screenHeight float64) {
	c.X = targetX - screenWidth/4
	c.Y = targetY - screenHeight/2
}
func (c *Camera) Constrain(mapWidth, mapHeight, screenWidth, screenHeight float64) {
	c.X = math.Max(c.X, 0)
	c.Y = math.Max(c.Y, 0)
	c.X = math.Min(c.X, mapWidth-screenWidth)
	c.Y = math.Min(c.Y, mapHeight-screenHeight)
}
