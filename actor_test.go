package actor

import (
	"testing"
	"time"
)

//==============================================================================
// Tests
//==============================================================================
func TestExecution(t *testing.T) {
	var userID int64 = 1
	testChan := make(chan interface{})
	fun := func() interface{} {
		close(testChan)
		return nil
	}

	userActor := New()
	// run this under separate routine to avoid deadlock - makes test failures harder to decipher
	go userActor.Execute(userID, fun)

	select {
	case <-testChan:
	case <-time.After(1 * time.Second):
		failNow(t, "Actor failed to execute function.")
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
	go func() {
		userActor.Execute(userIDA, testFunA)
		userActor.Execute(userIDB, testFunB)
	}()

	select {
	case <-stateChanA:
	case <-time.After(1 * time.Second):
		failNow(t, "Actor failed to execute function A.")
	}

	select {
	case <-stateChanB:
	case <-time.After(1 * time.Second):
		failNow(t, "Actor failed to execute function B.")
	}

	close(signalChanA)
	close(signalChanB)

	select {
	case <-stateChanA:
	case <-time.After(1 * time.Second):
		failNow(t, "Actor failed to execute function A.")
	}

	select {
	case <-stateChanB:
	case <-time.After(1 * time.Second):
		failNow(t, "Actor failed to execute function B.")
	}
}

func TestSameIdSerialised(t *testing.T) {
	var userID int64 = 1
	userActor := New()

	stateChan := make(chan string)
	signalChan := make(chan interface{})

	funA := func() interface{} {
		stateChan <- "executingA"
		<-signalChan
		stateChan <- "doneA"
		return nil
	}

	funB := func() interface{} {
		stateChan <- "executingB"
		return nil
	}

	go func() {
		userActor.Execute(userID, funA)
		userActor.Execute(userID, funB)
	}()

	select {
	case <-stateChan:
	case <-time.After(1 * time.Second):
		failNow(t, "Actor failed to execute function A.")
	}

	select {
	case <-stateChan:
		failNow(t, "Function B started execution before A completed.")
	case <-time.After(1 * time.Second):
	}

	close(signalChan)

	select {
	case msg := <-stateChan:
		if msg != "doneA" {
			failNow(t, "Functions executing in wrong order.")
		}
	case <-time.After(1 * time.Second):
		failNow(t, "Function A failed to complete.")
	}

	select {
	case msg := <-stateChan:
		if msg != "executingB" {
			failNow(t, "Functions executing in wrong order.")
		}
	case <-time.After(1 * time.Second):
		failNow(t, "Function B failed to execute.")
	}
}

//==============================================================================
// Utilities
//==============================================================================
func failNow(t *testing.T, message string) {
	t.Error(message)
	t.FailNow()
}
