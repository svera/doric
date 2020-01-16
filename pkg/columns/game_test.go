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

var codeToStatusName = [5]string{
	"StatusUpdated",
	"StatusScored",
	"StatusRenewed",
	"StatusPaused",
	"StatusFinished",
}

func getConfig() columns.Config {
	return columns.Config{
		PointsPerTile:           10,
		NumberTilesForNextLevel: 10,
		InitialSlowdown:         10,
		Frequency:               200 * time.Millisecond,
	}
}

func TestGameOver(t *testing.T) {
	timeout := time.After(1 * time.Second)
	pit := columns.NewPit(1, pithWidth)
	current := columns.NewPiece(pit)
	next := columns.NewPiece(pit)
	r := &mocks.Randomizer{Values: []int{0}}
	game := columns.NewGame(pit, *current, *next, r, getConfig())
	updates := make(chan columns.Update)
	input := make(chan int)
	pit[0][3] = 1
	go game.Play(input, updates)
	select {
	case upd := <-updates:
		if upd.Status == columns.StatusFinished {
			break
		}
	case <-timeout:
		t.Errorf("Game should be over")
	}
}

func TestPause(t *testing.T) {
	timeout := time.After(1 * time.Second)
	pit := columns.NewPit(pitHeight, pithWidth)
	current := columns.NewPiece(pit)
	next := columns.NewPiece(pit)
	r := &mocks.Randomizer{Values: []int{1}}
	game := columns.NewGame(pit, *current, *next, r, getConfig())
	updates := make(chan columns.Update)
	input := make(chan int)
	go game.Play(input, updates)

	go func() {
		input <- columns.ActionPause
	}()

	select {
	case upd := <-updates:
		if upd.Status == columns.StatusPaused {
			break
		}
		t.Errorf("Game must be paused, got '%s'", codeToStatusName[upd.Status])
	case <-timeout:
		t.Errorf("Test timed out")
	}

	go func() {
		input <- columns.ActionPause
	}()

	select {
	case upd := <-updates:
		if upd.Status != columns.StatusPaused {
			break
		}
		t.Errorf("Game must not be paused, got '%s'", codeToStatusName[upd.Status])
	case <-timeout:
		t.Errorf("Test timed out")
	}

}

func TestInput(t *testing.T) {
	timeout := time.After(1 * time.Second)
	pit := columns.NewPit(pitHeight, pithWidth)
	current := columns.NewPiece(pit)
	next := columns.NewPiece(pit)
	r := &mocks.Randomizer{Values: []int{0, 1, 2}}
	game := columns.NewGame(pit, *current, *next, r, getConfig())
	updates := make(chan columns.Update)
	input := make(chan int)
	go game.Play(input, updates)

	go func() {
		input <- columns.ActionLeft
	}()

	select {
	case upd := <-updates:
		if upd.Current.X() == 2 {
			break
		}
		t.Errorf("Current piece must be at column %d but is at %d", 2, upd.Current.X())
	case <-timeout:
		t.Errorf("Test timed out")
	}

	go func() {
		input <- columns.ActionRight
	}()

	select {
	case upd := <-updates:
		if upd.Current.X() == 3 {
			break
		}
		t.Errorf("Current piece must be at column %d but is at %d", 3, upd.Current.X())
	case <-timeout:
		t.Errorf("Test timed out")
	}

	go func() {
		input <- columns.ActionDown
	}()

	select {
	case upd := <-updates:
		if upd.Current.Y() == 1 {
			break
		}
		t.Errorf("Current piece must be at row %d but is at %d", 1, upd.Current.Y())
	case <-timeout:
		t.Errorf("Test timed out")
	}

	go func() {
		input <- columns.ActionRotate
	}()

	select {
	case upd := <-updates:
		if upd.Current.Tiles() == [3]int{3, 1, 2} {
			break
		}
		t.Errorf("Current piece must be as %v but is as %v", [3]int{3, 1, 2}, upd.Current.Tiles())
	case <-timeout:
		t.Errorf("Test timed out")
	}

}

func TestConsolidated(t *testing.T) {
	timeout := time.After(1 * time.Second)
	pit := columns.NewPit(3, pithWidth)
	initialPit := columns.NewPit(3, pithWidth)
	current := columns.NewPiece(pit)
	next := columns.NewPiece(pit)
	r := &mocks.Randomizer{Values: []int{0, 1, 2}}
	cfg := getConfig()
	cfg.Frequency = 1 * time.Millisecond
	cfg.InitialSlowdown = 1
	game := columns.NewGame(pit, *current, *next, r, cfg)
	updates := make(chan columns.Update)
	input := make(chan int)
	go game.Play(input, updates)

	for {
		select {
		case upd := <-updates:
			if upd.Status == columns.StatusRenewed {
				if reflect.DeepEqual(initialPit, upd.Pit) {
					t.Errorf("Previous piece wasn't consolidated in pit")
				}
				if reflect.DeepEqual(current, upd.Current) {
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
			current := columns.NewPiece(pit)
			next := columns.NewPiece(pit)
			r := &mocks.Randomizer{Values: []int{0, 0, 0, 3, 4, 5}}
			cfg := getConfig()
			cfg.Frequency = 1 * time.Millisecond
			cfg.InitialSlowdown = 2
			cfg.NumberTilesForNextLevel = tt.numberTilesForNextLevel
			game := columns.NewGame(pit, *current, *next, r, cfg)
			updates := make(chan columns.Update)
			input := make(chan int)
			go game.Play(input, updates)

			for {
				select {
				case upd := <-updates:
					if upd.Status == columns.StatusScored {
						if upd.Points != tt.expectedScore {
							t.Errorf("Expected %d points but got %d", tt.expectedScore, upd.Points)
						}
						if upd.Level != tt.expectedLevel {
							t.Errorf("Expected level %d but got %d", tt.expectedLevel, upd.Level)
						}
						return
					}

					if upd.Status == columns.StatusRenewed {
						if upd.Current.Tiles() != tt.expectedTiles {
							t.Errorf(
								"Expected that the next piece was copied to the current one with values %v, got %v",
								tt.expectedTiles,
								upd.Current.Tiles(),
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
