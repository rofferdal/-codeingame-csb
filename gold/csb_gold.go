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
		return NewSmartVectorCartesian(0, 0)
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
		return NewSmartVectorCartesian(sv.x*factor, sv.y*factor)
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
	numlaps        int
	numcheckpoints int
	first          bool
	usedboost      bool
	lastlap        bool
	players        [2]gamer
	opponents      [2]gamer
}

func initGameState(track map[int]*checkpoint, numlaps int) gameState {
	state := gameState{
		numlaps:        numlaps,
		numcheckpoints: len(track),
		first:          true,
		usedboost:      false,
		lastlap:        false,
		players:        [2]gamer{gamer{0, 0, 0, 0, 0, 0, 0, 1}, gamer{0, 0, 0, 0, 0, 0, 0, 1}},
		opponents:      [2]gamer{gamer{0, 0, 0, 0, 0, 0, 0, 1}, gamer{0, 0, 0, 0, 0, 0, 0, 1}},
	}
	return state
}

type gamer struct {
	x, y, vx, vy, angle, nextCheckPointId, advancement, currentlap int
}

func (g gamer) currentSpeedV() (SmartVector) {
	return NewSmartVectorCartesian(float64(g.vx), float64(g.vy))
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

func futureCollisionCourse(main, candidate gamer, dangerzonedist, lookaheadsteps int) (int, SmartVector) {
	mainV := NewSmartVectorCartesian(float64(main.x), float64(main.y))
	mainSpeedV := main.currentSpeedV()
	candV := NewSmartVectorCartesian(float64(candidate.x), float64(candidate.y))
	candSpeedV := candidate.currentSpeedV()
	for step := 1; step <= lookaheadsteps; step++ {
		mainStepV := mainV.addVector(mainSpeedV.multiplyNumber(float64(step)))
		candStepV := candV.addVector(candSpeedV.multiplyNumber(float64(step)))
		distanceV := mainStepV.subtractVector(candStepV)
		if distanceV.length <= float64(dangerzonedist) {
			return step, mainSpeedV.multiplyNumber(float64(step))
		}
	}
	return 0, NewSmartVectorCartesian(0, 0)
}

func movePlayer(playerId int, isLeader bool, state gameState, track map[int]*checkpoint) bool {
	player := state.players[playerId]
	checkpoint := track[player.nextCheckPointId]
	var partner gamer
	if playerId == 0 {
		partner = state.players[1]
	} else {
		partner = state.players[0]
	}
	opponents := state.opponents
	x := player.x
	y := player.y
	nextCheckpointX := checkpoint.center.x
	nextCheckpointY := checkpoint.center.y
	//currentSpeedV := player.currentSpeedV()
	toOpponent0V := NewSmartVectorCartesian(float64(opponents[0].x-x), float64(opponents[0].y-y))
	toOpponent1V := NewSmartVectorCartesian(float64(opponents[1].x-x), float64(opponents[1].y-y))
	toCheckpointV := NewSmartVectorCartesian(float64(nextCheckpointX-x), float64(nextCheckpointY-y))
	toNextAimpointV := NewSmartVectorCartesian(float64(checkpoint.nextAimpoint.x-x), float64(checkpoint.nextAimpoint.y-y))
	targetV := toCheckpointV
	toLongDistanceAimV := NewSmartVectorCartesian(float64(checkpoint.longDistanceAimpoint.x-x), float64(checkpoint.longDistanceAimpoint.y-y))

	fmt.Fprintf(os.Stderr, "player: %+v\n", player)
	nextCheckpointAngle := normalizeAngleDegrees(int(toCheckpointV.angleDegrees) - player.angle)
	nextCheckpointDist := int(targetV.length)

	var thrust int
	var useShield bool
	useBoost := false
	firstStretch := (player.currentlap == 1 && player.nextCheckPointId == 1)
	thirdLap := partner.currentlap >= 3 || opponents[0].currentlap >= 3 || opponents[1].currentlap >= 3

	neverAgressive := false // true value only for debugging
	if isLeader || neverAgressive || firstStretch {
		targetV, thrust = normalMove(player, toCheckpointV, toLongDistanceAimV, toNextAimpointV)
		useShield = nextCheckpointDist < 1000 && (toOpponent0V.length < 900 || toOpponent1V.length < 900)
		if useShield {
			useShield = shouldRarelyUseShield(player, opponents, false)
		}
		useBoost = (state.first && isLeader) || (!state.usedboost && nextCheckpointDist > 5500 && nextCheckpointAngle < 3 && nextCheckpointAngle > -3 && toOpponent0V.length > 2000 && toOpponent1V.length > 2000)
	} else if opponentLeads(state.players, state.opponents) || thirdLap {
		targetV, thrust = fullDefenseMode(player, partner, track, opponents)
		useShield = shouldUseShield(player, opponents, useShield)
	} else {
		targetV, thrust = aggroMove(player, nextCheckpointAngle, targetV, toLongDistanceAimV, toNextAimpointV, nextCheckpointDist, toOpponent0V, toOpponent1V)
		useShield = shouldUseShield(player, opponents, useShield)
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

func shouldUseShield(player gamer, opponents [2]gamer, useShield bool) bool {
	oppCollSteps0, _ := futureCollisionCourse(player, opponents[0], 800, 2)
	oppCollSteps1, _ := futureCollisionCourse(player, opponents[1], 800, 2)
	useShield = useShield || oppCollSteps0 > 0 || oppCollSteps1 > 0
	return useShield
}

func shouldRarelyUseShield(player gamer, opponents [2]gamer, useShield bool) bool {
	oppCollSteps0, _ := futureCollisionCourse(player, opponents[0], 800, 1)
	oppCollSteps1, _ := futureCollisionCourse(player, opponents[1], 800, 1)
	useShield = useShield || oppCollSteps0 > 0 || oppCollSteps1 > 0
	return useShield
}

func normalMove(player gamer, toCheckpointV, toLongDistanceAimV, toNextAimpointV SmartVector) (SmartVector, int) {
	fmt.Fprintln(os.Stderr, "Normal player")

	smartDirectionV, smartThrust := getDirectionSmartVector(player, toCheckpointV, toLongDistanceAimV, toNextAimpointV)

	thrust := 100
	if smartThrust < 0 {
		checkpointDeltaAngle := normalizeAngleDegrees(int(toCheckpointV.angleDegrees) - player.angle)
		checkpointDist := toCheckpointV.length
		if checkpointDist < 2000 {
			thrust = int(100 * (checkpointDist + 100) / 2100)
			fmt.Fprintf(os.Stderr, "distance: %f, thrust: %f\n", checkpointDist, thrust)
		}
		if checkpointDeltaAngle > 45 || checkpointDeltaAngle < -45 {
			thrust = 60
		}
		if checkpointDeltaAngle > 90 || checkpointDeltaAngle < -90 {
			thrust = 1
		}
	} else {
		thrust = smartThrust
	}
	return smartDirectionV, thrust
}

func aggroMove(player gamer, nextCheckpointAngle int, defaultTargetV, longDistanceAimV, toNextAimpointV SmartVector, nextCheckpointDist int, toOpponent0V, toOpponent1V SmartVector) (SmartVector, int) {
	fmt.Fprintln(os.Stderr, "AGGRO MODE")
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
	thrust := 100
	if !aggressive {
		var smartThrust int
		aggroTargetV, smartThrust = getDirectionSmartVector(player, defaultTargetV, longDistanceAimV, toNextAimpointV)
		if smartThrust < 0 {
			if nextCheckpointDist < 2000 {
				thrust = 100 * (nextCheckpointDist + 100) / 2100
				fmt.Fprintf(os.Stderr, "distance: %f, thrust: %f\n", nextCheckpointDist, thrust)
			}
			if nextCheckpointAngle > 90 || nextCheckpointAngle < -90 {
				thrust = 5
			}
		} else {
			thrust = smartThrust
		}
	}

	return aggroTargetV, thrust
}

func fullDefenseMode(player, partner gamer, track map[int]*checkpoint, opponents [2]gamer) (SmartVector, int) {
	fmt.Fprintln(os.Stderr, "FULL DEFENSE MODE")
	x := player.x
	y := player.y
	currentSpeedV := player.currentSpeedV()
	opponentLeaderId := determineLeader(opponents)
	opponentLeader := opponents[opponentLeaderId]
	opponentCheckpoint := track[opponentLeader.nextCheckPointId].center
	opponentNextCheckpoint := track[opponentLeader.nextCheckPointId].nextAimpoint
	toOpponentV := NewSmartVectorCartesian(float64(opponentLeader.x-x), float64(opponentLeader.y-y))
	toOpponentTargetV := NewSmartVectorCartesian(float64(opponentCheckpoint.x-x), float64(opponentCheckpoint.y-y))
	opponentToTargetV := NewSmartVectorCartesian(float64(opponentCheckpoint.x-opponentLeader.x), float64(opponentCheckpoint.y-opponentLeader.x))
	toOpponentNextTargetV := NewSmartVectorCartesian(float64(opponentNextCheckpoint.x-x), float64(opponentNextCheckpoint.y-y))

	distanceToOpponent := toOpponentV.length
	magicAngleRadians := (180 - float64(normalizeAngleDegrees(int(toOpponentV.angleDegrees-opponentToTargetV.angleDegrees)))) / 180 * math.Pi
	distanceToIntersection := math.Abs(distanceToOpponent / 2 * math.Tan(magicAngleRadians))
	fmt.Fprintf(os.Stderr, "magicAngle: %f, distToOpp: %f, distToInt: %f\n", magicAngleRadians, distanceToOpponent, distanceToIntersection)
	opponentToIntersectionV := NewSmartVectorPolar(distanceToIntersection, opponentToTargetV.angleDegrees)
	toIntersectionV := opponentToIntersectionV.addVector(toOpponentV)

	aggroTargetV := toIntersectionV
	if aggroTargetV.length > opponentToTargetV.length {
		if opponentLeader.nextCheckPointId > 0 || opponentLeader.currentlap < 3 {
			aggroTargetV = toOpponentNextTargetV
			if toOpponentNextTargetV.length < 1200 {
				aggroTargetV = toOpponentV
			}
		} else {
			//aggroTargetV = toOpponentTargetV
			if toOpponentTargetV.length < 1200 {
				aggroTargetV = toOpponentV
			}
		}
	}

	collsteps, collisionV := futureCollisionCourse(player, partner, 800, 10)
	if collsteps > 0 && collisionV.length < aggroTargetV.length {
		diffAngle := normalizeAngleDegrees(int(aggroTargetV.angleDegrees - collisionV.angleDegrees))
		if diffAngle > 0 {
			aggroTargetV = NewSmartVectorPolar(aggroTargetV.length, collisionV.angleDegrees+30)
		} else {
			aggroTargetV = NewSmartVectorPolar(aggroTargetV.length, collisionV.angleDegrees-30)
		}
		fmt.Fprintf(os.Stderr, "COLLISIONPATH  V: %v\n", collisionV)
		fmt.Fprintf(os.Stderr, "COLLISIONAVOID V: %v\n", aggroTargetV)
	}

	if math.Abs(float64(normalizeAngleDegrees(int(aggroTargetV.angleDegrees - currentSpeedV.angleDegrees)))) < 40 {
		aggroTargetV = smartDirectionChangeVector(aggroTargetV, currentSpeedV)
	}

	thrust := 100

	targetAngle := normalizeAngleDegrees(int(aggroTargetV.angleDegrees) - player.angle)
	if targetAngle > 90 || targetAngle < -90 {
		thrust = 1
	} else if targetAngle > 60 || targetAngle < -60 {
		thrust = 10
	} else if targetAngle > 30 || targetAngle < -30 {
		thrust = 40
	}
	return aggroTargetV, thrust
}

func getDirectionSmartVector(player gamer, toCheckpointV, longDistanceAimV, toNextAimpointV SmartVector) (SmartVector, int) {
	smartThrust := -1
	smartDirectionV := toCheckpointV
	currentSpeedV := player.currentSpeedV()
	checkpointAngle := toCheckpointV.angleDegrees
	checkpointDeltaAngle := normalizeAngleDegrees(int(toCheckpointV.angleDegrees) - player.angle)
	viabilityAngle := normalizeAngleDegrees(int(longDistanceAimV.angleDegrees - toCheckpointV.angleDegrees))
	playerAngleVsNextAimpointAngle := normalizeAngleDegrees(player.angle - int(toNextAimpointV.angleDegrees))
	turnRounds := math.Abs(float64(playerAngleVsNextAimpointAngle / 18))
	roundsToTargetCurrentSpeedWFriction, distWFriction := toTargetGivenSpeedWFriction(toCheckpointV.length, currentSpeedV.length)
	//roundsToTargetCurrentSpeed := toCheckpointV.length / currentSpeedV.length
	//dist := toCheckpointV.length
	predictedPathV := NewSmartVectorPolar(float64(distWFriction), currentSpeedV.angleDegrees)
	willProbablyHit := (toCheckpointV.subtractVector(predictedPathV)).length < 500
	fmt.Fprintf(os.Stderr, "RTTCSWF: %f, TRDS: %f, WPH: %t\n", roundsToTargetCurrentSpeedWFriction, turnRounds, willProbablyHit)
	fmt.Fprintf(os.Stderr, "predictedPathV: %v\n", predictedPathV)
	fmt.Fprintf(os.Stderr, "toCheckpointV: %v\n", toCheckpointV)
	fmt.Fprintf(os.Stderr, "currentSpeedV: %v\n", currentSpeedV)
	if math.Abs(float64(viabilityAngle)) < 45 && toCheckpointV.length > 5500 {
		smartDirectionV = longDistanceAimV
		fmt.Fprintf(os.Stderr, "USING SMARTDIRECTION: %+v\n", smartDirectionV)
    } else if roundsToTargetCurrentSpeedWFriction >= 1 && turnRounds >= 1 && willProbablyHit {
		smartDirectionV = toNextAimpointV
		smartThrust = 0
		potentialSpeedChangeV := NewSmartVectorPolar(100, float64(player.angle + (playerAngleVsNextAimpointAngle / int(turnRounds))))
		newSpeedV := currentSpeedV.addVector(potentialSpeedChangeV).multiplyNumber(0.85)
		_, newDistWFriction := toTargetGivenSpeedWFriction(toCheckpointV.length, newSpeedV.length)
		predictedNewPathV := NewSmartVectorPolar(float64(newDistWFriction), newSpeedV.angleDegrees)
		willProbablyHitNew := (toCheckpointV.subtractVector(predictedNewPathV)).length < 500
		if willProbablyHitNew {
			smartThrust = 100
		}
		fmt.Fprintln(os.Stderr, "Cut the curve with thrust: %d", smartThrust)
		smartDirectionV = toNextAimpointV
	} else if toCheckpointV.length > 1500 && (math.Abs(float64(checkpointDeltaAngle)) < 20 || (toCheckpointV.length < 2000 && math.Abs(float64(checkpointAngle)) < 45)) {
		smartDirectionV = smartDirectionChangeVector(toCheckpointV, currentSpeedV)
	} else if (toCheckpointV.length < 1500) && (math.Abs(float64(checkpointDeltaAngle)) < 10) {
		fmt.Fprintln(os.Stderr, "Oh so close, target next")
		smartDirectionV = toNextAimpointV
	}
	return smartDirectionV, smartThrust
}

func smartDirectionChangeVector(targetV SmartVector, currentSpeedV SmartVector) SmartVector {
	desiredAngle := targetV.angleDegrees
	deltaAngle := normalizeAngleDegrees(int(desiredAngle - currentSpeedV.angleDegrees))
	fmt.Fprintf(os.Stderr, "deltaAngle: %f, lastMoveV.angleDegrees: %f\n", deltaAngle, currentSpeedV.angleDegrees)
	newTargetAngle := desiredAngle + (float64(deltaAngle))
	smartDirectionV := NewSmartVectorPolar(targetV.length, newTargetAngle)
	fmt.Fprintf(os.Stderr, "desiredAngle: %f, newTargetAngle: %f\n", desiredAngle, newTargetAngle)
	fmt.Fprintf(os.Stderr, "smartDirectionV.x: %d, smartDirectionV.y: %d\n", int(smartDirectionV.x), int(smartDirectionV.y))
	return smartDirectionV
}

func toTargetGivenSpeedWFriction(targetLength, speedLength float64) (float64, float64) {
	roundsToTarget := 0
	distance := 0.0
	for round := 1; round < 10; round++ {
		distance = distance + speedLength * math.Pow(0.85, float64(round))
		if (distance >= targetLength) {
			roundsToTarget = round
			return float64(roundsToTarget), distance
		}
	}
	return 0, 0
}

func readPlayers(state gameState, track map[int]*checkpoint) [2]gamer {
	var players [2]gamer
	for i := 0; i < 2; i++ {
		var x, y, vx, vy, angle, nextCheckPointId int
		fmt.Scan(&x, &y, &vx, &vy, &angle, &nextCheckPointId)
		players[i] = gamer{x, y, vx, vy, angle, nextCheckPointId, state.players[i].advancement, state.players[i].currentlap}
		if state.players[i].nextCheckPointId != nextCheckPointId {
			// new checkpoint
			fmt.Fprintf(os.Stderr, "NEW nextCheckPointId %d for player %d\n", nextCheckPointId, i)
			if nextCheckPointId == 0 {
				players[i].currentlap = players[i].currentlap + 1
			}
		}
		toCheckPointV := NewSmartVectorCartesian(float64(track[nextCheckPointId].center.x-x), float64(track[nextCheckPointId].center.y-y))
		players[i].advancement = players[i].currentlap*1000000 + players[i].nextCheckPointId*100000 - int(toCheckPointV.length)
	}
	return players
}

func determineLeader(players [2]gamer) int {
	leaderId := 0
	for i := 0; i < 2; i++ {
		if players[i].advancement > players[leaderId].advancement {
			leaderId = i
		}
	}
	return leaderId
}

func opponentLeads(players [2]gamer, opponents [2]gamer) bool {
	playerLeadId := determineLeader(players)
	opponentLeadId := determineLeader(opponents)

	return players[playerLeadId].advancement < opponents[opponentLeadId].advancement
}

func readOpponents(state gameState, track map[int]*checkpoint) [2]gamer {
	var opponents [2]gamer
	for i := 0; i < 2; i++ {
		var x2, y2, vx2, vy2, angle2, nextCheckPointId2 int
		fmt.Scan(&x2, &y2, &vx2, &vy2, &angle2, &nextCheckPointId2)
		opponents[i] = gamer{x2, y2, vx2, vy2, angle2, nextCheckPointId2, state.opponents[i].advancement, state.opponents[i].currentlap}
		if state.opponents[i].nextCheckPointId != nextCheckPointId2 {
			// new checkpoint
			if nextCheckPointId2 == 0 {
				opponents[i].currentlap = opponents[i].currentlap + 1
			}
		}
		toCheckPointV := NewSmartVectorCartesian(float64(track[nextCheckPointId2].center.x-x2), float64(track[nextCheckPointId2].center.y-y2))
		opponents[i].advancement = opponents[i].currentlap*1000000 + opponents[i].nextCheckPointId*100000 - int(toCheckPointV.length)
	}
	return opponents
}

func readTrack() map[int]*checkpoint {
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
	return track
}

func main() {
	var laps int
	fmt.Scan(&laps)
	track := readTrack()

	state := initGameState(track, laps)

	for {
		players := readPlayers(state, track)
		state.players = players

		leaderId := determineLeader(players)
		opponents := readOpponents(state, track)
		state.opponents = opponents

		for playerId := 0; playerId < 2; playerId++ {
			isLeader := playerId == leaderId

			boostUsed := movePlayer(playerId, isLeader, state, track)

			if boostUsed && !state.first {
				state.usedboost = true
			}
		}

		state.first = false
	}
}
