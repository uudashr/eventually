package eventually

// Recorder records the events.
//
// The events will be recorded the the Events slice.
type Recorder struct {
	Events []Event
}

// NewRecorder allocates and returns a new [Recorder].
func NewRecorder() *Recorder {
	return &Recorder{}
}

// Publish the event to the recorder.
func (r *Recorder) Publish(event Event) {
	r.Events = append(r.Events, event)
}
