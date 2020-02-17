package doric_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/svera/doric"
)

// mockTilesFactory implements the TilesFactory function to generate random tilesets
type mockTilesFactory struct {
	Tilesets [][3]int
	current  int
}

// build returns a tileset in the Tilesets property in the same order
// If all tilesets inside Tilesets were returned, the slice is ran again from the beginning
func (m *mockTilesFactory) build(n int) [3]int {
	if m.current == len(m.Tilesets) {
		m.current = 0
	}
	val := m.Tilesets[m.current]
	m.current++
	return val
}

func getConfig() doric.Config {
	return doric.Config{
		NumberTilesForNextLevel: 10,
		InitialSpeed:            5,
		SpeedIncrement:          1,
		MaxSpeed:                13,
	}
}

func TestGameOver(t *testing.T) {
	timeout := time.After(1 * time.Second)
	well := doric.NewWell(1, doric.StandardWidth)
	factory := &mockTilesFactory{
		Tilesets: [][3]int{
			[3]int{1, 1, 1},
		},
	}
	command := make(chan int)
	well[3][0] = 1
	events := doric.Play(well, factory.build, getConfig(), command)

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
	well := doric.NewWell(doric.StandardHeight, doric.StandardWidth)
	factory := &mockTilesFactory{
		Tilesets: [][3]int{
			[3]int{1, 1, 1},
		},
	}
	command := make(chan int)
	events := doric.Play(well, factory.build, getConfig(), command)

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
	well := doric.NewWell(doric.StandardHeight, doric.StandardWidth)
	factory := &mockTilesFactory{
		Tilesets: [][3]int{
			[3]int{1, 2, 3},
		},
	}
	command := make(chan int)
	events := doric.Play(well, factory.build, getConfig(), command)

	// First event received is just before game logic loop begins
	// the actual test will happen after that
	<-events

	command <- doric.CommandPauseSwitch

	command <- doric.CommandLeft

	select {
	case ev := <-events:
		if et, ok := ev.(doric.EventUpdated); ok {
			if et.Current.X == 3 {
				break
			}
			t.Errorf("Current column must not be moved left if game is paused")
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
			t.Errorf("Current column must not be moved right if game is paused")
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
			t.Errorf("Current column must not be moved down if game is paused")
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
			t.Errorf("Current column must not be rotated if game is paused")
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}
}

func TestWait(t *testing.T) {
	timeout := time.After(1 * time.Second)
	well := doric.NewWell(doric.StandardHeight, doric.StandardWidth)
	factory := &mockTilesFactory{
		Tilesets: [][3]int{
			[3]int{1, 2, 3},
		},
	}
	command := make(chan int)
	cfg := getConfig()
	cfg.InitialSpeed = 0.5
	events := doric.Play(well, factory.build, cfg, command)

	// First event received is just before game logic loop begins
	// the actual test will happen after that
	<-events

	command <- doric.CommandWaitSwitch

	command <- doric.CommandLeft

	select {
	case ev := <-events:
		if et, ok := ev.(doric.EventUpdated); ok {
			if et.Current.X == 3 {
				break
			}
			t.Errorf("Current column must not be moved left if game is waiting")
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
			t.Errorf("Current column must not be moved right if game is waiting")
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
			t.Errorf("Current column must not be moved down if game is waiting")
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
			t.Errorf("Current column must not be rotated if game is waiting")
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}

	command <- doric.CommandWaitSwitch

	command <- doric.CommandLeft

	select {
	case ev := <-events:
		if et, ok := ev.(doric.EventUpdated); ok {
			if et.Current.X == 2 {
				break
			}
			t.Errorf("Current column must be moved left if game is not waiting")
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}
}

func TestCommands(t *testing.T) {
	timeout := time.After(1 * time.Second)
	well := doric.NewWell(doric.StandardHeight, doric.StandardWidth)
	factory := &mockTilesFactory{
		Tilesets: [][3]int{
			[3]int{1, 2, 3},
		},
	}
	command := make(chan int)
	events := doric.Play(well, factory.build, getConfig(), command)

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
			t.Errorf("Current column must be at column %d but is at %d", 2, et.Current.X)
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
			t.Errorf("Current column must be at column %d but is at %d", 3, et.Current.X)
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
			t.Errorf("Current column must be at row %d but is at %d", 1, et.Current.Y)
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
			t.Errorf("Current column must be as %v but is as %v", [3]int{3, 1, 2}, et.Current.Tiles)
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}

}

func TestWellBounds(t *testing.T) {
	timeout := time.After(1 * time.Second)
	well := doric.NewWell(1, 1)
	factory := &mockTilesFactory{
		Tilesets: [][3]int{
			[3]int{1, 2, 3},
		},
	}
	command := make(chan int)
	events := doric.Play(well, factory.build, getConfig(), command)

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
			t.Errorf("Current column must be at column %d but is at %d", 0, et.Current.X)
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
			t.Errorf("Current column must be at column %d but is at %d", 0, et.Current.X)
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
			t.Errorf("Current column must be at row %d but is at %d", 0, et.Current.Y)
		}
	case <-timeout:
		t.Errorf("Test timed out")
	}
}

