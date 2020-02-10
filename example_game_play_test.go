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
	command := make(chan int)
	pit := doric.NewPit(doric.StandardHeight, doric.StandardWidth)
	source := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(source)

	// Start the game and return game events in the events channel
	events := doric.Play(pit, rnd, cfg, command)

	defer func() {
		close(command)
		// events channel will be closed when game is over
	}()

	// Update game every 16 ms ~ 60 hz
	tick := time.Tick(16 * time.Millisecond)

	// Game loop
	for {
		// Listen for game events and act accordingly
		select {
		case ev, open := <-events:
			if !open {
				// game over, do whatever
				break
			}
			switch ev.(type) {
			case doric.EventScored:
				// Do whatever
			case doric.EventUpdated:
				// Do whatever
			case doric.EventRenewed:
				// Do whatever
			}
		case <-tick:
			// Update screen, send commands to game through the
			// command channel, etc.
		}
	}
}
