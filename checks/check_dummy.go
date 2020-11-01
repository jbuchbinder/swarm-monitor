package checks

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	RegisterChecker("check_dummy", func() Checker { return CheckDummy{} })
}

// CheckDummy defines a check type which does nothing other than randomly
// wait and spit out OK.
type CheckDummy struct {
}

func (cd CheckDummy) Check(host string) (int, string) {
	w := rand.Int31n(3000)
	time.Sleep(time.Duration(w) * time.Millisecond)
	return StatusOK, fmt.Sprintf("slept for %d ms", w)
}
