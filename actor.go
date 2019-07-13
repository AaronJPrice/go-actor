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
		regi: make(map[int64]chan<- actorMessage),
		lock: sync.Mutex{},
	}
}

// Actor is type returned by the New function. It holds a registry of actor
// goRoutines.
type Actor struct {
	regi map[int64]chan<- actorMessage
	lock sync.Mutex
}

// Execute causes the actor identified by the 'hashID' to execute the function
// 'fun'.
func (a *Actor) Execute(hashID int64, fun func() interface{}) interface{} {
	returnChan := make(chan interface{})

	msg := actorMessage{
		fun:        fun,
		returnChan: returnChan,
	}

	a.lock.Lock()
	actorChan, ok := a.regi[hashID]
	if !ok {
		actorChan = newActor()
		a.regi[hashID] = actorChan
	}
	a.lock.Unlock()

	actorChan <- msg

	return <-returnChan
}

//==============================================================================
// Internal
//==============================================================================
type actorMessage struct {
	fun        func() interface{}
	returnChan chan<- interface{}
}

func newActor() chan<- actorMessage {
	actorChan := make(chan actorMessage)
	go actorRoutine(actorChan)
	return actorChan
}

func actorRoutine(actorChan <-chan actorMessage) {
	for {
		msg := <-actorChan
		msg.returnChan <- msg.fun()
	}
}
