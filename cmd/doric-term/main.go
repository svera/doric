package main

import (
	"fmt"
	"math/rand"
	"time"

	tl "github.com/JoelOtter/termloop"
	"github.com/svera/doric/pkg/columns"
)

const (
	offsetX   = 32
	offsetY   = 5
	pithWidth = 6
	pitHeight = 13
)

var game *columns.Game
var mainLevel *tl.BaseLevel
var score *tl.Text
var level *tl.Text
var events <-chan columns.Event

func main() {
	cfg := columns.Config{
		PointsPerTile:           10,
		NumberTilesForNextLevel: 10,
		InitialSlowdown:         10,
		Frequency:               200 * time.Millisecond,
	}

	actions := make(chan int)
	app := tl.NewGame()
	app.Screen().SetFps(60)
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	pit := columns.NewPit(pitHeight, pithWidth)
	game, events = columns.NewGame(pit, r, cfg)
	score = tl.NewText(offsetX+15, offsetY, fmt.Sprintf("Score: %d", 0), tl.ColorWhite, tl.ColorBlack)
	level = tl.NewText(offsetX+15, offsetY+1, fmt.Sprintf("Level: %d", 1), tl.ColorWhite, tl.ColorBlack)
	pitEntity := NewPit(pit, offsetX, offsetY)
	message := tl.NewText(offsetX+1, offsetY+5, "", tl.ColorBlack, tl.ColorWhite)
	playerEntity := NewPlayer(actions, message, offsetX, offsetY)
	nextPieceEntity := NewNext(offsetX+15, offsetY+5)

	setUpMainLevel(pitEntity, playerEntity, nextPieceEntity, message)
	app.Screen().SetLevel(mainLevel)
	startGameLogic(actions, pitEntity, playerEntity, nextPieceEntity, message)
	app.Start()
}

func setUpMainLevel(pitEntity *Pit, playerEntity *Player, nextPieceEntity *Next, message *tl.Text) {
	mainLevel = tl.NewBaseLevel(tl.Cell{
		Bg: tl.ColorBlack,
	})
	mainLevel.AddEntity(pitEntity)
	mainLevel.AddEntity(playerEntity)
	mainLevel.AddEntity(nextPieceEntity)
	mainLevel.AddEntity(score)
	mainLevel.AddEntity(level)
	mainLevel.AddEntity(message)
}

func startGameLogic(actions chan int, pitEntity *Pit, playerEntity *Player, nextPieceEntity *Next, message *tl.Text) {
	score.SetText(fmt.Sprintf("Score: %d", 0))
	level.SetText(fmt.Sprintf("Level: %d", 1))
	go game.Play(actions)

	firstUpdate := <-events
	playerEntity.Current = &firstUpdate.Status.Current
	nextPieceEntity.Piece = &firstUpdate.Status.Next

	go func() {
		defer func() {
			playerEntity.Finished = true
			close(actions)
		}()
		for ev := range events {
			if ev.ID == columns.EventScored {
				score.SetText(fmt.Sprintf("Score: %d", ev.Status.Points))
				level.SetText(fmt.Sprintf("Level: %d", ev.Status.Level))
				pitEntity.Pit = ev.Status.Pit
			}
			if ev.ID == columns.EventUpdated {
				pitEntity.Pit = ev.Status.Pit
				playerEntity.Current = &ev.Status.Current
				playerEntity.Paused = false
				if ev.Status.Paused {
					playerEntity.Paused = true
				}
			}
			if ev.ID == columns.EventRenewed {
				playerEntity.Current = &ev.Status.Current
				nextPieceEntity.Piece = &ev.Status.Next
			}
		}
	}()
}
