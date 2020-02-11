package doric_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/svera/doric"
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
	pit := doric.NewPit(1, doric.StandardWidth)
	r := &doric.MockRandomizer{Values: []int{0}}
	command := make(chan int)
	pit[3][0] = 1
	events := doric.Play(pit, r, getConfig(), command)

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

func TestQuit(t *testing.T) {
	timeout := time.After(1 * time.Second)
	pit := doric.NewPit(doric.StandardHeight, doric.StandardWidth)
	r := &doric.MockRandomizer{Values: []int{0}}
	command := make(chan int)
	events := doric.Play(pit, r, getConfig(), command)

	// First event received is just before game logic loop begins
	// the actual test will happen after that
	<-events

	command <- doric.CommandQuit

	for {
		select {
		case _, open := <-events:
			if !open {
				return
			}
		case <-timeout:
			t.Errorf("Game should be quitted")
		}
	}
}

func TestPause(t *testing.T) {
	timeout := time.After(1 * time.Second)
	pit := doric.NewPit(doric.StandardHeight, doric.StandardWidth)
	r := &doric.MockRandomizer{Values: []int{0, 1, 2}}
	command := make(chan int)
	events := doric.Play(pit, r, getConfig(), command)

	// First event received is just before game logic loop begins
	// the actual test will happen after that
	<-events

	command <- doric.CommandPause

	command <- doric.CommandLeft

	select {
	case ev := <-events:
		if et, ok := ev.(doric.EventUpdated); ok {
			if et.Current.X == 3 {
				break
			}
			t.Errorf("Current piece must not be moved left if game is paused")
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}

	command <- doric.CommandRight

	select {
	case ev := <-events:
		if et, ok := ev.(doric.EventUpdated); ok {
			if et.Current.X == 3 {
				break
			}
			t.Errorf("Current piece must not be moved right if game is paused")
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}

	command <- doric.CommandDown

	select {
	case ev := <-events:
		if et, ok := ev.(doric.EventUpdated); ok {
			if et.Current.Y == 0 {
				break
			}
			t.Errorf("Current piece must not be moved down if game is paused")
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}

	command <- doric.CommandRotate

	select {
	case ev := <-events:
		if et, ok := ev.(doric.EventUpdated); ok {
			if et.Current.Tiles == [3]int{1, 2, 3} {
				break
			}
			t.Errorf("Current piece must not be rotated if game is paused")
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}
}

func TestWait(t *testing.T) {
	timeout := time.After(1 * time.Second)
	pit := doric.NewPit(doric.StandardHeight, doric.StandardWidth)
	r := &doric.MockRandomizer{Values: []int{0, 1, 2}}
	command := make(chan int)
	cfg := getConfig()
	cfg.Frequency = 10 * time.Second
	events := doric.Play(pit, r, cfg, command)

	// First event received is just before game logic loop begins
	// the actual test will happen after that
	<-events

	command <- doric.CommandWait

	command <- doric.CommandLeft

	select {
	case ev := <-events:
		if et, ok := ev.(doric.EventUpdated); ok {
			if et.Current.X == 3 {
				break
			}
			t.Errorf("Current piece must not be moved left if game is waiting")
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}

	command <- doric.CommandRight

	select {
	case ev := <-events:
		if et, ok := ev.(doric.EventUpdated); ok {
			if et.Current.X == 3 {
				break
			}
			t.Errorf("Current piece must not be moved right if game is waiting")
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}

	command <- doric.CommandDown

	select {
	case ev := <-events:
		if et, ok := ev.(doric.EventUpdated); ok {
			if et.Current.Y == 0 {
				break
			}
			t.Errorf("Current piece must not be moved down if game is waiting")
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}

	command <- doric.CommandRotate

	select {
	case ev := <-events:
		if et, ok := ev.(doric.EventUpdated); ok {
			if et.Current.Tiles == [3]int{1, 2, 3} {
				break
			}
			t.Errorf("Current piece must not be rotated if game is waiting")
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}

	command <- doric.CommandWait

	command <- doric.CommandLeft

	select {
	case ev := <-events:
		if et, ok := ev.(doric.EventUpdated); ok {
			if et.Current.X == 2 {
				break
			}
			t.Errorf("Current piece must be moved left if game is not waiting")
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}
}

func TestInput(t *testing.T) {
	timeout := time.After(1 * time.Second)
	pit := doric.NewPit(doric.StandardHeight, doric.StandardWidth)
	r := &doric.MockRandomizer{Values: []int{0, 1, 2}}
	command := make(chan int)
	events := doric.Play(pit, r, getConfig(), command)

	// First event received is just before game logic loop begins
	// the actual test will happen after that
	<-events

	command <- doric.CommandLeft

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

	command <- doric.CommandRight

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

	command <- doric.CommandDown

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

	command <- doric.CommandRotate

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
	command := make(chan int)
	events := doric.Play(pit, r, getConfig(), command)

	// First event received is just before game logic loop begins
	// the actual test will happen after that
	<-events

	command <- doric.CommandLeft

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

	command <- doric.CommandRight

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

	command <- doric.CommandDown

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
		rand                    *doric.MockRandomizer
		expectedPit             doric.Pit
		expectedRenewedPit      doric.Pit
		expectedRemoved         int
		expectedLevel           int
		expectedCurrent         [3]int
	}{
		{
			name:                    "Scored with no level up",
			numberTilesForNextLevel: 20,
			rand:                    &doric.MockRandomizer{Values: []int{0, 0, 0, 3, 4, 5}},
			pit: transpose(doric.Pit{
				[]int{0, 1, 0, 0, 0, 0},
				[]int{1, 1, 0, 0, 1, 1},
				[]int{1, 1, 1, 0, 1, 1},
			}),
			expectedPit: transpose(doric.Pit{
				[]int{0, -1, 0, -1, 0, 0},
				[]int{1, -1, 0, -1, -1, -1},
				[]int{-1, -1, -1, -1, -1, -1},
			}),
			expectedRenewedPit: transpose(doric.Pit{
				[]int{0, 0, 0, 0, 0, 0},
				[]int{0, 0, 0, 0, 0, 0},
				[]int{1, 0, 0, 0, 0, 0},
			}),
			expectedRemoved: 12,
			expectedLevel:   1,
			expectedCurrent: [3]int{4, 5, 6},
		},
		{
			name:                    "Scored with level up",
			numberTilesForNextLevel: 1,
			rand:                    &doric.MockRandomizer{Values: []int{0, 0, 0, 3, 4, 5}},
			pit: transpose(doric.Pit{
				[]int{0, 1, 0, 0, 0, 0},
				[]int{1, 1, 0, 0, 1, 1},
				[]int{1, 1, 1, 0, 1, 1},
			}),
			expectedPit: transpose(doric.Pit{
				[]int{0, -1, 0, -1, 0, 0},
				[]int{1, -1, 0, -1, -1, -1},
				[]int{-1, -1, -1, -1, -1, -1},
			}),
			expectedRenewedPit: transpose(doric.Pit{
				[]int{0, 0, 0, 0, 0, 0},
				[]int{0, 0, 0, 0, 0, 0},
				[]int{1, 0, 0, 0, 0, 0},
			}),
			expectedRemoved: 12,
			expectedLevel:   2,
			expectedCurrent: [3]int{4, 5, 6},
		},
		{
			name:                    "Diagonal lines",
			numberTilesForNextLevel: 20,
			rand:                    &doric.MockRandomizer{Values: []int{0, 0, 0, 3, 4, 5}},
			pit: transpose(doric.Pit{
				[]int{1, 0, 0, 0, 0, 1},
				[]int{2, 1, 0, 0, 1, 2},
				[]int{3, 2, 1, 0, 2, 3},
			}),
			expectedPit: transpose(doric.Pit{
				[]int{-1, 0, 0, -1, 0, -1},
				[]int{2, -1, 0, -1, -1, 2},
				[]int{3, 2, -1, -1, 2, 3},
			}),
			expectedRenewedPit: transpose(doric.Pit{
				[]int{0, 0, 0, 0, 0, 0},
				[]int{2, 0, 0, 0, 0, 2},
				[]int{3, 2, 0, 0, 2, 3},
			}),
			expectedRemoved: 8,
			expectedLevel:   1,
			expectedCurrent: [3]int{4, 5, 6},
		},
	}

	for _, tt := range scoredTests {
		t.Run(tt.name, func(t *testing.T) {
			timeout := time.After(1 * time.Second)
			cfg := getConfig()
			cfg.InitialSlowdown = 2
			cfg.NumberTilesForNextLevel = tt.numberTilesForNextLevel
			command := make(chan int)
			events := doric.Play(tt.pit, tt.rand, cfg, command)

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
						if !reflect.DeepEqual(tt.expectedRenewedPit, asserted.Pit) {
							t.Errorf("Expected pit %v but got %v", tt.expectedRenewedPit, asserted.Pit)
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

func TestScoredCombo(t *testing.T) {
	comboTests := []struct {
		name                    string
		numberTilesForNextLevel int
		pit                     doric.Pit
		rand                    *doric.MockRandomizer
		expectedPits            []doric.Pit
	}{
		{
			name:                    "Scored with combo",
			numberTilesForNextLevel: 20,
			rand:                    &doric.MockRandomizer{Values: []int{0, 1, 2, 3, 4, 5}},
			pit: transpose(doric.Pit{
				[]int{0, 0, 0, 0, 0, 0},
				[]int{0, 0, 0, 0, 0, 0},
				[]int{0, 2, 2, 0, 1, 1},
			}),
			expectedPits: []doric.Pit{
				transpose(doric.Pit{
					[]int{0, 0, 0, 3, 0, 0},
					[]int{0, 0, 0, 2, 0, 0},
					[]int{0, 2, 2, -1, -1, -1},
				}),
				transpose(doric.Pit{
					[]int{0, 0, 0, 0, 0, 0},
					[]int{0, 0, 0, 3, 0, 0},
					[]int{0, -1, -1, -1, 0, 0},
				}),
			},
		},
	}

	for _, tt := range comboTests {
		t.Run(tt.name, func(t *testing.T) {
			timeout := time.After(1 * time.Second)
			cfg := getConfig()
			cfg.InitialSlowdown = 2
			cfg.NumberTilesForNextLevel = tt.numberTilesForNextLevel
			command := make(chan int)
			events := doric.Play(tt.pit, tt.rand, cfg, command)

			<-events

			count := 0
			for {
				select {
				case ev := <-events:
					switch asserted := ev.(type) {
					case doric.EventScored:
						if asserted.Combo != count+1 {
							t.Errorf("Expected combo value %d but got %d", count, asserted.Combo)
						}
						if !reflect.DeepEqual(tt.expectedPits[count], asserted.Pit) {
							t.Errorf("Expected pit %v but got %v", tt.expectedPits[count], asserted.Pit)
						}
						if count == len(tt.expectedPits)-1 {
							return
						}
						count++
					}

				case <-timeout:
					t.Fatalf("Test timed out and no scored update reached")
				}
			}
		})
	}
}

// transpose swaps values between columns and rows in a pit,
// as visually is more natural to put all X values per row in a single line
// (as they appear on actual games)
// and this can only be done by putting those in the Y index of the pit type.
func transpose(slice doric.Pit) doric.Pit {
	xl := len(slice[0])
	yl := len(slice)
	result := doric.NewPit(yl, xl)
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = slice[j][i]
		}
	}
	return result
}
