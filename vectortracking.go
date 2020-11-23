package main

import (
	"fmt"
	"math"
)

type SmartVector struct {
	x, y         float64
	length       float64
	angleDegrees float64
	angleRadians float64
}

func NewSmartVectorCartesian(x, y float64) SmartVector {
	angleRadians := cartesianToRadian(x, y)
	angleRadians = normalizeAngleRadian(angleRadians)
	angleDegrees := angleRadians * 180.0 / math.Pi
	smartVector := SmartVector{
		x:            x,
		y:            y,
		length:       math.Sqrt(x*x + y*y),
		angleDegrees: angleDegrees,
		angleRadians: angleRadians,
	}
	return smartVector
}

func NewSmartVectorPolar(length, angleDegrees float64) SmartVector {
	if length == 0 {
		return NewSmartVectorCartesian(0,0)
	}
	angleRadians := angleDegrees * math.Pi / 180
	angleRadians = normalizeAngleRadian(angleRadians)
	angleDegrees = angleRadians * 180.0 / math.Pi
	smartVector := SmartVector{
		x:            length * math.Cos(angleRadians),
		y:            length * math.Sin(angleRadians),
		length:       length,
		angleDegrees: angleDegrees,
		angleRadians: angleRadians,
	}
	return smartVector
}

func (sv SmartVector) GetXYAsInts() (int, int) {
	return int(sv.x), int(sv.y)
}

func (sv SmartVector) multiplyNumber(factor float64) SmartVector {
	if sv.length == 0 || math.IsNaN(sv.angleRadians) {
		return NewSmartVectorCartesian(0, 0)
	} else {
		return NewSmartVectorCartesian(sv.x * factor, sv.y * factor)
	}
}

func (sv SmartVector) addVector(otherVector SmartVector) SmartVector {
	return NewSmartVectorCartesian(sv.x+otherVector.x, sv.y+otherVector.y)
}

func (sv SmartVector) subtractVector(otherVector SmartVector) SmartVector {
	return NewSmartVectorCartesian(sv.x-otherVector.x, sv.y-otherVector.y)
}

func cartesianToRadian(x, y float64) float64 {
	angleRadians := math.Atan(y / x)
	if x < 0 && y >= 0 {
		return angleRadians + math.Pi
	} // quadrant 2
	if x < 0 && y < 0 {
		return angleRadians - math.Pi
	} // quadrant 3
	return angleRadians // default for quadrant 1 and quadrant 4
}

func normalizeAngleRadian(angle float64) float64 {
	// Normalize to game standard between -180 and 180 degrees
	if angle > math.Pi {
		return angle - (math.Pi * 2)
	}
	if angle < (math.Pi * -1) {
		return angle + (math.Pi * 2)
	}
	return angle
}

func main() {
	fmt.Print("Hello!")
}
