package actor

import (
	"testing"
	"time"
)

func TestExecution(t *testing.T) {
	var userID int64 = 1
	testChan := make(chan interface{})
	fun := func() interface{} {
		close(testChan)
		return nil
	}

	userActor := New()
	userActor.Execute(userID, fun)

	select {
	case <-testChan:
	case <-time.After(1 * time.Second):
		t.Error("Actor failed to execute function.")
	}
}

func TestDiffHashIDConcurrentExecution(t *testing.T) {
	var userIDA int64 = 1
	var userIDB int64 = 2

	concurrentTestFunGen := func(stateChan chan<- string, signalChan <-chan interface{}) func() interface{} {
		return func() interface{} {
			stateChan <- "executing"
			<-signalChan
			stateChan <- "done"
			return nil
		}
	}

	stateChanA := make(chan string)
	signalChanA := make(chan interface{})
	testFunA := concurrentTestFunGen(stateChanA, signalChanA)

	stateChanB := make(chan string)
	signalChanB := make(chan interface{})
	testFunB := concurrentTestFunGen(stateChanB, signalChanB)

	userActor := New()
	userActor.Execute(userIDA, testFunA)
	userActor.Execute(userIDB, testFunB)

	select {
	case <-stateChanA:
	case <-time.After(1 * time.Second):
		t.Error("Actor failed to execute function A.")
	}

	select {
	case <-stateChanB:
	case <-time.After(1 * time.Second):
		t.Error("Actor failed to execute function B.")
	}

	close(signalChanA)
	close(signalChanB)

	select {
	case <-stateChanA:
	case <-time.After(1 * time.Second):
		t.Error("Actor failed to execute function A.")
	}

	select {
	case <-stateChanB:
	case <-time.After(1 * time.Second):
		t.Error("Actor failed to execute function B.")
	}
}
