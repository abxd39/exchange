package random

import (
	"fmt"
	"math/rand"
	"time"
)

func Random6dec() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06v", rnd.Int31n(1000000))
}
