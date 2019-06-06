# Doric

A Columns game implementation written by Sergio Vera.

## Features

* The classic SEGA arcade game.
* Game logic completely isolated from presentation, running in its own thread. [pkg/columns]() library can be used in other implementations with minimal effort. Basically:
```go
    import "github.com/svera/doric/pkg/columns"

    func main() {
        events := make(chan int)
        pit := columns.NewPit(13, 6)
        player := columns.NewPlayer(pit)
        player.Play(events)
    }
```
* A sample client using the library can be found at [cmd/doric-term]().

## Build from sources

### Requirements

* Go 1.11 or higher

### Instructions

 1. In a terminal, run `go get github.com/svera/doric`
 2. From the source code directory, run `go build ./cmd/doric-term`.

## How to play

The objective of the game is to get the maximum possible score. To do that, player must eliminate falling pieces from the pit, aligning
3 or more tiles of the same color vertically, horizontally or diagonally. Every 10 tiles removed the falling speed increases slightly.