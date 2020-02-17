package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	tl "github.com/JoelOtter/termloop"
	"github.com/svera/doric"
)

const (
	offsetX       = 32
	offsetY       = 5
	pointsPerTile = 10
)

func main() {
	commands := make(chan int)
	app := tl.NewGame()
	app.Screen().SetFps(60)

	mainLevel := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorBlack,
	})
	app.Screen().SetLevel(mainLevel)
	entities := startGameLogic(commands)
	setUpMainLevel(mainLevel, entities)
	app.Start()
}

func setUpMainLevel(mainLevel *tl.BaseLevel, entities []tl.Drawable) {
	for _, ent := range entities {
		mainLevel.AddEntity(ent)
	}
}

func startGameLogic(commands chan int) []tl.Drawable {
	well := doric.NewWell(doric.StandardHeight, doric.StandardWidth)
	cfg := doric.Config{
		NumberTilesForNextLevel: 10,
		InitialSpeed:            0.5,
		SpeedIncrement:          0.25,
		MaxSpeed:                13,
	}

	factory := func(n int) [3]int {
		source := rand.NewSource(time.Now().UnixNano())
		rand.New(source)
		return [3]int{
			rand.Intn(n) + 1,
			rand.Intn(n) + 1,
			rand.Intn(n) + 1,
		}
	}
	events := doric.Play(well, factory, cfg, commands)

	firstUpdate := <-events
	cur := firstUpdate.(doric.EventRenewed).Current
	nxt := firstUpdate.(doric.EventRenewed).Next

	mux := &sync.Mutex{}
	message := tl.NewText(offsetX+1, offsetY+5, "", tl.ColorBlack, tl.ColorWhite)
	wellEntity := NewWell(well, offsetX, offsetY, doric.StandardHeight, doric.StandardWidth, mux)
	playerEntity := NewPlayer(&cur, commands, message, offsetX, offsetY, mux)
	nextColumnEntity := NewNext(&nxt, offsetX+15, offsetY+5, mux)
	score := tl.NewText(offsetX+15, offsetY, fmt.Sprintf("Score: %d", 0), tl.ColorWhite, tl.ColorBlack)
	level := tl.NewText(offsetX+15, offsetY+1, fmt.Sprintf("Level: %d", 1), tl.ColorWhite, tl.ColorBlack)

	go func() {
		points := 0
		defer func() {
			playerEntity.Finished = true
			close(commands)
		}()
		for ev := range events {
			switch t := ev.(type) {
			case doric.EventScored:
				mux.Lock()
				points += t.Removed * t.Combo * pointsPerTile
				score.SetText(fmt.Sprintf("Score: %d", points))
				level.SetText(fmt.Sprintf("Level: %d", t.Level))
				mux.Unlock()
			case doric.EventUpdated:
				mux.Lock()
				playerEntity.Current = &t.Current
				mux.Unlock()
			case doric.EventRenewed:
				mux.Lock()
				wellEntity.Well = t.Well
				playerEntity.Current = &t.Current
				nextColumnEntity.Column = &t.Next
				mux.Unlock()
			}
		}
	}()

	return []tl.Drawable{wellEntity, playerEntity, nextColumnEntity, message, score, level}
}
