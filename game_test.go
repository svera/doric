package doric_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/svera/doric"
)

// mockTilesetBuilder implements the TilesFactory function to generate random tilesets
type mockTilesetBuilder struct {
	Tilesets [][3]int
	current  int
}

// build returns a tileset in the Tilesets property in the same order
// If all tilesets inside Tilesets were returned, the slice is ran again from the beginning
func (m *mockTilesetBuilder) build(n int) [3]int {
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
	factory := &mockTilesetBuilder{
		Tilesets: [][3]int{
			{1, 1, 1},
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
	factory := &mockTilesetBuilder{
		Tilesets: [][3]int{
			{1, 1, 1},
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
	setup := func() (chan<- int, <-chan interface{}) {
		well := doric.NewWell(doric.StandardHeight, doric.StandardWidth)
		factory := &mockTilesetBuilder{
			Tilesets: [][3]int{
				{1, 2, 3},
			},
		}
		command := make(chan int)
		events := doric.Play(well, factory.build, getConfig(), command)

		// First event received is just before game logic loop begins
		// the actual test will happen after that
		<-events
		command <- doric.CommandPauseSwitch

		return command, events
	}

	timeout := time.After(1 * time.Second)

	t.Run("Must not move left if paused", func(t *testing.T) {
		command, events := setup()

		command <- doric.CommandLeft

		select {
		case ev := <-events:
			if et, ok := ev.(doric.EventUpdated); ok {
				if et.Column.X == 3 {
					break
				}
				t.Errorf("Current column must not be moved left if game is paused")
			}
		case <-timeout:
			t.Errorf("Test timed out")
		}
	})

	t.Run("Must not move right if paused", func(t *testing.T) {
		command, events := setup()

		command <- doric.CommandRight

		select {
		case ev := <-events:
			if et, ok := ev.(doric.EventUpdated); ok {
				if et.Column.X == 3 {
					break
				}
				t.Errorf("Current column must not be moved right if game is paused")
			}
		case <-timeout:
			t.Errorf("Test timed out")
		}
	})

	t.Run("Must not move down if paused", func(t *testing.T) {
		command, events := setup()

		command <- doric.CommandDown

		select {
		case ev := <-events:
			if et, ok := ev.(doric.EventUpdated); ok && et.Column.Y == 0 {
				break
			}
			t.Errorf("Current column must not be moved down if game is paused")
		case <-timeout:
			t.Errorf("Test timed out")
		}
	})

	t.Run("Must not rotate if paused", func(t *testing.T) {
		command, events := setup()

		command <- doric.CommandRotate

		select {
		case ev := <-events:
			if et, ok := ev.(doric.EventUpdated); ok && et.Column.Tileset == [3]int{1, 2, 3} {
				break
			}
			t.Errorf("Current column must not be rotated if game is paused")
		case <-timeout:
			t.Errorf("Test timed out")
		}
	})
}

func TestWait(t *testing.T) {
	setup := func() (chan<- int, <-chan interface{}) {
		well := doric.NewWell(doric.StandardHeight, doric.StandardWidth)
		factory := &mockTilesetBuilder{
			Tilesets: [][3]int{
				{1, 2, 3},
			},
		}
		command := make(chan int)
		events := doric.Play(well, factory.build, getConfig(), command)

		// First event received is just before game logic loop begins
		// the actual test will happen after that
		<-events
		command <- doric.CommandWaitSwitch

		return command, events
	}

	timeout := time.After(1 * time.Second)

	t.Run("Must not move left if waiting", func(t *testing.T) {
		command, events := setup()

		command <- doric.CommandLeft

		select {
		case ev := <-events:
			if et, ok := ev.(doric.EventUpdated); ok && et.Column.X == 3 {
				break
			}
			t.Errorf("Current column must not be moved left if game is waiting")

		case <-timeout:
			t.Errorf("Test timed out")
		}
	})

	t.Run("Must not move right if waiting", func(t *testing.T) {
		command, events := setup()

		command <- doric.CommandRight

		select {
		case ev := <-events:
			if et, ok := ev.(doric.EventUpdated); ok && et.Column.X == 3 {
				break
			}
			t.Errorf("Current column must not be moved right if game is waiting")
		case <-timeout:
			t.Errorf("Test timed out")
		}
	})

	t.Run("Must not move down if waiting", func(t *testing.T) {
		command, events := setup()

		command <- doric.CommandDown

		select {
		case ev := <-events:
			if et, ok := ev.(doric.EventUpdated); ok && et.Column.Y == 0 {
				break
			}
			t.Errorf("Current column must not be moved down if game is waiting")
		case <-timeout:
			t.Errorf("Test timed out")
		}
	})

	t.Run("Must not rotate if waiting", func(t *testing.T) {
		command, events := setup()

		command <- doric.CommandRotate

		select {
		case ev := <-events:
			if et, ok := ev.(doric.EventUpdated); ok && et.Column.Tileset == [3]int{1, 2, 3} {
				break
			}
			t.Errorf("Current column must not be rotated if game is waiting")
		case <-timeout:
			t.Errorf("Test timed out")
		}
	})

	t.Run("Commands must work after unwaiting", func(t *testing.T) {
		command, events := setup()

		command <- doric.CommandWaitSwitch
		command <- doric.CommandLeft

		select {
		case ev := <-events:
			if et, ok := ev.(doric.EventUpdated); ok && et.Column.X == 2 {
				break
			}
			t.Errorf("Current column must be moved left if game is not waiting")
		case <-timeout:
			t.Errorf("Test timed out")
		}
	})
}

func TestCommands(t *testing.T) {
	setup := func() (chan<- int, <-chan interface{}) {
		well := doric.NewWell(doric.StandardHeight, doric.StandardWidth)
		factory := &mockTilesetBuilder{
			Tilesets: [][3]int{
				{1, 2, 3},
			},
		}
		command := make(chan int)
		events := doric.Play(well, factory.build, getConfig(), command)
		return command, events
	}

	timeout := time.After(1 * time.Second)

	t.Run("Must move left", func(t *testing.T) {
		command, events := setup()
		// First event received is just before game logic loop begins
		// the actual test will happen after that
		<-events

		command <- doric.CommandLeft

		select {
		case ev := <-events:
			if et, ok := ev.(doric.EventUpdated); ok && et.Column.X == 2 {
				break
			}
			t.Errorf("Current column must have been moved left")

		case <-timeout:
			t.Errorf("Test timed out")
		}
	})

	t.Run("Must move right", func(t *testing.T) {
		command, events := setup()

		<-events

		command <- doric.CommandRight

		select {
		case ev := <-events:
			if et, ok := ev.(doric.EventUpdated); ok && et.Column.X == 4 {
				break
			}
			t.Errorf("Current column must have been moved right")

		case <-timeout:
			t.Errorf("Test timed out")
		}
	})

	t.Run("Must move down", func(t *testing.T) {
		command, events := setup()

		<-events

		command <- doric.CommandDown

		select {
		case ev := <-events:
			if et, ok := ev.(doric.EventUpdated); ok && et.Column.Y == 1 {
				break
			}
			t.Errorf("Current column must have been moved down")
		case <-timeout:
			t.Errorf("Test timed out")
		}
	})

	t.Run("Must rotate", func(t *testing.T) {
		command, events := setup()

		<-events

		command <- doric.CommandRotate

		select {
		case ev := <-events:
			if et, ok := ev.(doric.EventUpdated); ok && et.Column.Tileset == [3]int{3, 1, 2} {
				break
			}
			t.Errorf("Current column must have been rotate")
		case <-timeout:
			t.Errorf("Test timed out")
		}
	})
}

func TestWellBounds(t *testing.T) {
	timeout := time.After(1 * time.Second)
	well := doric.NewWell(1, 1)
	factory := &mockTilesetBuilder{
		Tilesets: [][3]int{
			{1, 2, 3},
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
		if et, ok := ev.(doric.EventUpdated); ok && et.Column.X == 0 {
			break
		}
		t.Errorf("Current column must not move left as it clashes with well's left border")
	case <-timeout:
		t.Errorf("Test timed out")
	}

	command <- doric.CommandRight

	select {
	case ev := <-events:
		if et, ok := ev.(doric.EventUpdated); ok && et.Column.X == 0 {
			break
		}
		t.Errorf("Current column must not move right as it clashes with well's right border")
	case <-timeout:
		t.Errorf("Test timed out")
	}

	command <- doric.CommandDown

	select {
	case ev := <-events:
		if et, ok := ev.(doric.EventUpdated); ok && et.Column.Y == 0 {
			break
		}
		t.Errorf("Current column must not move down as it clashes with well's bottom")
	case <-timeout:
		t.Errorf("Test timed out")
	}
}

func TestScored(t *testing.T) {
	scoredTests := []struct {
		name                    string
		numberTilesForNextLevel int
		well                    doric.Well
		rand                    *mockTilesetBuilder
		expectedWell            doric.Well
		expectedRenewedWell     doric.Well
		expectedRemoved         int
		expectedLevel           int
		expectedCurrent         [3]int
	}{
		{
			name:                    "Scored with no level up",
			numberTilesForNextLevel: 20,
			rand: &mockTilesetBuilder{
				Tilesets: [][3]int{
					{1, 1, 1},
					{4, 5, 6},
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
			rand: &mockTilesetBuilder{
				Tilesets: [][3]int{
					{1, 1, 1},
					{4, 5, 6},
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
			rand: &mockTilesetBuilder{
				Tilesets: [][3]int{
					{1, 1, 1},
					{4, 5, 6},
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
						if asserted.Column.Tileset != tt.expectedCurrent {
							t.Errorf(
								"Expected that the next column was copied to the current one with values %v, got %v",
								tt.expectedCurrent,
								asserted.Column.Tileset,
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
		rand                    *mockTilesetBuilder
		expectedWells           []doric.Well
	}{
		{
			name:                    "Scored with combo",
			numberTilesForNextLevel: 20,
			rand: &mockTilesetBuilder{
				Tilesets: [][3]int{
					{1, 2, 3},
					{4, 5, 6},
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
