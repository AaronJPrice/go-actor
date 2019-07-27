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

	userActor := NewWithDefaults()
	// run this under separate routine to avoid deadlock - makes test failures harder to decipher
	go userActor.Execute(userID, fun)

	select {
	case <-testChan:
	case <-time.After(1 * time.Second):
		t.Fatal("Actor failed to execute function.")
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

	userActor := NewWithDefaults()
	go userActor.Execute(userIDA, testFunA)
	go userActor.Execute(userIDB, testFunB)

	select {
	case <-stateChanA:
	case <-time.After(1 * time.Second):
		t.Fatal(t, "Actor failed to execute function A.")
	}

	select {
	case <-stateChanB:
	case <-time.After(1 * time.Second):
		t.Fatal(t, "Actor failed to execute function B.")
	}

	close(signalChanA)
	close(signalChanB)

	select {
	case <-stateChanA:
	case <-time.After(1 * time.Second):
		t.Fatal(t, "Actor failed to execute function A.")
	}

	select {
	case <-stateChanB:
	case <-time.After(1 * time.Second):
		t.Fatal(t, "Actor failed to execute function B.")
	}
}

func TestSameIdSerialised(t *testing.T) {
	var userID int64 = 1
	userActor := NewWithDefaults()

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

	go userActor.Execute(userID, funA)

	select {
	case <-stateChan:
	case <-time.After(1 * time.Second):
		t.Fatal(t, "Actor failed to execute function A.")
	}

	go userActor.Execute(userID, funB)

	select {
	case <-stateChan:
		t.Fatal(t, "Function B started execution before A completed.")
	case <-time.After(1 * time.Second):
	}

	close(signalChan)

	select {
	case msg := <-stateChan:
		if msg != "doneA" {
			t.Fatal(t, "Functions executing in wrong order.")
		}
	case <-time.After(1 * time.Second):
		t.Fatal(t, "Function A failed to complete.")
	}

	select {
	case msg := <-stateChan:
		if msg != "executingB" {
			t.Fatal(t, "Functions executing in wrong order.")
		}
	case <-time.After(1 * time.Second):
		t.Fatal(t, "Function B failed to execute.")
	}
}

func TestReturn(t *testing.T) {
	expected := 17
	userActor := NewWithDefaults()

	fun := func() interface{} {
		return expected
	}

	result := userActor.Execute(1, fun)

	assert(t, result, expected, "Execute did not return expected value.")
}

func TestExecutionAfterTimeout(t *testing.T) {
	var userID int64 = 1
	testChan := make(chan interface{})
	fun := func() interface{} {
		testChan <- 0
		return nil
	}

	userActor := NewWithDefaults()
	go userActor.Execute(userID, fun)

	select {
	case <-testChan:
	case <-time.After(1 * time.Second):
		t.Fatal(t, "Actor failed to execute function.")
	}

	<-time.After(50 * time.Millisecond)

	go userActor.Execute(userID, fun)

	select {
	case <-testChan:
	case <-time.After(1 * time.Second):
		t.Fatal(t, "Actor failed to execute second function sent after timeout.")
	}
}

// TODO
// Not sure how to test these (without depending on internal implementation details)...
// func TestActorDiesAfterTimeout(t *testing.T) {}
// func TestActorDoesNotDieBeforeTimeout(t *testing.T) {}

//==============================================================================
// Utilities
//==============================================================================
func assert(t *testing.T, value interface{}, expected interface{}, failMsg string) {
	if value != expected {
		t.Fatal(failMsg)
	}
}
