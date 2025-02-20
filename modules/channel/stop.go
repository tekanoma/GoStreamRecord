package channel

type Stop struct {
	stop chan bool
}

// NewStopChannel initializes a new stop channel.
func New() Stop {
	return Stop{stop: make(chan bool)}
}

// GetChannel returns the stop channel to be used by goroutines.
func (s *Stop) Get() <-chan bool {
	return s.stop
}

// Stop closes the stop channel, signaling all listeners to stop.
func (s *Stop) 	Stop() {
	close(s.stop) // Broadcasts the stop signal
}
