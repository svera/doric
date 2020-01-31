package doric_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/svera/doric"
)

const (
	pithWidth = 6
	pitHeight = 13
)

func getConfig() doric.Config {
	return doric.Config{
		NumberTilesForNextLevel: 10,
		InitialSlowdown:         1,
		Frequency:               1 * time.Millisecond,
	}
}

func TestGameOver(t *testing.T) {
	timeout := time.After(1 * time.Second)
	pit := doric.NewPit(1, pithWidth)
	r := &doric.MockRandomizer{Values: []int{0}}
	input := make(chan int)
	pit[0][3] = 1
	events := doric.Play(pit, r, getConfig(), input)

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
	pit := doric.NewPit(pitHeight, pithWidth)
	r := &doric.MockRandomizer{Values: []int{1}}
	input := make(chan int)
	events := doric.Play(pit, r, getConfig(), input)

	// First event received is just before game logic loop begins
	// the actual test will happen after that
	<-events

	input <- doric.ActionPause

	select {
	case ev := <-events:
		if et, ok := ev.(doric.EventUpdated); ok {
			if !et.Paused {
				t.Errorf("Game must be paused")
			}
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}

	input <- doric.ActionPause

	select {
	case ev := <-events:
		if et, ok := ev.(doric.EventUpdated); ok {
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
	pit := doric.NewPit(pitHeight, pithWidth)
	r := &doric.MockRandomizer{Values: []int{0, 1, 2}}
	input := make(chan int)
	events := doric.Play(pit, r, getConfig(), input)

	// First event received is just before game logic loop begins
	// the actual test will happen after that
	<-events

	input <- doric.ActionLeft

	select {
	case ev := <-events:
		if et, ok := ev.(doric.EventUpdated); ok {
			if et.Current.X == 2 {
				break
			}
			t.Errorf("Current piece must be at column %d but is at %d", 2, et.Current.X)
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}

	input <- doric.ActionRight

	select {
	case ev := <-events:
		if et, ok := ev.(doric.EventUpdated); ok {
			if et.Current.X == 3 {
				break
			}
			t.Errorf("Current piece must be at column %d but is at %d", 3, et.Current.X)
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}

	input <- doric.ActionDown

	select {
	case ev := <-events:
		if et, ok := ev.(doric.EventUpdated); ok {
			if et.Current.Y == 1 {
				break
			}
			t.Errorf("Current piece must be at row %d but is at %d", 1, et.Current.Y)
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}

	input <- doric.ActionRotate

	select {
	case ev := <-events:
		switch et := ev.(type) {
		case doric.EventUpdated:
			if et.Current.Tiles == [3]int{3, 1, 2} {
				break
			}
			t.Errorf("Current piece must be as %v but is as %v", [3]int{3, 1, 2}, et.Current.Tiles)
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}

}

func TestPitBounds(t *testing.T) {
	timeout := time.After(1 * time.Second)
	pit := doric.NewPit(1, 1)
	r := &doric.MockRandomizer{Values: []int{0, 1, 2}}
	input := make(chan int)
	events := doric.Play(pit, r, getConfig(), input)

	// First event received is just before game logic loop begins
	// the actual test will happen after that
	<-events

	input <- doric.ActionLeft

	select {
	case ev := <-events:
		if et, ok := ev.(doric.EventUpdated); ok {
			if et.Current.X == 0 {
				break
			}
			t.Errorf("Current piece must be at column %d but is at %d", 0, et.Current.X)
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}

	input <- doric.ActionRight

	select {
	case ev := <-events:
		if et, ok := ev.(doric.EventUpdated); ok {
			if et.Current.X == 0 {
				break
			}
			t.Errorf("Current piece must be at column %d but is at %d", 0, et.Current.X)
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}

	input <- doric.ActionDown

	select {
	case ev := <-events:
		if et, ok := ev.(doric.EventUpdated); ok {
			if et.Current.Y == 0 {
				break
			}
			t.Errorf("Current piece must be at row %d but is at %d", 0, et.Current.Y)
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}
}

func TestScored(t *testing.T) {
	scoredTests := []struct {
		name                    string
		numberTilesForNextLevel int
		pit                     doric.Pit
		expectedPit             doric.Pit
		expectedRemoved         int
		expectedLevel           int
		expectedCombo           int
		expectedCurrent         [3]int
	}{
		{
			name:                    "Scored with no level up",
			numberTilesForNextLevel: 20,
			pit: doric.Pit{
				[]int{0, 1, 0, 0, 0, 0},
				[]int{1, 1, 0, 0, 1, 1},
				[]int{1, 1, 1, 0, 1, 1},
			},
			expectedPit: doric.Pit{
				[]int{0, -1, 0, -1, 0, 0},
				[]int{1, -1, 0, -1, -1, -1},
				[]int{-1, -1, -1, -1, -1, -1},
			},
			expectedRemoved: 12,
			expectedLevel:   1,
			expectedCombo:   1,
			expectedCurrent: [3]int{4, 5, 6},
		},
		{
			name:                    "Scored with level up",
			numberTilesForNextLevel: 1,
			pit: doric.Pit{
				[]int{0, 1, 0, 0, 0, 0},
				[]int{1, 1, 0, 0, 1, 1},
				[]int{1, 1, 1, 0, 1, 1},
			},
			expectedPit: doric.Pit{
				[]int{0, -1, 0, -1, 0, 0},
				[]int{1, -1, 0, -1, -1, -1},
				[]int{-1, -1, -1, -1, -1, -1},
			},
			expectedRemoved: 12,
			expectedLevel:   2,
			expectedCombo:   1,
			expectedCurrent: [3]int{4, 5, 6},
		},
		{
			name:                    "Diagonal lines",
			numberTilesForNextLevel: 20,
			pit: doric.Pit{
				{1, 0, 0, 0, 0, 1},
				{2, 1, 0, 0, 1, 2},
				{3, 2, 1, 0, 2, 3},
			},
			expectedPit: doric.Pit{
				{-1, 0, 0, -1, 0, -1},
				{2, -1, 0, -1, -1, 2},
				{3, 2, -1, -1, 2, 3},
			},
			expectedRemoved: 8,
			expectedLevel:   1,
			expectedCombo:   1,
			expectedCurrent: [3]int{4, 5, 6},
		},
	}

	for _, tt := range scoredTests {
		t.Run(tt.name, func(t *testing.T) {
			timeout := time.After(1 * time.Second)
			r := &doric.MockRandomizer{Values: []int{0, 0, 0, 3, 4, 5}}
			cfg := getConfig()
			cfg.InitialSlowdown = 2
			cfg.NumberTilesForNextLevel = tt.numberTilesForNextLevel
			input := make(chan int)
			events := doric.Play(tt.pit, r, cfg, input)

			<-events

			for {
				select {
				case ev := <-events:
					switch asserted := ev.(type) {
					case doric.EventScored:
						if asserted.Removed != tt.expectedRemoved {
							t.Errorf("Expected %d removed tiles but got %d", tt.expectedRemoved, asserted.Removed)
						}
						if asserted.Level != tt.expectedLevel {
							t.Errorf("Expected level %d but got %d", tt.expectedLevel, asserted.Level)
						}
						if asserted.Combo != tt.expectedCombo {
							t.Errorf("Expected combo value %d but got %d", tt.expectedCombo, asserted.Combo)
						}
						if !reflect.DeepEqual(tt.expectedPit, asserted.Pit) {
							t.Errorf("Expected pit %v but got %v", tt.expectedPit, asserted.Pit)
						}
					case doric.EventRenewed:
						if asserted.Current.Tiles != tt.expectedCurrent {
							t.Errorf(
								"Expected that the next piece was copied to the current one with values %v, got %v",
								tt.expectedCurrent,
								asserted.Current.Tiles,
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
