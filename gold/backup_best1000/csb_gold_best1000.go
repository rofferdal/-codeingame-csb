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

type checkpoint struct {
	center               point
	longDistanceAimpoint point
	nextAimpoint         point
}
type point struct {
	x int
	y int
}

type gameState struct {
	first                  bool
	usedboost              bool
	playerpos              [2]point
	playeradvance          [2]int
	playerprevcheckpointid [2]int
}

func initGameState(track map[int]*checkpoint) gameState {
	state := gameState{
		first:                  true,
		usedboost:              false,
		playerpos:              [2]point{track[0].center, track[0].center},
		playeradvance:          [2]int{0, 0},
		playerprevcheckpointid: [2]int{0, 0},
	}
	return state
}

type gamer struct {
	x, y, vx, vy, angle, nextCheckPointId, advancement int
}

func calculateAimpoints(track map[int]*checkpoint) {
	nextpoint := track[0]
	for id := len(track) - 1; id >= 0; id-- {
		currpoint := track[id]
		ldaX := currpoint.center.x + ((currpoint.center.x - nextpoint.center.x) / 3)
		ldaY := currpoint.center.y + ((currpoint.center.y - nextpoint.center.y) / 3)
		track[id].longDistanceAimpoint.x = ldaX
		track[id].longDistanceAimpoint.y = ldaY
		track[id].nextAimpoint = nextpoint.center
		nextpoint = currpoint
	}
	for id := 0; id < len(track); id++ {
		fmt.Fprintf(os.Stderr, "checkpoint %d: %+v\n", id, *track[id])
	}
}

func normalizeAngleDegrees(angle int) int {
	// Normalize to game standard between -180 and 180 degrees
	if angle > 180 {
		return angle - 360
	}
	if angle < -180 {
		return angle + 360
	}
	return angle
}

func movePlayer(playerid int, isLeader bool, players [2]gamer, opponents [2]gamer, checkpoint *checkpoint, state gameState) bool {
	player := players[playerid]
	x := player.x
	y := player.y
	nextCheckpointX := checkpoint.center.x
	nextCheckpointY := checkpoint.center.y
	toOpponent0V := NewSmartVectorCartesian(float64(opponents[0].x-x), float64(opponents[0].y-y))
	toOpponent1V := NewSmartVectorCartesian(float64(opponents[1].x-x), float64(opponents[1].y-y))
	checkpointV := NewSmartVectorCartesian(float64(nextCheckpointX-x), float64(nextCheckpointY-y))
	theCheckpointAfterV := NewSmartVectorCartesian(float64(checkpoint.nextAimpoint.x-x), float64(checkpoint.nextAimpoint.y-y))
	targetV := checkpointV
	longDistanceAimV := NewSmartVectorCartesian(float64(checkpoint.longDistanceAimpoint.x-x), float64(checkpoint.longDistanceAimpoint.y-y))
	lastMoveV := NewSmartVectorCartesian(float64(x-state.playerpos[playerid].x), float64(y-state.playerpos[playerid].y))
	if lastMoveV.length < 5 {
		lastMoveV = checkpointV
	}
	fmt.Fprintf(os.Stderr, "player: %+v\n", player)
	fmt.Fprintf(os.Stderr, "targetV: %+v\n", targetV)
	nextCheckpointAngle := normalizeAngleDegrees(int(checkpointV.angleDegrees) - player.angle)
	fmt.Fprintf(os.Stderr, "nextCheckpointAngle: %d\n", nextCheckpointAngle)
	nextCheckpointDist := int(targetV.length)

	var thrust int
	var useShield bool
	useBoost := false
	if isLeader || player.advancement < 100000 {
		targetV, thrust = normalMove(nextCheckpointAngle, targetV, longDistanceAimV, theCheckpointAfterV, lastMoveV, x, y, nextCheckpointDist)
		useShield = nextCheckpointDist < 1000 && (toOpponent0V.length < 900 || toOpponent1V.length < 900)
		useShield = useShield || ((toOpponent0V.length < 800 || toOpponent1V.length < 800) && lastMoveV.length > 20 && player.advancement > 100000)
		useBoost = (state.first && isLeader) || (!state.usedboost && nextCheckpointDist > 5000 && nextCheckpointAngle < 3 && nextCheckpointAngle > -3 && toOpponent0V.length > 2000 && toOpponent1V.length > 2000)
	} else {
		targetV, thrust = aggroMove(nextCheckpointAngle, targetV, longDistanceAimV, theCheckpointAfterV, lastMoveV, x, y, nextCheckpointDist, toOpponent0V, toOpponent1V)
		useShield = useShield || ((toOpponent0V.length < 1000 || toOpponent1V.length < 1000) && lastMoveV.length > 20)
	}
	fmt.Fprintf(os.Stderr, "targetV: %v\n", targetV)

	// You have to output the target position
	// followed by the power (0 <= thrust <= 100) or "BOOST"
	// i.e.: "x y thrust"
	targetX, targetY := targetV.GetXYAsInts()
	fmt.Fprintf(os.Stderr, "usedboost: %t\n", state.usedboost)
	if useBoost {
		fmt.Fprintf(os.Stderr, "BOOOOOOOOOOOST!!!!!!!!!!!!!!!\n")
		fmt.Printf("%d %d BOOST\n", x+targetX, y+targetY)
	} else if useShield {
		fmt.Printf("%d %d SHIELD\n", x+targetX, y+targetY)
	} else {
		fmt.Printf("%d %d %d\n", x+targetX, y+targetY, thrust)
	}
	return useBoost
}

