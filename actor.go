package actor

import (
	"sync"
)

//==============================================================================
// API
//==============================================================================

// New returns a new registry of actor goRoutines.
func New() Actor {
	return Actor{
		regi: make(map[int64]chan<- func() interface{}),
		lock: sync.Mutex{},
	}
}

// Actor is type returned by the New function. It holds a registry of actor goRoutines.
type Actor struct {
	regi map[int64]chan<- func() interface{}
	lock sync.Mutex
}

// Execute causes the actor identified by the 'hashID' to execute the function 'fun'.
func (a *Actor) Execute(hashID int64, fun func() interface{}) interface{} {

	a.lock.Lock()
	actorChan, ok := a.regi[hashID]
	if !ok {
		actorChan = newActor()
		a.regi[hashID] = actorChan
	}
	a.lock.Unlock()

	actorChan <- fun

	return nil
}

//==============================================================================
// Internal
//==============================================================================
func newActor() chan<- func() interface{} {
	actorChan := make(chan func() interface{})
	go actorRoutine(actorChan)
	return actorChan
}

func actorRoutine(actorChan <-chan func() interface{}) {
	for {
		fun := <-actorChan
		fun()
	}
}
