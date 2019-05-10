package columns

import (
	"math/rand"
	"time"
)

var source rand.Source

func init() {
	source = rand.NewSource(time.Now().UnixNano())
}