func normalMove(nextCheckpointAngle int, targetV, longDistanceAimV, theCheckpointAfterV, lastMoveV SmartVector, x, y, nextCheckpointDist int) (SmartVector, int) {
	fmt.Fprintln(os.Stderr, "Normal player")
	smartDirectionV := targetV

	smartDirectionV = getDirectionSmartVector (nextCheckpointAngle, targetV, longDistanceAimV, theCheckpointAfterV, lastMoveV, x, y, smartDirectionV)

	thrust := 100
	if nextCheckpointDist < 2000 {
		thrust = 100 * (nextCheckpointDist + 100) / 2100
		fmt.Fprintf(os.Stderr, "distance: %f, thrust: %f\n", nextCheckpointDist, thrust)
	}
	if nextCheckpointAngle > 90 || nextCheckpointAngle < -90 {
		thrust = 5
	}
	return smartDirectionV, thrust
}

func aggroMove(nextCheckpointAngle int, defaultTargetV, longDistanceAimV, theCheckpointAfterV, lastMoveV SmartVector, x, y, nextCheckpointDist int, toOpponent0V, toOpponent1V SmartVector) (SmartVector, int) {
	fmt.Fprintln(os.Stderr, "AGGRO PLAYER!!")
	aggroTargetV := defaultTargetV
	aggressive := false
	if defaultTargetV.length < 6000 {
		if toOpponent0V.length*2 < defaultTargetV.length && toOpponent0V.length < aggroTargetV.length {
			aggroTargetV = toOpponent0V
			aggressive = true
		}
		if toOpponent1V.length*2 < defaultTargetV.length && toOpponent1V.length < aggroTargetV.length {
			aggroTargetV = toOpponent1V
			aggressive = true
		}
	}
	if !aggressive {
		aggroTargetV = getDirectionSmartVector (nextCheckpointAngle, defaultTargetV, longDistanceAimV, theCheckpointAfterV, lastMoveV, x, y, aggroTargetV)
	}

	thrust := 100
	if !aggressive {
		if nextCheckpointDist < 2000 {
			thrust = 100 * (nextCheckpointDist + 100) / 2100
			fmt.Fprintf(os.Stderr, "distance: %f, thrust: %f\n", nextCheckpointDist, thrust)
		}
		if nextCheckpointAngle > 90 || nextCheckpointAngle < -90 {
			thrust = 5
		}
	}
	return aggroTargetV, thrust
}

