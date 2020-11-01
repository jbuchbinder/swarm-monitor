package checks

import (
	"errors"
	"sync"
)

const (
	CheckTypeBuiltIn = 0
	CheckTypeNagios  = 1

	StatusOK       = 0
	StatusWarning  = 1
	StatusCritical = 2
	StatusUnknown  = 3
)

var (
	checkRegistry     = map[string]func() Checker{}
	checkRegistryLock = new(sync.Mutex)
)

// Checker is an interface type which defines a SWARM internal check.
type Checker interface {
	// Check runs a check against a specified host and returns a result
	Check(string) (int, string)
}

// RegisterChecker adds a new Checker instance to the registry
func RegisterChecker(name string, m func() Checker) {
	checkRegistryLock.Lock()
	defer checkRegistryLock.Unlock()
	checkRegistry[name] = m
}

// InstantiateChecker instantiates a Checker by name
func InstantiateChecker(name string) (c Checker, err error) {
	var f func() Checker
	var found bool
	if f, found = checkRegistry[name]; !found {
		err = errors.New("unable to locate check " + name)
		return
	}
	c = f()
	err = nil
	return
}
