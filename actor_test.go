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