func getDirectionSmartVector(nextCheckpointAngle int, targetV SmartVector, longDistanceAimV SmartVector, theCheckpointAfterV SmartVector, lastMoveV SmartVector, x int, y int, smartDirectionV SmartVector) SmartVector {
	viabilityAngle := normalizeAngleDegrees(int(longDistanceAimV.angleDegrees - targetV.angleDegrees))
	if math.Abs(float64(viabilityAngle)) < 45 && targetV.length > 5500 {
		smartDirectionV = longDistanceAimV
		fmt.Fprintf(os.Stderr, "USING SMARTDIRECTION: %+v\n", smartDirectionV)
	} else if targetV.length > 1500 && (math.Abs(float64(nextCheckpointAngle)) < 20 || (targetV.length < 2000 && math.Abs(float64(nextCheckpointAngle)) < 45)) {
		desiredAngle := targetV.angleDegrees
		deltaAngle := normalizeAngleDegrees(int(desiredAngle - lastMoveV.angleDegrees))
		fmt.Fprintf(os.Stderr, "deltaAngle: %f, lastMoveV.angleDegrees: %f\n", deltaAngle, lastMoveV.angleDegrees)
		newTargetAngle := desiredAngle + (float64(deltaAngle))
		smartDirectionV = NewSmartVectorPolar(targetV.length, newTargetAngle)
		fmt.Fprintf(os.Stderr, "desiredAngle: %f, newTargetAngle: %f\n", desiredAngle, newTargetAngle)
		fmt.Fprintf(os.Stderr, "smartDirectionV.x: %d, smartDirectionV.y: %d\n", int(smartDirectionV.x), int(smartDirectionV.y))
		fmt.Fprintf(os.Stderr, "nextx: %d, nexty: %d\n", x+int(smartDirectionV.x), y+int(smartDirectionV.y))
	} else if (targetV.length < 1500) && (math.Abs(float64(nextCheckpointAngle)) < 10) {
		fmt.Fprintln(os.Stderr, "Oh so close, target next")
		smartDirectionV = theCheckpointAfterV
	}
	return smartDirectionV
}

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/

func main() {
	var laps int
	fmt.Scan(&laps)

	var checkpointCount int
	fmt.Scan(&checkpointCount)

	var track map[int]*checkpoint = make(map[int]*checkpoint)
	for id := 0; id < checkpointCount; id++ {
		var checkpointX, checkpointY int
		fmt.Scan(&checkpointX, &checkpointY)
		track[id] = &checkpoint{
			center:               point{checkpointX, checkpointY},
			longDistanceAimpoint: point{checkpointX, checkpointY},
		}
	}
	calculateAimpoints(track)
	state := initGameState(track)

	for {
		var players [2]gamer
		for i := 0; i < 2; i++ {
			// x: x position of your pod
			// y: y position of your pod
			// vx: x speed of your pod
			// vy: y speed of your pod
			// angle: angle of your pod
			// nextCheckPointId: next check point id of your pod
			var x, y, vx, vy, angle, nextCheckPointId int
			fmt.Scan(&x, &y, &vx, &vy, &angle, &nextCheckPointId)
			if state.playerprevcheckpointid[i] != nextCheckPointId {
				state.playeradvance[i] = state.playeradvance[i] + 1
				state.playerprevcheckpointid[i] = nextCheckPointId
			}
			toCheckPointV := NewSmartVectorCartesian(float64(track[nextCheckPointId].center.x-x), float64(track[nextCheckPointId].center.y-y))
			advancement := state.playeradvance[i]*100000 - int(toCheckPointV.length)
			players[i] = gamer{x, y, vx, vy, angle, nextCheckPointId, advancement}
		}
		// determine leader
		leaderId := 0
		for i := 0; i < 2; i++ {
			if players[i].advancement > players[leaderId].advancement {
				leaderId = i
			}
		}

		var opponents [2]gamer
		for i := 0; i < 2; i++ {
			// x2: x position of the opponent's pod
			// y2: y position of the opponent's pod
			// vx2: x speed of the opponent's pod
			// vy2: y speed of the opponent's pod
			// angle2: angle of the opponent's pod
			// nextCheckPointId2: next check point id of the opponent's pod
			var x2, y2, vx2, vy2, angle2, nextCheckPointId2 int
			fmt.Scan(&x2, &y2, &vx2, &vy2, &angle2, &nextCheckPointId2)
			opponents[i] = gamer{x2, y2, vx2, vy2, angle2, nextCheckPointId2, 0}
		}

		for i := 0; i < 2; i++ {
			isLeader := (i == leaderId)
			usedboost := movePlayer(i, isLeader, players, opponents, track[players[i].nextCheckPointId], state)
			state.playerpos[i].x = players[i].x
			state.playerpos[i].y = players[i].y
			if usedboost && !state.first {
				state.usedboost = true
			}
		}
		// fmt.Fprintln(os.Stderr, "Debug messages...")

		// You have to output the target position
		// followed by the power (0 <= thrust <= 100)
		// i.e.: "x y thrust"
		state.first = false
	}
}
