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

func getConfig() columns.Config {
	return columns.Config{
		NumberTilesForNextLevel: 10,
		InitialSlowdown:         1,
		Frequency:               1 * time.Millisecond,
	}
}

func TestGameOver(t *testing.T) {
	timeout := time.After(1 * time.Second)
	pit := columns.NewPit(1, pithWidth)
	r := &mocks.Randomizer{Values: []int{0}}
	input := make(chan int)
	pit[0][3] = 1
	events := columns.Play(pit, r, getConfig(), input)

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
	input := make(chan int)
	events := columns.Play(pit, r, getConfig(), input)

	// First event received is just before game logic loop begins
	// the actual test will happen after that
	<-events

	input <- columns.ActionPause

	select {
	case ev := <-events:
		if et, ok := ev.(columns.EventUpdated); ok {
			if !et.Paused {
				t.Errorf("Game must be paused")
			}
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}

	input <- columns.ActionPause

	select {
	case ev := <-events:
		if et, ok := ev.(columns.EventUpdated); ok {
			if et.Paused {
				t.Errorf("Game must not be paused")
			}
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}
}

func TestInput(t *testing.T) {
	timeout := time.After(1 * time.Second)
	pit := columns.NewPit(pitHeight, pithWidth)
	r := &mocks.Randomizer{Values: []int{0, 1, 2}}
	input := make(chan int)
	events := columns.Play(pit, r, getConfig(), input)

	// First event received is just before game logic loop begins
	// the actual test will happen after that
	<-events

	input <- columns.ActionLeft

	select {
	case ev := <-events:
		if et, ok := ev.(columns.EventUpdated); ok {
			if et.Current.X == 2 {
				break
			}
			t.Errorf("Current piece must be at column %d but is at %d", 2, et.Current.X)
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}

	input <- columns.ActionRight

	select {
	case ev := <-events:
		if et, ok := ev.(columns.EventUpdated); ok {
			if et.Current.X == 3 {
				break
			}
			t.Errorf("Current piece must be at column %d but is at %d", 3, et.Current.X)
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}

	input <- columns.ActionDown

	select {
	case ev := <-events:
		if et, ok := ev.(columns.EventUpdated); ok {
			if et.Current.Y == 1 {
				break
			}
			t.Errorf("Current piece must be at row %d but is at %d", 1, et.Current.Y)
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}

	input <- columns.ActionRotate

	select {
	case ev := <-events:
		switch et := ev.(type) {
		case columns.EventUpdated:
			if et.Current.Tiles == [3]int{3, 1, 2} {
				break
			}
			t.Errorf("Current piece must be as %v but is as %v", [3]int{3, 1, 2}, et.Current.Tiles)
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}

}

func TestConsolidated(t *testing.T) {
	timeout := time.After(1 * time.Second)
	pit := columns.NewPit(3, pithWidth)
	initialPit := columns.NewPit(3, pithWidth)
	r := &mocks.Randomizer{Values: []int{0, 1, 2}}
	input := make(chan int)
	var previous columns.Piece
	events := columns.Play(pit, r, getConfig(), input)

	for {
		select {
		case ev := <-events:
			switch et := ev.(type) {
			case columns.EventUpdated:
				previous = et.Current
			case columns.EventScored:
				if reflect.DeepEqual(initialPit, et.Pit) {
					t.Errorf("Previous piece wasn't consolidated in pit")
				}
			case columns.EventRenewed:
				if reflect.DeepEqual(previous, et.Current) {
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
		expectedRemoved         int
		expectedLevel           int
		expectedCombo           int
		expectedTiles           [3]int
	}{
		{
			name:                    "Scored with no level up",
			numberTilesForNextLevel: 10,
			expectedRemoved:         3,
			expectedLevel:           1,
			expectedCombo:           1,
			expectedTiles:           [3]int{4, 5, 6},
		},
		{
			name:                    "Scored with level up",
			numberTilesForNextLevel: 1,
			expectedRemoved:         3,
			expectedLevel:           2,
			expectedCombo:           1,
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
			input := make(chan int)
			events := columns.Play(pit, r, cfg, input)

			<-events

			for {
				select {
				case ev := <-events:
					switch et := ev.(type) {
					case columns.EventScored:
						if et.Removed != tt.expectedRemoved {
							t.Errorf("Expected %d removed tiles but got %d", tt.expectedRemoved, et.Removed)
						}
						if et.Level != tt.expectedLevel {
							t.Errorf("Expected level %d but got %d", tt.expectedLevel, et.Level)
						}
						if et.Combo != tt.expectedCombo {
							t.Errorf("Expected combo value %d but got %d", tt.expectedCombo, et.Combo)
						}
						return
					case columns.EventRenewed:
						if et.Current.Tiles != tt.expectedTiles {
							t.Errorf(
								"Expected that the next piece was copied to the current one with values %v, got %v",
								tt.expectedTiles,
								et.Current.Tiles,
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
