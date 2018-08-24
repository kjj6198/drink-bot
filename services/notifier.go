package services

// Notifier has Notify method, which should implement Notify logic.
type Notifier interface {
	Notify() (success bool, err error)
	Done() <-chan bool
}
