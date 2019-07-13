package actor

// New returns a new registry of actor goRoutines.
func New() Actor {
	return Actor{}
}

// Actor is type returned by the New function. It holds a registry of actor goRoutines.
type Actor struct {
}

// Execute causes the actor identified by the 'hashID' to execute the function 'fun'.
func (A Actor) Execute(hashID int64, fun func() interface{}) interface{} {
	return fun()
}
