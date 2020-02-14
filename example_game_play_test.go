package doric_test

import (
	"math/rand"
	"time"

	"github.com/svera/doric"
)

func Example() {
	cfg := doric.Config{
		NumberTilesForNextLevel: 10,
		InitialSpeed:            0.5,
		SpeedIncrement:          0.25,
		MaxSpeed:                13,
	}
	command := make(chan int)
	well := doric.NewWell(doric.StandardHeight, doric.StandardWidth)
	rnd := func(n int) [3]int {
		source := rand.NewSource(time.Now().UnixNano())
		rand.New(source)
		return [3]int{
			rand.Intn(n) + 1,
			rand.Intn(n) + 1,
			rand.Intn(n) + 1,
		}
	}

	// Start the game and return game events in the events channel
	events := doric.Play(well, rnd, cfg, command)

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
