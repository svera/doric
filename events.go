package doric

// EventUpdated is sent as a response to a current piece movement
type EventUpdated struct {
	Current Piece
}

// EventScored is sent when the three or more tiles of the same color are aligned in the pit,
// thus scoring points for the player
type EventScored struct {
	Pit     Pit
	Combo   int
	Removed int
	Level   int
}

// EventRenewed is sent when the current and next pieces are renewed
type EventRenewed struct {
	Pit     Pit
	Current Piece
	Next    Piece
}
