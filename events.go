package doric

// EventUpdated is sent as a response to a current column movement
type EventUpdated struct {
	Column Column
}

// EventScored is sent when the three or more tiles of the same color are aligned in the well,
// thus scoring points for the player
type EventScored struct {
	Well    Well
	Combo   int
	Removed int
	Level   int
}

// EventRenewed is sent when the current and next columns are renewed
type EventRenewed struct {
	Well        Well
	Column      Column
	NextTileset [3]int
}
