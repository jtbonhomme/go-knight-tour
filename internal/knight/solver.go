package knight

import (
	"log"
	"slices"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func (k *Knight) handleSpeed() {
	select {
	case s := <-k.slowMotionChange:
		log.Println("slowMotion change")
		k.slowMotion = s
	default:
	}

	time.Sleep(time.Millisecond * time.Duration(float64(1000*k.slowMotion)/ebiten.ActualTPS()))
}

func (k *Knight) NaiveSolver(tour int, positions []Position, stopChannel chan struct{}) bool {
	select {
	case <-stopChannel:
		log.Println("solver stopped")

		return false
	default:
		k.handleSpeed()

		k.tour = tour
		if tour == 8*8 {
			log.Println("win!")
			return true
		}

		moves := RandomMoves()
		// pick successively random moves
		for _, m := range moves {
			p := positions[len(positions)-1]
			p.X += m.X
			p.Y += m.Y
			if !slices.Contains(positions, p) && !k.OutOfRange(p) {
				positions = append(positions, p)
				k.Positions = positions
				return k.NaiveSolver(tour+1, positions, stopChannel)
			}
		}

		return false
	}
}

func (k *Knight) BacktrackingSolver(tour int, positions []Position, stopChannel chan struct{}) bool {
	select {
	case <-stopChannel:
		log.Println("solver stopped")

		return false
	default:
		k.handleSpeed()

		k.tour = tour
		if tour == 8*8 {
			log.Println("win!")

			return true
		}

		moves := RandomMoves()
		// pick successively random moves
		for _, m := range moves {
			p := positions[len(positions)-1]
			p.X += m.X
			p.Y += m.Y
			if !slices.Contains(positions, p) && !k.OutOfRange(p) {
				positions = append(positions, p)
				k.Positions = positions
				if k.BacktrackingSolver(tour+1, positions, stopChannel) {
					return true
				}
			}
		}

		return false
	}
}

func (k *Knight) OptimizedSolver(tour int, positions []Position, stopChannel chan struct{}) bool {
	select {
	case <-stopChannel:
		log.Println("solver stopped")

		return false
	default:
		k.handleSpeed()

		k.tour = tour
		if tour == 8*8 {
			log.Println("win!")
			return true
		}

		currentPosition := positions[len(positions)-1]
		rankedPositions := k.RankedPositions(currentPosition)
		// pick best moves
		for _, p := range rankedPositions {
			if k.IsValid(p) {
				positions = append(positions, p)
				k.Invalidate(p)
				k.Positions = positions
				if k.OptimizedSolver(tour+1, positions, stopChannel) {
					return true
				}
			}
		}

		return false
	}
}
