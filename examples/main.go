package main

import (
	"fmt"
	"math/rand"
	"time"

	tl "github.com/JoelOtter/termloop"
	"github.com/svera/doric"
)

const (
	offsetX       = 32
	offsetY       = 5
	pitWidth      = 6
	pitHeight     = 13
	pointsPerTile = 10
)

var score *tl.Text
var level *tl.Text

func main() {
	actions := make(chan int)
	app := tl.NewGame()
	app.Screen().SetFps(60)
	score = tl.NewText(offsetX+15, offsetY, fmt.Sprintf("Score: %d", 0), tl.ColorWhite, tl.ColorBlack)
	level = tl.NewText(offsetX+15, offsetY+1, fmt.Sprintf("Level: %d", 1), tl.ColorWhite, tl.ColorBlack)
	pitEntity := NewPit(offsetX, offsetY, pitWidth, pitHeight)
	message := tl.NewText(offsetX+1, offsetY+5, "", tl.ColorBlack, tl.ColorWhite)
	playerEntity := NewPlayer(actions, message, offsetX, offsetY)
	nextPieceEntity := NewNext(offsetX+15, offsetY+5)

	mainLevel := tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorBlack,
	})
	setUpMainLevel(mainLevel, pitEntity, playerEntity, nextPieceEntity, message)
	app.Screen().SetLevel(mainLevel)
	startGameLogic(actions, pitEntity, playerEntity, nextPieceEntity)
	app.Start()
}

func setUpMainLevel(mainLevel *tl.BaseLevel, entities ...tl.Drawable) {
	for _, ent := range entities {
		mainLevel.AddEntity(ent)
	}
	mainLevel.AddEntity(score)
	mainLevel.AddEntity(level)
}

func startGameLogic(actions chan int, pitEntity *Pit, playerEntity *Player, nextPieceEntity *Next) {
	pit := doric.NewPit(pitHeight, pitWidth)
	cfg := doric.Config{
		NumberTilesForNextLevel: 10,
		InitialSlowdown:         10,
		Frequency:               200 * time.Millisecond,
	}

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	events := doric.Play(pit, r, cfg, actions)

	firstUpdate := <-events
	cur := firstUpdate.(doric.EventRenewed).Current
	nxt := firstUpdate.(doric.EventRenewed).Next
	pitEntity.Pit = pit
	playerEntity.Current = &cur
	nextPieceEntity.Piece = &nxt

	go func() {
		points := 0
		defer func() {
			playerEntity.Finished = true
			close(actions)
		}()
		for ev := range events {
			switch t := ev.(type) {
			case doric.EventScored:
				points += t.Removed * t.Combo * pointsPerTile
				score.SetText(fmt.Sprintf("Score: %d", points))
				level.SetText(fmt.Sprintf("Level: %d", t.Level))
			case doric.EventUpdated:
				playerEntity.Current = &t.Current
				playerEntity.Paused = t.Paused
			case doric.EventRenewed:
				pitEntity.Pit = t.Pit
				playerEntity.Current = &t.Current
				nextPieceEntity.Piece = &t.Next
			}
		}
	}()
}
