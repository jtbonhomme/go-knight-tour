package knight

import (
	"log"
	"math/rand"
)

var accessibility = [8 * 8]int{
	2, 3, 4, 4, 4, 4, 3, 2,
	3, 4, 6, 6, 6, 6, 4, 3,
	4, 6, 8, 8, 8, 8, 6, 4,
	4, 6, 8, 8, 8, 8, 6, 4,
	4, 6, 8, 8, 8, 8, 6, 4,
	4, 6, 8, 8, 8, 8, 6, 4,
	3, 4, 6, 6, 6, 6, 4, 3,
	2, 3, 4, 4, 4, 4, 3, 2}

type Knight struct {
	Positions      []Position
	slowMotion          int
	tour           int
	implementation string
	grid           [8 * 8]int
	slowMotionChange    chan int
}

func New(slowMotion int, implementation string, slowMotionChange chan int) *Knight {
	return &Knight{
		Positions:      []Position{},
		slowMotion:          slowMotion,
		implementation: implementation,
		grid:           accessibility,
		slowMotionChange:    slowMotionChange,
	}
}

func (k *Knight) Tour() int {
	return k.tour
}

func (k *Knight) Update(positions []Position, tour int) {
	k.Positions = positions
	k.tour = tour
}

func (k *Knight) Run(stopChannel chan struct{}) chan bool {
	var result bool
	c := make(chan bool)

	// pick initial position
	x := rand.Intn(8)
	y := rand.Intn(8)

	k.Positions = append(k.Positions, Position{x, y})
	k.tour = 1

	go func() {
		log.Println("start knight's tour solver: starting from ", Position{x, y})
		switch k.implementation {
		case "naive":
			result = k.NaiveSolver(k.tour, k.Positions, stopChannel)
		case "backtracking":
			result = k.BacktrackingSolver(k.tour, k.Positions, stopChannel)
		case "optimized":
			result = k.OptimizedSolver(k.tour, k.Positions, stopChannel)
		default:
			log.Fatalf("%s implementation does not exist", k.implementation)
		}

		log.Printf("solver result: %v at position %v", result, k.Positions[len(k.Positions)-1])
		c <- result
	}()

	log.Printf("knight: run")
	return c
}

func (p Position) Distance(q Position) int {
	dx := q.X - p.X
	dy := q.Y - p.Y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}
