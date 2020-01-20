package columns_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/svera/doric/pkg/columns"
	"github.com/svera/doric/pkg/columns/mocks"
)

const (
	pithWidth = 6
	pitHeight = 13
)

var codeToEventName = [3]string{
	"EventUpdated",
	"EventScored",
	"EventRenewed",
}

func getConfig() columns.Config {
	return columns.Config{
		PointsPerTile:           10,
		NumberTilesForNextLevel: 10,
		InitialSlowdown:         1,
		Frequency:               1 * time.Millisecond,
	}
}

func TestGameOver(t *testing.T) {
	timeout := time.After(1 * time.Second)
	pit := columns.NewPit(1, pithWidth)
	r := &mocks.Randomizer{Values: []int{0}}
	game := columns.NewGame(pit, r, getConfig())
	events := make(chan columns.Event)
	input := make(chan int)
	pit[0][3] = 1
	go game.Play(input, events)
	for {
		select {
		case _, open := <-events:
			if !open {
				return
			}
		case <-timeout:
			t.Errorf("Game should be over")
		}
	}
}

func TestPause(t *testing.T) {
	timeout := time.After(1 * time.Second)
	pit := columns.NewPit(pitHeight, pithWidth)
	r := &mocks.Randomizer{Values: []int{1}}
	game := columns.NewGame(pit, r, getConfig())
	events := make(chan columns.Event)
	input := make(chan int)
	go game.Play(input, events)

	go func() {
		input <- columns.ActionPause
	}()

	// First event received is just before game logic loop begins
	// the actual test will happen after that
	<-events

	select {
	case ev := <-events:
		if !ev.Status.Paused {
			t.Errorf("Game must be paused, got '%s'", codeToEventName[ev.ID])
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}

	go func() {
		input <- columns.ActionPause
	}()

	select {
	case ev := <-events:
		if ev.Status.Paused {
			t.Errorf("Game must not be paused, got '%s'", codeToEventName[ev.ID])
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}

}

func TestInput(t *testing.T) {
	timeout := time.After(1 * time.Second)
	pit := columns.NewPit(pitHeight, pithWidth)
	r := &mocks.Randomizer{Values: []int{0, 1, 2}}
	game := columns.NewGame(pit, r, getConfig())
	events := make(chan columns.Event)
	input := make(chan int)
	go game.Play(input, events)

	go func() {
		input <- columns.ActionLeft
	}()

	// First event received is just before game logic loop begins
	// the actual test will happen after that
	<-events

	select {
	case ev := <-events:
		if ev.Status.Current.X() == 2 {
			break
		}
		t.Errorf("Current piece must be at column %d but is at %d", 2, ev.Status.Current.X())
	case <-timeout:
		t.Errorf("Test timed out")
	}

	go func() {
		input <- columns.ActionRight
	}()

	select {
	case ev := <-events:
		if ev.Status.Current.X() == 3 {
			break
		}
		t.Errorf("Current piece must be at column %d but is at %d", 3, ev.Status.Current.X())
	case <-timeout:
		t.Errorf("Test timed out")
	}

	go func() {
		input <- columns.ActionDown
	}()

	select {
	case ev := <-events:
		if ev.Status.Current.Y() == 1 {
			break
		}
		t.Errorf("Current piece must be at row %d but is at %d", 1, ev.Status.Current.Y())
	case <-timeout:
		t.Errorf("Test timed out")
	}

	go func() {
		input <- columns.ActionRotate
	}()

	select {
	case ev := <-events:
		if ev.Status.Current.Tiles() == [3]int{3, 1, 2} {
			break
		}
		t.Errorf("Current piece must be as %v but is as %v", [3]int{3, 1, 2}, ev.Status.Current.Tiles())
	case <-timeout:
		t.Errorf("Test timed out")
	}

}

func TestConsolidated(t *testing.T) {
	timeout := time.After(1 * time.Second)
	pit := columns.NewPit(3, pithWidth)
	initialPit := columns.NewPit(3, pithWidth)
	r := &mocks.Randomizer{Values: []int{0, 1, 2}}
	cfg := getConfig()
	game := columns.NewGame(pit, r, cfg)
	events := make(chan columns.Event)
	input := make(chan int)
	var previous columns.Piece
	go game.Play(input, events)

	for {
		select {
		case ev := <-events:
			if ev.ID == columns.EventUpdated {
				previous = ev.Status.Current
			}
			if ev.ID == columns.EventRenewed {
				if reflect.DeepEqual(initialPit, ev.Status.Pit) {
					t.Errorf("Previous piece wasn't consolidated in pit")
				}
				if reflect.DeepEqual(previous, ev.Status.Current) {
					t.Errorf("Current piece wasn't renewed")
				}
				return
			}
		case <-timeout:
			t.Errorf("Test timed out and current piece wasn't renewed")
		}
	}

}

func TestScored(t *testing.T) {
	scoredTests := []struct {
		name                    string
		numberTilesForNextLevel int
		expectedScore           int
		expectedLevel           int
		expectedTiles           [3]int
	}{
		{
			name:                    "Scored with no level up",
			numberTilesForNextLevel: 10,
			expectedScore:           30,
			expectedLevel:           1,
			expectedTiles:           [3]int{4, 5, 6},
		},
		{
			name:                    "Scored with level up",
			numberTilesForNextLevel: 1,
			expectedScore:           30,
			expectedLevel:           2,
			expectedTiles:           [3]int{4, 5, 6},
		},
	}

	for _, tt := range scoredTests {
		t.Run(tt.name, func(t *testing.T) {
			timeout := time.After(1 * time.Second)
			pit := columns.NewPit(3, pithWidth)
			r := &mocks.Randomizer{Values: []int{0, 0, 0, 3, 4, 5}}
			cfg := getConfig()
			cfg.InitialSlowdown = 2
			cfg.NumberTilesForNextLevel = tt.numberTilesForNextLevel
			game := columns.NewGame(pit, r, cfg)
			events := make(chan columns.Event)
			input := make(chan int)
			go game.Play(input, events)

			for {
				select {
				case ev := <-events:
					if ev.ID == columns.EventScored {
						if ev.Status.Points != tt.expectedScore {
							t.Errorf("Expected %d points but got %d", tt.expectedScore, ev.Status.Points)
						}
						if ev.Status.Level != tt.expectedLevel {
							t.Errorf("Expected level %d but got %d", tt.expectedLevel, ev.Status.Level)
						}
						return
					}

					if ev.ID == columns.EventRenewed {
						if ev.Status.Current.Tiles() != tt.expectedTiles {
							t.Errorf(
								"Expected that the next piece was copied to the current one with values %v, got %v",
								tt.expectedTiles,
								ev.Status.Current.Tiles(),
							)
						}
						return
					}

				case <-timeout:
					t.Fatalf("Test timed out and no scored update reached")
				}
			}
		})
	}

}
