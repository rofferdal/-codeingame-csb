package main

import (
	"fmt"
	"math"
	"os"
)

type SmartVector struct {
	x, y         float64
	length       float64
	angleDegrees float64
	angleRadians float64
}

func NewSmartVectorCartesian(x, y float64) SmartVector {
	angleRadians := cartesianToRadian(x, y)
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
	angleRadians := angleDegrees * math.Pi / 180
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

type gameState struct {
	first     bool
	usedboost bool
	prevx     int
	prevy     int
}

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/

func main() {
	state := gameState{
		first:     true,
		usedboost: false,
		prevx:     0,
		prevy:     0,
	}
	for {
		// nextCheckpointX: x position of the next check point
		// nextCheckpointY: y position of the next check point
		// nextCheckpointDist: distance to the next checkpoint
		// nextCheckpointAngle: angle between your pod orientation and the direction of the next checkpoint
		var x, y, nextCheckpointX, nextCheckpointY, nextCheckpointDist, nextCheckpointAngle int
		fmt.Scan(&x, &y, &nextCheckpointX, &nextCheckpointY, &nextCheckpointDist, &nextCheckpointAngle)

		var opponentX, opponentY int
		fmt.Scan(&opponentX, &opponentY)

		if state.first {
			state.prevx = x
			state.prevy = y
		}
		toOpponentV := NewSmartVectorCartesian(float64(opponentX-x), float64(opponentY-y))
		targetV := NewSmartVectorCartesian(float64(nextCheckpointX-x), float64(nextCheckpointY-y))
		lastMoveV := NewSmartVectorCartesian(float64(x-state.prevx), float64(y-state.prevy))
		if lastMoveV.length < 10 {
			lastMoveV = targetV
		}
		fmt.Fprintf(os.Stderr, "nextCheckpointAngle: %d\n", nextCheckpointAngle)
		if math.Abs(float64(nextCheckpointAngle)) < 20 {
			desiredAngle := targetV.angleDegrees
			deltaAngle := desiredAngle - lastMoveV.angleDegrees
			fmt.Fprintf(os.Stderr, "deltaAngle: %f, lastMoveV.angleDegrees: %f\n", deltaAngle, lastMoveV.angleDegrees)
			newTargetAngle := desiredAngle + (float64(deltaAngle))
			targetV = NewSmartVectorPolar(targetV.length, newTargetAngle)
			fmt.Fprintf(os.Stderr, "desiredAngle: %f, newTargetAngle: %f\n", desiredAngle, newTargetAngle)
			fmt.Fprintf(os.Stderr, "nextCheckpointX: %d, nextCheckpointX: %d,\n", nextCheckpointX, nextCheckpointY)
			fmt.Fprintf(os.Stderr, "targetV.x: %d, targetV.y: %d\n", int(targetV.x), int(targetV.y))
			fmt.Fprintf(os.Stderr, "nextx: %d, nexty: %d\n", x+int(targetV.x), y+int(targetV.y))
		}

		// fmt.Fprintln(os.Stderr, "Debug messages...")
		thrust := 100
		if nextCheckpointDist < 2000 {
			thrust = 100 * (nextCheckpointDist + 100) / 2100
			fmt.Fprintln(os.Stderr, "distancethrust:", thrust)
		}
		if nextCheckpointAngle > 90 || nextCheckpointAngle < -90 {
			thrust = 1
		}
		// You have to output the target position
		// followed by the power (0 <= thrust <= 100) or "BOOST"
		// i.e.: "x y thrust"
		targetX, targetY := targetV.GetXYAsInts()
		fmt.Fprintf(os.Stderr, "usedboost: %t", state.usedboost)
		useboost := !state.usedboost && nextCheckpointDist > 4500 && nextCheckpointAngle < 5 && nextCheckpointAngle > -5 && toOpponentV.length > 2500
		useshield := nextCheckpointDist+int(toOpponentV.length) < 2000
		if useboost {
			fmt.Printf("%d %d BOOST\n", x+targetX, y+targetY)
			state.usedboost = true
		} else if useshield {
			fmt.Printf("%d %d SHIELD\n", x+targetX, y+targetY)
		} else {
			fmt.Printf("%d %d %d\n", x+targetX, y+targetY, thrust)
		}
		state.prevx = x
		state.prevy = y
		state.first = false
	}
}
