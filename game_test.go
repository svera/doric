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

func defaultConfig() doric.Config {
	return doric.Config{
		NumberTilesForNextLevel: 10,
		InitialSpeed:            5,
		SpeedIncrement:          1,
		MaxSpeed:                13,
	}
}

func setup(cfg doric.Config, well doric.Well, ts [][3]int) (chan<- int, <-chan interface{}, <-chan time.Time) {
	timeout := time.After(1 * time.Second)
	factory := &mockTilesetBuilder{
		Tilesets: ts,
	}
	commands := make(chan int)
	events := doric.Play(well, factory.build, cfg, commands)

	// First event received is just before game logic loop begins
	// the actual test will happen after that
	<-events

	return commands, events, timeout
}

func TestGameOver(t *testing.T) {
	well := doric.NewWell(doric.StandardWidth, 1)
	well[3][0] = 1
	_, events, timeout := setup(
		defaultConfig(),
		well,
		[][3]int{{1, 1, 1}},
	)

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
	commands, events, timeout := setup(
		defaultConfig(),
		doric.NewWell(doric.StandardWidth, doric.StandardHeight),
		[][3]int{{1, 2, 3}},
	)
	commands <- doric.CommandQuit

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
	tests := []struct {
		Name           string
		Command        int
		ExpectedUpdate doric.EventUpdated
	}{
		{
			Name:    "Must not move left if paused",
			Command: doric.CommandLeft,
			ExpectedUpdate: doric.EventUpdated{
				Column: doric.Column{
					Tileset: [3]int{1, 2, 3},
					X:       3,
					Y:       0,
				},
			},
		},
		{
			Name:    "Must not move right if paused",
			Command: doric.CommandRight,
			ExpectedUpdate: doric.EventUpdated{
				Column: doric.Column{
					Tileset: [3]int{1, 2, 3},
					X:       3,
					Y:       0,
				},
			},
		},
		{
			Name:    "Must not move down if paused",
			Command: doric.CommandDown,
			ExpectedUpdate: doric.EventUpdated{
				Column: doric.Column{
					Tileset: [3]int{1, 2, 3},
					X:       3,
					Y:       0,
				},
			},
		},
		{
			Name:    "Must not rotate if paused",
			Command: doric.CommandRotate,
			ExpectedUpdate: doric.EventUpdated{
				Column: doric.Column{
					Tileset: [3]int{1, 2, 3},
					X:       3,
					Y:       0,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			commands, events, timeout := setup(
				defaultConfig(),
				doric.NewWell(doric.StandardWidth, doric.StandardHeight),
				[][3]int{{1, 2, 3}},
			)

			commands <- doric.CommandPauseSwitch
			commands <- test.Command

			select {
			case ev := <-events:
				if upd, ok := ev.(doric.EventUpdated); ok && reflect.DeepEqual(upd, test.ExpectedUpdate) {
					break
				}
				t.Errorf("Current column must not move or rotate if game is paused")

			case <-timeout:
				t.Errorf("Test timed out")
			}
		})
	}
}

func TestWait(t *testing.T) {
	tests := []struct {
		name           string
		command        int
		expectedUpdate doric.EventUpdated
	}{
		{
			name:    "Must not move left if waiting",
			command: doric.CommandLeft,
			expectedUpdate: doric.EventUpdated{
				Column: doric.Column{
					Tileset: [3]int{1, 2, 3},
					X:       3,
					Y:       0,
				},
			},
		},
		{
			name:    "Must not move right if waiting",
			command: doric.CommandRight,
			expectedUpdate: doric.EventUpdated{
				Column: doric.Column{
					Tileset: [3]int{1, 2, 3},
					X:       3,
					Y:       0,
				},
			},
		},
		{
			name:    "Must not move down if waiting",
			command: doric.CommandDown,
			expectedUpdate: doric.EventUpdated{
				Column: doric.Column{
					Tileset: [3]int{1, 2, 3},
					X:       3,
					Y:       0,
				},
			},
		},
		{
			name:    "Must not rotate if waiting",
			command: doric.CommandRotate,
			expectedUpdate: doric.EventUpdated{
				Column: doric.Column{
					Tileset: [3]int{1, 2, 3},
					X:       3,
					Y:       0,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			commands, events, timeout := setup(
				defaultConfig(),
				doric.NewWell(doric.StandardWidth, doric.StandardHeight),
				[][3]int{{1, 2, 3}},
			)

			commands <- doric.CommandWaitSwitch
			commands <- test.command

			select {
			case ev := <-events:
				if upd, ok := ev.(doric.EventUpdated); ok && reflect.DeepEqual(upd, test.expectedUpdate) {
					break
				}
				t.Errorf("Current column must not move or rotate if game is waiting")

			case <-timeout:
				t.Errorf("Test timed out")
			}
		})
	}
}

func TestCommands(t *testing.T) {
	tests := []struct {
		name           string
		command        int
		expectedUpdate doric.EventUpdated
	}{
		{
			name:    "Must move left",
			command: doric.CommandLeft,
			expectedUpdate: doric.EventUpdated{
				Column: doric.Column{
					Tileset: [3]int{1, 2, 3},
					X:       2,
					Y:       0,
				},
			},
		},
		{
			name:    "Must move right",
			command: doric.CommandRight,
			expectedUpdate: doric.EventUpdated{
				Column: doric.Column{
					Tileset: [3]int{1, 2, 3},
					X:       4,
					Y:       0,
				},
			},
		},
		{
			name:    "Must move down",
			command: doric.CommandDown,
			expectedUpdate: doric.EventUpdated{
				Column: doric.Column{
					Tileset: [3]int{1, 2, 3},
					X:       3,
					Y:       1,
				},
			},
		},
		{
			name:    "Must rotate",
			command: doric.CommandRotate,
			expectedUpdate: doric.EventUpdated{
				Column: doric.Column{
					Tileset: [3]int{3, 1, 2},
					X:       3,
					Y:       0,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			commands, events, timeout := setup(
				defaultConfig(),
				doric.NewWell(doric.StandardWidth, doric.StandardHeight),
				[][3]int{{1, 2, 3}},
			)
			commands <- test.command

			select {
			case ev := <-events:
				if upd, ok := ev.(doric.EventUpdated); ok && reflect.DeepEqual(upd, test.expectedUpdate) {
					break
				}
				t.Errorf("Current column must move or rotate")

			case <-timeout:
				t.Errorf("Test timed out")
			}
		})
	}
}

func TestWellBounds(t *testing.T) {
	commands, events, timeout := setup(
		defaultConfig(),
		doric.NewWell(1, 1),
		[][3]int{{1, 2, 3}},
	)

	commands <- doric.CommandLeft

	select {
	case ev := <-events:
		if upd, ok := ev.(doric.EventUpdated); ok && upd.Column.X == 0 {
			break
		}
		t.Errorf("Current column must not move left as it clashes with well's left border")
	case <-timeout:
		t.Errorf("Test timed out")
	}

	commands <- doric.CommandRight

	select {
	case ev := <-events:
		if upd, ok := ev.(doric.EventUpdated); ok && upd.Column.X == 0 {
			break
		}
		t.Errorf("Current column must not move right as it clashes with well's right border")
	case <-timeout:
		t.Errorf("Test timed out")
	}

	commands <- doric.CommandDown

	select {
	case ev := <-events:
		if upd, ok := ev.(doric.EventUpdated); ok && upd.Column.Y == 0 {
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
		tilesets                [][3]int
		expectedWell            doric.Well
		expectedRenewedWell     doric.Well
		expectedRemoved         int
		expectedLevel           int
		expectedCurrent         [3]int
	}{
		{
			name:                    "Scored with no level up",
			numberTilesForNextLevel: 20,
			tilesets: [][3]int{
				{1, 1, 1},
				{4, 5, 6},
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
			tilesets: [][3]int{
				{1, 1, 1},
				{4, 5, 6},
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
			tilesets: [][3]int{
				{1, 1, 1},
				{4, 5, 6},
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

	for _, test := range scoredTests {
		t.Run(test.name, func(t *testing.T) {
			cfg := defaultConfig()
			cfg.InitialSpeed = 20
			cfg.MaxSpeed = 40
			cfg.NumberTilesForNextLevel = test.numberTilesForNextLevel
			_, events, timeout := setup(
				cfg,
				test.well,
				test.tilesets,
			)

			<-events

			for {
				select {
				case ev := <-events:
					switch asserted := ev.(type) {
					case doric.EventScored:
						if asserted.Removed != test.expectedRemoved {
							t.Errorf("Expected %d removed tiles but got %d", test.expectedRemoved, asserted.Removed)
						}
						if asserted.Level != test.expectedLevel {
							t.Errorf("Expected level %d but got %d", test.expectedLevel, asserted.Level)
						}
						if !reflect.DeepEqual(test.expectedWell, asserted.Well) {
							t.Errorf("Expected well %v but got %v", test.expectedWell, asserted.Well)
						}
					case doric.EventRenewed:
						if asserted.Column.Tileset != test.expectedCurrent {
							t.Errorf(
								"Expected that the next column was copied to the current one with values %v, got %v",
								test.expectedCurrent,
								asserted.Column.Tileset,
							)
						}
						if !reflect.DeepEqual(test.expectedRenewedWell, asserted.Well) {
							t.Errorf("Expected well %v but got %v", test.expectedRenewedWell, asserted.Well)
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
		tilesets                [][3]int
		expectedWells           []doric.Well
	}{
		{
			name:                    "Scored with combo",
			numberTilesForNextLevel: 20,
			tilesets: [][3]int{
				{1, 2, 3},
				{4, 5, 6},
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

	for _, test := range comboTests {
		t.Run(test.name, func(t *testing.T) {
			cfg := defaultConfig()
			cfg.InitialSpeed = 20
			cfg.NumberTilesForNextLevel = test.numberTilesForNextLevel
			_, events, timeout := setup(
				cfg,
				test.well,
				test.tilesets,
			)

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
						if !reflect.DeepEqual(test.expectedWells[count], asserted.Well) {
							t.Errorf("Expected well %v but got %v", test.expectedWells[count], asserted.Well)
						}
						if count == len(test.expectedWells)-1 {
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
	result := doric.NewWell(xl, yl)
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = slice[j][i]
		}
	}
	return result
}
