package doric_test

import (
	"math/rand"
	"time"

	"github.com/svera/doric"
)

func Example() {
	cfg := doric.Config{
		NumberTilesForNextLevel: 10,
		InitialSlowdown:         10,
		Frequency:               200 * time.Millisecond,
	}
	input := make(chan int)
	pit := doric.NewPit(doric.StandardHeight, doric.StandardWidth)
	source := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(source)

	// Start the game and return game events in the events channel
	events := doric.Play(pit, rnd, cfg, input)

	// Here you would need to start the game loop, manage input,
	// show graphics on screen, etc.

	// Listen for game events and act accordingly
	go func() {
		defer func() {
			close(input)
			// events channel will be closed when game is over
		}()
		for ev := range events {
			switch ev.(type) {
			case doric.EventScored:
				// Do whatever
			case doric.EventUpdated:
				// Do whatever
			case doric.EventRenewed:
				// Do whatever
			}
		}
	}()
}