func TestScored(t *testing.T) {
	scoredTests := []struct {
		name                    string
		numberTilesForNextLevel int
		well                    doric.Well
		rand                    *mockTilesFactory
		expectedWell            doric.Well
		expectedRenewedWell     doric.Well
		expectedRemoved         int
		expectedLevel           int
		expectedCurrent         [3]int
	}{
		{
			name:                    "Scored with no level up",
			numberTilesForNextLevel: 20,
			rand: &mockTilesFactory{
				Tilesets: [][3]int{
					[3]int{1, 1, 1},
					[3]int{4, 5, 6},
				},
			},
			well: transpose(doric.Well{
				[]int{0, 1, 0, 0, 0, 0},
				[]int{1, 1, 0, 0, 1, 1},
				[]int{1, 1, 1, 0, 1, 1},
			}),
			expectedWell: transpose(doric.Well{
				[]int{0, -1, 0, -1, 0, 0},
				[]int{1, -1, 0, -1, -1, -1},
				[]int{-1, -1, -1, -1, -1, -1},
			}),
			expectedRenewedWell: transpose(doric.Well{
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
			rand: &mockTilesFactory{
				Tilesets: [][3]int{
					[3]int{1, 1, 1},
					[3]int{4, 5, 6},
				},
			},
			well: transpose(doric.Well{
				[]int{0, 1, 0, 0, 0, 0},
				[]int{1, 1, 0, 0, 1, 1},
				[]int{1, 1, 1, 0, 1, 1},
			}),
			expectedWell: transpose(doric.Well{
				[]int{0, -1, 0, -1, 0, 0},
				[]int{1, -1, 0, -1, -1, -1},
				[]int{-1, -1, -1, -1, -1, -1},
			}),
			expectedRenewedWell: transpose(doric.Well{
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
			rand: &mockTilesFactory{
				Tilesets: [][3]int{
					[3]int{1, 1, 1},
					[3]int{4, 5, 6},
				},
			},
			well: transpose(doric.Well{
				[]int{1, 0, 0, 0, 0, 1},
				[]int{2, 1, 0, 0, 1, 2},
				[]int{3, 2, 1, 0, 2, 3},
			}),
			expectedWell: transpose(doric.Well{
				[]int{-1, 0, 0, -1, 0, -1},
				[]int{2, -1, 0, -1, -1, 2},
				[]int{3, 2, -1, -1, 2, 3},
			}),
			expectedRenewedWell: transpose(doric.Well{
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
			cfg.InitialSpeed = 20
			cfg.MaxSpeed = 40
			cfg.NumberTilesForNextLevel = tt.numberTilesForNextLevel
			command := make(chan int)
			events := doric.Play(tt.well, tt.rand.build, cfg, command)

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
						if !reflect.DeepEqual(tt.expectedWell, asserted.Well) {
							t.Errorf("Expected well %v but got %v", tt.expectedWell, asserted.Well)
						}
					case doric.EventRenewed:
						if asserted.Current.Tiles != tt.expectedCurrent {
							t.Errorf(
								"Expected that the next column was copied to the current one with values %v, got %v",
								tt.expectedCurrent,
								asserted.Current.Tiles,
							)
						}
						if !reflect.DeepEqual(tt.expectedRenewedWell, asserted.Well) {
							t.Errorf("Expected well %v but got %v", tt.expectedRenewedWell, asserted.Well)
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
		well                    doric.Well
		rand                    *mockTilesFactory
		expectedWells           []doric.Well
	}{
		{
			name:                    "Scored with combo",
			numberTilesForNextLevel: 20,
			rand: &mockTilesFactory{
				Tilesets: [][3]int{
					[3]int{1, 2, 3},
					[3]int{4, 5, 6},
				},
			},
			well: transpose(doric.Well{
				[]int{0, 0, 0, 0, 0, 0},
				[]int{0, 0, 0, 0, 0, 0},
				[]int{0, 2, 2, 0, 1, 1},
			}),
			expectedWells: []doric.Well{
				transpose(doric.Well{
					[]int{0, 0, 0, 3, 0, 0},
					[]int{0, 0, 0, 2, 0, 0},
					[]int{0, 2, 2, -1, -1, -1},
				}),
				transpose(doric.Well{
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
			cfg.InitialSpeed = 20
			cfg.NumberTilesForNextLevel = tt.numberTilesForNextLevel
			command := make(chan int)
			events := doric.Play(tt.well, tt.rand.build, cfg, command)

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
						if !reflect.DeepEqual(tt.expectedWells[count], asserted.Well) {
							t.Errorf("Expected well %v but got %v", tt.expectedWells[count], asserted.Well)
						}
						if count == len(tt.expectedWells)-1 {
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

// transpose swaps values between columns and rows in a well,
// as visually is more natural to put all X values per row in a single line
// (as they appear on actual games)
// and this can only be done by putting those in the Y index of the well type.
func transpose(slice doric.Well) doric.Well {
	xl := len(slice[0])
	yl := len(slice)
	result := doric.NewWell(yl, xl)
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = slice[j][i]
		}
	}
	return result
}
