package actor

//TODO: Add comments explaining buffered message channel, checking for messages
//TODO: before timing out, etc., and why they are necessary to avoid race-conditions

import (
	"sync"
	"time"
)

//==============================================================================
// API
//==============================================================================

// NewWithDefaults calls New/3 with default values
func NewWithDefaults() Actor {
	return New(30*time.Second, 10, time.Second)
}

// New returns a new registry of actor goRoutines.
func New(timeout time.Duration, bufferSize int, pauseTime time.Duration) Actor {
	return Actor{
		regi:       make(map[int64]chan<- actorMessage, bufferSize),
		lock:       sync.Mutex{},
		timeout:    timeout,
		bufferSize: bufferSize,
		pauseTime:  pauseTime,
	}
}

// Actor is type returned by the New function. It holds a registry of actor
// goRoutines.
type Actor struct {
	regi       map[int64]chan<- actorMessage
	lock       sync.Mutex
	timeout    time.Duration
	bufferSize int
	pauseTime  time.Duration
}

// Execute causes the actor identified by the 'hashID' to execute the function
// 'fun'.
func (a *Actor) Execute(hashID int64, fun func() interface{}) interface{} {
	returnChan := make(chan interface{})

	msg := actorMessage{
		fun:        fun,
		returnChan: returnChan,
	}

	trySendMessage(a, hashID, msg)

	return <-returnChan
}

//==============================================================================
// Internal
//==============================================================================
type actorMessage struct {
	fun        func() interface{}
	returnChan chan<- interface{}
}

func trySendMessage(a *Actor, hashID int64, msg actorMessage) {
LOOP:
	for {
		a.lock.Lock()
		actorChan, ok := a.regi[hashID]
		if !ok {
			actorChan = newActor(a, hashID)
			a.regi[hashID] = actorChan
		}
		select {
		case actorChan <- msg:
			a.lock.Unlock()
			break LOOP
		default:
			a.lock.Unlock()
			<-time.After(a.pauseTime)
		}
	}
}

func newActor(a *Actor, hashID int64) chan<- actorMessage {
	actorChan := make(chan actorMessage, a.bufferSize)
	go actorRoutine(a, hashID, actorChan)
	return actorChan
}

func actorRoutine(a *Actor, hashID int64, actorChan <-chan actorMessage) {
LOOP:
	for {
		select {
		case <-time.After(a.timeout):
			// Prep to stop
			a.lock.Lock()
			select {
			// Check if any routine has attempted to send a message in the
			// intervening time. If so unlock and handle it
			case msg := <-actorChan:
				a.lock.Unlock()
				msg.returnChan <- msg.fun()
			// No messages have been received, so remove self from registry
			default:
				delete(a.regi, hashID)
				a.lock.Unlock()
				break LOOP
			}
		case msg := <-actorChan:
			msg.returnChan <- msg.fun()
		}
	}
}
