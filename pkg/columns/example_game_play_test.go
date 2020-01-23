package columns_test

import (
	"math/rand"
	"time"

	"github.com/svera/doric/pkg/columns"
)

func Example() {
	cfg := columns.Config{
		NumberTilesForNextLevel: 10,
		InitialSlowdown:         10,
		Frequency:               200 * time.Millisecond,
	}
	input := make(chan int)
	pit := columns.NewPit(13, 6)
	source := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(source)

	game, events := columns.NewGame(pit, rnd, cfg)
	// Start the game and return game events in the events channel
	go game.Play(input)

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
			case columns.EventScored:
				// Do whatever
			case columns.EventUpdated:
				// Do whatever
			case columns.EventRenewed:
				// Do whatever
			}
		}
	}()
}
